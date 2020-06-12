// Package service 提供一个默认的处理模型
package service

import (
	"chat/config"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

// ErrClientIDNotFound 找不到 client id
var ErrClientIDNotFound error

func init() {
	ErrClientIDNotFound = errors.New("找不到 client id")
}

// DefaultDispatcher 默认适配器
type DefaultDispatcher struct {
	ConnList  map[string]*Client
	register  chan *Client
	destroy   chan *Client
	broadcast chan Message
	cancel    context.CancelFunc
	ctx       context.Context
	Len       int
}

// Init 初始化
func (th *DefaultDispatcher) Init() {
	th.ctx, th.cancel = context.WithCancel(context.Background())
	th.ConnList = make(map[string]*Client)
	th.register = make(chan *Client, 1000)
	th.destroy = make(chan *Client, 1000)
	th.broadcast = make(chan Message, 1000)
	log.Println("Init")
	go th.run()
	go th.HeartBeat()
}

// Close 收尾
func (th *DefaultDispatcher) Close() {
	log.Println("Close")
	th.cancel()
}

// Context 获取文本流
func (th *DefaultDispatcher) Context() context.Context {
	return th.ctx
}

// Status 当前的状态
func (th *DefaultDispatcher) Status() map[string]interface{} {
	list := []string{}

	for id := range th.ConnList {
		list = append(list, id)
	}

	tmp := map[string]interface{}{
		"register_stacking":  fmt.Sprintf("%d/%d", len(th.register), cap(th.register)),
		"broadcast_stacking": fmt.Sprintf("%d/%d", len(th.broadcast), cap(th.broadcast)),
		"destroy_stacking":   fmt.Sprintf("%d/%d", len(th.destroy), cap(th.destroy)),
		"connect_client_id":  list,
	}
	return tmp
}

// RegisterEvent 注册client事件
func (th *DefaultDispatcher) RegisterEvent(cli *Client) error {
	log.Printf("RegisterEvent => cli: %#+v \n", cli)
	th.register <- cli
	return nil
}

// DestroyEvent 销毁client事件
func (th *DefaultDispatcher) DestroyEvent(cli *Client) error {
	log.Printf("DestroyEvent =>  cli: %#+v \n", cli)
	th.destroy <- cli
	return nil
}

// HeartBeatEvent 心跳检测事件
func (th *DefaultDispatcher) HeartBeatEvent(msg Message, cli *Client) error {
	// 这里客户端发来的心跳ack不处理
	log.Printf("HeartBeatEvent =>  msg: %#+v, cli: %#+v \n", msg, cli)
	return nil
	// return cli.HeartBeat()
}

// BroadcastEvent 广播事件
func (th *DefaultDispatcher) BroadcastEvent(msg Message, cli *Client) error {
	log.Printf("BroadcastEvent =>  msg: %#+v, cli: %#+v \n", msg, cli)
	msg.From = cli.ID
	th.broadcast <- msg
	return nil
}

// GetClientByID 获取指定客户端
func (th *DefaultDispatcher) GetClientByID(id string) (cli *Client, err error) {
	cli, ok := th.ConnList[id]
	if !ok {
		return nil, ErrClientIDNotFound
	}
	return
}

// DefaultMessageEvent 默认发送消息事件
func (th *DefaultDispatcher) DefaultMessageEvent(MessageType int, msg Message, cli *Client) error {
	log.Printf("DefaultEvent => msgType: %#+v  msg: %#+v, cli: %#+v \n", MessageType, msg, cli)
	return nil
}

// HeartBeat 定时检测连接健康程度, 失联的就断开链接
func (th *DefaultDispatcher) HeartBeat() {
	ticker := time.NewTicker(config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("heartBeat")
			for id := range th.ConnList {
				th.ConnList[id].HeartBeat()
			}
		case <-th.ctx.Done():
			log.Println("关闭 heart beat")
			return
		}
	}
}

// 启动消息接收
func (th *DefaultDispatcher) run() {
	log.Println("开始监听注册事件")

LOOP:
	for {
		select {
		case client := <-th.register:
			log.Println("注册事件:", client.ID)
			th.ConnList[client.ID] = client
			th.Len++
		case client := <-th.destroy:
			log.Println("销毁事件:", client.ID)
			err := client.conn.Close()
			if err != nil {
				log.Printf("destroy Error: %v \n", err)
			}
			delete(th.ConnList, client.ID)
			th.Len--
		case message := <-th.broadcast:
			log.Printf("广播事件: %#+v \n", message)
			for id := range th.ConnList {
				if id != message.From { // 自己发送的消息自己不会收到
					client := th.ConnList[id]
					err := client.SendText(message)
					if err != nil {
						log.Println("broadcastError: ", err)
					}
				}
			}
		case <-th.ctx.Done():
			log.Println("终止事件")
			for id := range th.ConnList {
				conn := th.ConnList[id]
				// 向所有在线链接发送断开提示
				conn.SendMessage(DisconnectedMessage, "")
				conn.Close()
			}
			break LOOP
		}
	}
}
