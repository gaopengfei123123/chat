package service

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
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

// Client socket客户端
type Client struct {
	ID         string          // 链接的唯一标识
	conn       *websocket.Conn // 链接实体
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

// 跨域配置
var upgrader = websocket.Upgrader{}

func init() {
	// 允许跨域请求
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}
