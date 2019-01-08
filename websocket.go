package chat

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// SystemMessage 系统消息
	SystemMessage = iota
	// BroadcastMessage 广播消息(正常的消息)
	BroadcastMessage
	// HeartBeatMessage 心跳消息
	HeartBeatMessage
	// ConnectedMessage 上线通知
	ConnectedMessage
	// DisconnectedMessage 下线通知
	DisconnectedMessage
)

var aliveList *AliveList
var upgrader = websocket.Upgrader{}

// AliveList 当前在线列表
type AliveList struct {
	ConnList  map[string]*Client
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
	Content string
	SentAt  int64
	Type    int
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
		ConnList:  make(map[string]*Client, 100),
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
			al.ConnList[client.ID] = client
			al.Len++
			al.SysBroadcast(ConnectedMessage, Message{
				ID:      client.ID,
				Content: "connected",
				SentAt:  time.Now().Unix(),
			})

		case client := <-al.destroy:
			log.Println("销毁事件:", client.ID)
			err := client.conn.Close()
			if err != nil {
				log.Printf("destroy Error: %v \n", err)
			}
			delete(al.ConnList, client.ID)
			al.Len--

		case message := <-al.broadcast:
			log.Printf("广播事件: %s %s %d \n", message.ID, message.Content, message.Type)
			for id := range al.ConnList {
				if id != message.ID {

					err := al.sendMessage(id, message)
					if err != nil {
						log.Println("broadcastError: ", err)
					}
				}
			}

		case sign := <-al.cancel:
			log.Println("终止事件: ", sign)
			os.Exit(0)
		}
	}
}

func (al *AliveList) sendMessage(id string, msg Message) error {
	if conn, ok := al.ConnList[id]; ok {
		return conn.SendMessage(msg.Type, msg.Content)
	}
	return fmt.Errorf("conn not found: %v", msg)
}

// Register 注册
func (al *AliveList) Register(client *Client) {
	al.register <- client
}

// Destroy 销毁
func (al *AliveList) Destroy(client *Client) {
	al.destroy <- client
}

// Broadcast 个人广播消息
func (al *AliveList) Broadcast(message Message) {
	al.broadcast <- message
}

// SysBroadcast 系统广播 这里加了一个消息类型, 正常的broadcast应该就是 BroadcastMessage 类型消息
func (al *AliveList) SysBroadcast(messageType int, message Message) {
	message.Type = messageType
	al.Broadcast(message)
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
		conn:   conn,
		ID:     id,
		cancel: make(chan int, 1),
	}

	aliveList.Register(client)
	return
}

// Broadcast 单个客户端的广播事件
func (cli *Client) Broadcast(msg string) {
	aliveList.Broadcast(Message{
		ID:      cli.ID,
		Content: msg,
		Type:    BroadcastMessage,
		SentAt:  time.Now().Unix(),
	})
}

// SendMessage 单个链接发送消息
func (cli *Client) SendMessage(messageType int, message string) error {

	msg := Message{
		ID:      cli.ID,
		Content: message,
		SentAt:  time.Now().Unix(),
		Type:    messageType,
	}
	// 这里固定是
	err := cli.conn.WriteJSON(msg)
	if err != nil {
		log.Println("sendMessageError :", err)
		log.Println("message: ", msg)
		log.Println("cli: ", cli)
		cli.Close()
	}
	return err
}

// Close 单个链接断开 (这里可以加一个参数, 进行区分关闭链接时的状态, 比如0:正常关闭,1:非正常关闭 etc..)
func (cli *Client) Close() {
	cli.cancel <- 1
	aliveList.Broadcast(Message{
		ID:      cli.ID,
		Content: "",
		Type:    DisconnectedMessage,
	})
	aliveList.Destroy(cli)
}

// HeartBeat 服务端检测链接是否正常
func (cli *Client) HeartBeat() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cli.SendMessage(HeartBeatMessage, "heart beat")
		case <-cli.cancel:
			log.Println("即将关闭定时器", cli)
			close(cli.cancel)
			return
		}
	}
}
