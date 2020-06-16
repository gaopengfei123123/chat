package config

import (
	"flag"
	"time"
)

// ADDR 服务启动地址
var ADDR = flag.String("addr", "0.0.0.0:8080", "聊天室地址,eg  0.0.0.0:8080")

// HeartbeatInterval 心跳检测间隔
var HeartbeatInterval = 30 * time.Second

// MaxRetryTime 心跳重试次数
var MaxRetryTime = 2

// Parse 解析命令行参数
func Parse() {
	flag.Parse()
}
