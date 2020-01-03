// Package chat 处理socket请求
package chat

import (
	"encoding/json"
	// "fmt"
	"log"
)

// 默认适配器
var adapter EventInterface

// SetAdapter 设置业务事件
func SetAdapter(adp EventInterface) {
	adapter = adp
}

// EventInterface 事件接口
type EventInterface interface {
	RegisterEvent(msg Message, cli *Client) error
	HeartBeatEvent(msg Message, cli *Client) error
	BroadcastEvent(msg Message, cli *Client) error
	DefaultEvent(msg Message, cli *Client) error
}

func init() {
	SetAdapter(DefaultAdapter{})
}

// HandleRequest 处理
func HandleRequest(cli *Client, msg []byte) (err error) {
	log.Printf("获取信息: %s \n", msg)

	var msgBody Message
	err = json.Unmarshal(msg, &msgBody)

	if err != nil {
		return
	}

	switch msgBody.Type {
	case BroadcastMessage, SystemMessage:
		adapter.BroadcastEvent(msgBody, cli)
	case HeartBeatMessage:
		adapter.HeartBeatEvent(msgBody, cli)
	default:
		// 默认通知所有消息(虽然这不对)
		adapter.DefaultEvent(msgBody, cli)
	}
	return
}
