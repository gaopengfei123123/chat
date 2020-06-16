// Package service 处理socket请求
package service

import (
	"context"
)

// 默认适配器 全局存储, 只能根据方法修改
var socketHandler EventInterface

// EventInterface 事件接口
type EventInterface interface {
	RegisterEvent(cli *Client) error                                     // 将链接加入到整个list当中
	DestroyEvent(cli *Client) error                                      // 从list中移除
	HeartBeatEvent(msg Message, cli *Client) error                       // 心跳检测
	BroadcastEvent(msg Message, cli *Client) error                       // 广播
	DefaultMessageEvent(MessageType int, msg Message, cli *Client) error // 识别不出type的都扔这里
	GetClientByID(clientID string) (cli *Client, err error)              // 根据clientID获取实体
	Context() context.Context                                            // 所有的链接应该都继承自这个context
	Init()                                                               // 初始化用
	Close()                                                              // 关闭服务的收尾
	Status() map[string]interface{}                                      // 观察适配器运行状态
}

// 出于扩展的考虑, init里面只进行一些无状态的操作, 和外部有交互的放在方法里面
// 没错, 就是说的  flag.Parse() 这个玩意
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
