package service

import (
	"chat/config"
	"chat/library"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// NewSocketClient 新建客户端链接
func NewSocketClient(id string, w http.ResponseWriter, r *http.Request) (client *Client, err error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 获取整个集合的文本流
	dispatcher := GetDispatcher()
	ctx, cancel := context.WithCancel(dispatcher.Context())
	client = &Client{
		conn:       conn,
		ID:         id,
		Ctx:        ctx,
		CancelFunc: cancel,
	}

	dispatcher.RegisterEvent(client)
	return
}

// Broadcast 单个客户端的广播事件
func (cli *Client) Broadcast(msg string) {
	GetDispatcher().BroadcastEvent(Message{
		ID:      library.RandSeq(32),
		Content: msg,
		From:    cli.ID,
		Type:    BroadcastMessage,
		SentAt:  time.Now().Unix(),
	}, cli)
}

// SysBroadcast 单个客户端的系统广播事件
func (cli *Client) SysBroadcast(msg string) {
	cli.DispatchRequest(Message{
		ID:      library.RandSeq(32),
		Content: msg,
		From:    cli.ID,
		Type:    SystemMessage,
		SentAt:  time.Now().Unix(),
	})
}

// ReadMessage 读消息
func (cli *Client) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = cli.conn.ReadMessage()
	cli.lastRequestTime = time.Now()
	cli.retryTime = 0
	log.Printf("更新时间 %d \n", cli.lastRequestTime)
	return
}

// HeartBeat 心跳检测
// 最近接收过消息就不做心跳检测
func (cli *Client) HeartBeat() (err error) {
	log.Printf("start check heartBeat, ID: %s \n", cli.ID)
	last := time.Now().Sub(cli.lastRequestTime)

	if last < config.HeartbeatInterval {
		return nil
	}

	if cli.retryTime > config.MaxRetryTime {
		cli.Close()
	}

	log.Printf("当前客户端: %#+v \n", cli)
	msg := Message{
		Type:    HeartBeatMessage,
		Content: "",
	}
	cli.retryTime++
	err = cli.SendText(msg)

	return nil
}

// DispatchRequest 分发消息
func (cli *Client) DispatchRequest(msg Message) (err error) {
	log.Printf("获取信息: %#+v \n", msg)

	dispatcher := GetDispatcher()
	switch msg.Type {
	case BroadcastMessage, SystemMessage:
		err = dispatcher.BroadcastEvent(msg, cli)
	case HeartBeatMessage:
		err = dispatcher.HeartBeatEvent(msg, cli)
	default:
		// 自定义的type就层层套娃一下
		err = dispatcher.DefaultMessageEvent(msg.Type, msg, cli)
	}
	return
}

// SendMessage 单个链接快速发送消息, 默认模板
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
	msg.SentAt = time.Now().Unix()
	return cli.conn.WriteJSON(msg)
}

// Close 单个链接断开 (这里可以加一个参数, 进行区分关闭链接时的状态, 比如0:正常关闭,1:非正常关闭 etc..)
func (cli *Client) Close() {
	cli.CancelFunc()
	dispatcher := GetDispatcher()
	dispatcher.DestroyEvent(cli)
}

// DispatchRequest 分发请求
func DispatchRequest(cli *Client, msg []byte) (err error) {
	log.Printf("获取信息: %s \n", msg)

	var msgBody Message
	err = json.Unmarshal(msg, &msgBody)

	if err != nil {
		return
	}

	// log.Printf("MessageBody: %#+v \n", msgBody)

	switch msgBody.Type {
	case BroadcastMessage, SystemMessage:
		socketHandler.BroadcastEvent(msgBody, cli)
	case HeartBeatMessage:
		socketHandler.HeartBeatEvent(msgBody, cli)
	default:
		// 自定义的type就层层套娃一下
		socketHandler.DefaultMessageEvent(msgBody.Type, msgBody, cli)
	}
	return
}
