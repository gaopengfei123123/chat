package chat

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// SystemMessage 系统消息 0
	SystemMessage = iota
	// BroadcastMessage 广播消息(正常的消息) 1
	BroadcastMessage
	// HeartBeatMessage 心跳消息(暂时不处理)  2
	HeartBeatMessage
	// ConnectedMessage 上线通知 3
	ConnectedMessage
	// DisconnectedMessage 下线通知 4
	DisconnectedMessage
	// BreakMessage 服务断开链接通知(服务端关闭) 5
	BreakMessage
	// RegisterMessage 注册事件消息 6
	RegisterMessage
)

// 跨域配置
var upgrader = websocket.Upgrader{}

// Client socket客户端
type Client struct {
	ID         string          // 链接的唯一标识
	conn       *websocket.Conn // 链接实体
	cancel     chan int
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

// Message 消息体结构
type Message struct {
	ID         string      // 发送消息id
	Content    string      // 消息内容
	SentAt     int64       `json:"sent_at"` // 发送时间
	Type       int         // 消息类型, 如 BroadcastMessage
	From       string      // 发送人client id
	To         []string    // 接收人client id, 根据消息类型来说, 单发, 群发, 广播什么的, 具体处理在Event中处理
	FromUserID string      `json:"from_user_id"` // 发送者用户业务id
	ToUserID   string      `json:"to_user_id"`   // 接受者用户业务id
	Ext        interface{} `json:"ext"`          // 扩展字段, 按需使用
}

func init() {
	// 允许跨域请求
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}

// NewSocketClient 新建客户端链接
func NewSocketClient(id string, w http.ResponseWriter, r *http.Request) (client *Client, err error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 获取整个集合的文本流
	handler := GetHandler()
	ctx, cancel := context.WithCancel(handler.Context())
	client = &Client{
		conn:       conn,
		ID:         id,
		cancel:     make(chan int, 1),
		Ctx:        ctx,
		CancelFunc: cancel,
	}

	handler.RegisterEvent(client)
	return
}

// Broadcast 单个客户端的广播事件
func (cli *Client) Broadcast(msg string) {
	handler := GetHandler()
	handler.BroadcastEvent(Message{
		ID:      randSeq(32),
		Content: msg,
		Type:    BroadcastMessage,
		SentAt:  time.Now().Unix(),
	}, cli)
}

// SendMessage 单个链接发送消息, 默认模板
func (cli *Client) SendMessage(messageType int, message string) error {

	if messageType == BreakMessage {
		err := cli.conn.WriteMessage(websocket.CloseMessage, []byte("close"))
		return err
	}

	msg := Message{
		ID:      cli.ID,
		Content: message,
		SentAt:  time.Now().Unix(),
		Type:    messageType,
	}

	err := cli.SendText(msg)
	if err != nil {
		log.Println("sendMessageError :", err)
		log.Println("message: ", msg)
		log.Printf("cli: %#+v \n", cli)
		cli.Close()
	}
	return err
}

// SendText 发送文本类消息
func (cli *Client) SendText(msg Message) error {
	return cli.conn.WriteJSON(msg)
}

// Close 单个链接断开 (这里可以加一个参数, 进行区分关闭链接时的状态, 比如0:正常关闭,1:非正常关闭 etc..)
func (cli *Client) Close() {
	cli.CancelFunc()
	handler := GetHandler()
	handler.DestroyEvent(cli)
}
