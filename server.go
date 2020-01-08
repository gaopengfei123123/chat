package chat

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var addr = flag.String("addr", "localhost:8080", "聊天室地址,eg  localhost:8080")

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

// ServerStar 启动
func ServerStar() {
	flag.Parse()

	gentleExit()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello`))
	})
	http.HandleFunc("/ws", socketServer)

	log.Printf("监听端口: %v", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// 监听关闭信号
func gentleExit() {
	// 创建监听退出信号的chan
	c := make(chan os.Signal)
	// 监听
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("退出", s)
				ExitFunc()
			case syscall.SIGUSR1:
				fmt.Println("usr1", s)
			case syscall.SIGUSR2:
				fmt.Println("usr2", s)
			default:
				fmt.Println("other", s)
			}
		}
	}()

}

// ExitFunc 退出函数
func ExitFunc() {
	fmt.Println("开始退出...")
	fmt.Println("执行清理...")
	handler := GetHandler()
	handler.Close()
	fmt.Println("结束退出...")
	os.Exit(0)
}

func socketServer(w http.ResponseWriter, r *http.Request) {

	if websocket.IsWebSocketUpgrade(r) {
		log.Println("收到websocket链接")
	} else {
		log.Println("您这也不是websocket啊")
		w.Write([]byte(`您这也不是websocket啊`))
		return
	}

	// 使用Sec-WebSocket-Key当链接key
	id := r.Header.Get("Sec-WebSocket-Key")
	log.Printf("header: Sec-WebSocket-Key is \" %v \" \n", id)
	client, err := NewSocketClient(id, w, r)
	defer client.Close()

	if err != nil {
		errMsg := []byte(`发生链接错误,错误原因` + err.Error())
		log.Println(errMsg)
		w.Write(errMsg)
		return
	}

	welcome2 := fmt.Sprintf("欢迎 %s", id)
	client.SendMessage(1, welcome2)

	for {
		_, message, err := client.conn.ReadMessage()
		// log.Printf("read in server:  %s  err: %v \n", message, err)
		if websocket.IsCloseError(err, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
			log.Println("主动断开链接")
			return
		}
		if err != nil {
			log.Println("error:", err)
			return
		}

		client.Broadcast(string(message))
		err = HandleRequest(client, message)
		if err != nil {
			errMsg := []byte(`发生链接错误,错误原因` + err.Error())
			log.Println(errMsg)
			w.Write(errMsg)
			return
		}
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
