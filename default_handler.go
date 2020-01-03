// Package chat 提供一个默认的处理模型
package chat

import (
	"fmt"
)

// DefaultAdapter 默认适配器
type DefaultAdapter struct {
}

// RegisterEvent 注册事件
func (th DefaultAdapter) RegisterEvent(msg Message, cli *Client) error {
	fmt.Printf("RegisterEvent =>  msg: %v, cli: %v \n", msg, cli)
	return nil
}

// HeartBeatEvent 心跳检测事件
func (th DefaultAdapter) HeartBeatEvent(msg Message, cli *Client) error {
	fmt.Printf("HeartBeatEvent =>  msg: %v, cli: %v \n", msg, cli)
	return nil
}

// BroadcastEvent 广播事件
func (th DefaultAdapter) BroadcastEvent(msg Message, cli *Client) error {
	fmt.Printf("BroadcastEvent =>  msg: %v, cli: %v \n", msg, cli)
	return nil
}

// DefaultEvent 默认事件处理
func (th DefaultAdapter) DefaultEvent(msg Message, cli *Client) error {
	fmt.Printf("DefaultEvent =>  msg: %v, cli: %v \n", msg, cli)
	return nil
}
