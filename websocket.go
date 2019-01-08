package chat

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

var aliveList *AliveList
var upgrader = websocket.Upgrader{}

// AliveList 当前在线列表
type AliveList struct {
	ConnList  map[string]*websocket.Conn
	register  chan *Client
	destroy   chan *Client
	broadcast chan Message
	cancel    chan int
	Len       int
}

// Client socket客户端
type Client struct {
	ID     string
	conn   *websocket.Conn
	cancel chan int
}

// Message 消息体结构
type Message struct {
	ID      string
	Content []byte
}

func init() {
	// 允许跨域请求
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	aliveList = NewAliveList()
	go aliveList.run()
}

// NewAliveList 初始化
func NewAliveList() *AliveList {
	return &AliveList{
		ConnList:  make(map[string]*websocket.Conn, 100),
		register:  make(chan *Client, 100),
		destroy:   make(chan *Client, 100),
		broadcast: make(chan Message, 100),
		cancel:    make(chan int),
		Len:       0,
	}
}

// 启动监听
func (al *AliveList) run() {
	log.Println("开始监听注册事件")
	for {
		select {
		case client := <-al.register:
			log.Println("注册事件:", client.ID)
			al.ConnList[client.ID] = client.conn
			al.Len++

		case client := <-al.destroy:
			log.Println("销毁事件:", client.ID)
			err := client.conn.Close()
			if err != nil {
				log.Printf("destroy Error: %v \n", err)
			}
			delete(al.ConnList, client.ID)
			al.Len--

		case message := <-al.broadcast:
			log.Printf("广播事件: %s %s \n", message.ID, message.Content)
			auth := []byte(message.ID + "说: ")
			msg := append(auth, message.Content...)

			for id := range al.ConnList {
				if id != message.ID {
					al.sendMessage(id, msg)
				}
			}

		case sign := <-al.cancel:
			log.Println("终止事件: ", sign)
			os.Exit(0)
		}
	}
}

func (al *AliveList) sendMessage(id string, msg []byte) error {
	return al.ConnList[id].WriteMessage(websocket.TextMessage, msg)
}

// Register 注册
func (al *AliveList) Register(client *Client) {
	al.register <- client
}

// Destroy 销毁
func (al *AliveList) Destroy(client *Client) {
	al.destroy <- client
}

// Broadcast 广播消息
func (al *AliveList) Broadcast(message Message) {
	al.broadcast <- message
}

// Cancel 关闭集合
func (al *AliveList) Cancel() {
	al.cancel <- 1
}

// NewWebSocket 新建客户端链接
func NewWebSocket(id string, w http.ResponseWriter, r *http.Request) (client *Client, err error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client = &Client{
		conn: conn,
		ID:   id,
	}

	aliveList.Register(client)
	return
}

// Broadcast 单个客户端的广播事件
func (cli *Client) Broadcast(msg []byte) {
	aliveList.Broadcast(Message{ID: cli.ID, Content: msg})
}

// SendMessage 单个链接发送消息
func (cli *Client) SendMessage(messageType int, message []byte) error {
	err := cli.conn.WriteMessage(messageType, message)
	if err != nil {
		log.Println("sendMessageError :", err)
		cli.Close()
	}
	return err
}

// Close 单个链接断开
func (cli *Client) Close() {
	cli.cancel <- 1
	aliveList.Broadcast(Message{ID: cli.ID, Content: []byte(cli.ID + "小老弟下线了")})
	aliveList.Destroy(cli)
}

// HeartBeat 服务端检测链接是否正常
func (cli *Client) HeartBeat() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cli.SendMessage(websocket.TextMessage, []byte("heart beat"))
		case <-cli.cancel:
			return
		}
	}
}
