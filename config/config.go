package config

import "flag"

// ADDR 服务启动地址
var ADDR = flag.String("addr", "localhost:8080", "聊天室地址,eg  localhost:8080")

func init() {
	flag.Parse()
}
