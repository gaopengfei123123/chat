// Package service 处理socket请求
package service

import (
	"context"
)

// 默认适配器 全局存储
var socketHandler EventInterface

// EventInterface 事件接口
type EventInterface interface {
	RegisterEvent(cli *Client) error
	DestroyEvent(cli *Client) error
	HeartBeatEvent(msg Message, cli *Client) error
	BroadcastEvent(msg Message, cli *Client) error
	DefaultMessageEvent(MessageType int, msg Message, cli *Client) error
	GetClientByID(id string) (cli *Client, err error)
	Context() context.Context
	Init()
	Close()
	Status() map[string]interface{}
}

func init() {
}

// Start 开启ws服务监听
func Start() {
	SetDispatcher(&DefaultDispatcher{})
}

// SetDispatcher 设置调度器
func SetDispatcher(adp EventInterface) {
	// 目前只能有一个去处理业务
	if socketHandler != nil {
		socketHandler.Close()
	}
	socketHandler = adp
	socketHandler.Init()
}

// GetDispatcher 获取业务处理器
func GetDispatcher() EventInterface {
	return socketHandler
}
