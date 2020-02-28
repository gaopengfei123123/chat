package main

import (
	"bufio"
	"chat/library"
	"chat/service"
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

var host = flag.String("host", "0.0.0.0:8080", "http service address")

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	start()
}

var inputReader *bufio.Reader

func init() {
	inputReader = bufio.NewReader(os.Stdin)
}

func start() {
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *host, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// 监听输出
	go output(conn)

	// 循环接收客户端输入
	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			log.Panic(err)
		}

		input = library.Trim(input)

		if input == "" {
			continue
		}

		// log.Printf("input: %#+v \n", input)
		msg := buildMsg(input)
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

// 监听输出
func output(conn *websocket.Conn) {
	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			if msgType == -1 {
				log.Println("服务端断开链接")
				os.Exit(0)
			}
			log.Println("read:", err)
			log.Println("msgType", msgType)
			return
		}
		var msg service.Message
		json.Unmarshal(message, &msg)
		// log.Printf("\n \n 接收到消息: %#+v \n \n", msg)

		switch msg.Type {
		case service.SystemMessage:
			log.Printf("系统广播 %v\n", msg.Content)
		case service.BroadcastMessage:
			log.Printf("%s 说: %s \n", msg.From, msg.Content)
		case service.ConnectedMessage:
			log.Printf("%s  上线\n", msg.From)
		case service.DisconnectedMessage:
			log.Printf("%s  下线\n", msg.From)
		case service.HeartBeatMessage:
			msg := buildMsg("", service.HeartBeatMessage)
			log.Println("heartBeat ack")
			conn.WriteMessage(websocket.TextMessage, msg)
		default:
			log.Printf("msyType: %d, content: %v\n", msg.Type, msg.Content)
		}
	}
}

func buildMsg(msg string, Type ...int) []byte {
	m := service.Message{
		Content: msg,
		Type:    service.BroadcastMessage,
	}

	if len(Type) > 0 {
		m.Type = Type[0]
	}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return b
}
