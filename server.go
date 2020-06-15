package chat

import (
	"chat/config"
	"chat/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

// ServerStar 启动
func ServerStar() {
	config.Parse()
	gentleExit()
	service.Start()
	mux := Routes()
	log.Fatal(http.ListenAndServe(*config.ADDR, mux))
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
	dispatcher := service.GetDispatcher()
	dispatcher.Close()
	fmt.Println("结束退出...")
	os.Exit(0)
}

// SocketStatus 状态查询
func SocketStatus(w http.ResponseWriter, r *http.Request) {
	dispatcher := service.GetDispatcher()
	statusMap := dispatcher.Status()
	b, _ := json.Marshal(statusMap)
	w.Write(b)
}

// SocketServer websocket服务
func SocketServer(w http.ResponseWriter, r *http.Request) {

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
	client, err := service.NewSocketClient(id, w, r)
	defer client.Close()

	if err != nil {
		errMsg := []byte(`发生链接错误,错误原因` + err.Error())
		log.Println(errMsg)
		w.Write(errMsg)
		return
	}

	client.SysBroadcast(fmt.Sprintf("欢迎 %s", id))

	for {
		_, message, err := client.ReadMessage()
		// log.Printf("read in server:  %s  err: %v \n", message, err)
		if websocket.IsCloseError(err, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
			log.Println("主动断开链接")
			return
		}
		if err != nil {
			log.Println("error:", err)
			return
		}

		var msgBody service.Message
		err = json.Unmarshal(message, &msgBody)

		if err != nil {
			return
		}

		// 将收到的消息广播给所有链接
		err = client.DispatchRequest(msgBody)
		if err != nil {
			errMsg := []byte(`发生链接错误,错误原因` + err.Error())
			log.Printf("%s \n", errMsg)
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
