package chat

import (
	"net/http"
)

// Routes 对外暴露的路由接口
func Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello`))
	})
	// websocket 监听接口
	mux.HandleFunc("/ws", SocketServer)
	// 服务状态接口, 观察一下消息积攒情况
	mux.HandleFunc("/status", SocketStatus)
	return mux
}
