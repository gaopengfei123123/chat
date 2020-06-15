package main

import (
	"bufio"
	"chat/library"
	"chat/service"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

var host = flag.String("host", "0.0.0.0:8080", "http service address")
var user userModel // 临时用户名

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)

	user = tmpUser()
	start()
}

var inputReader *bufio.Reader

func init() {
	inputReader = bufio.NewReader(os.Stdin)
}

// 临时的用户结构
type userModel struct {
	Name string
}

// 简单的写一个临时账号的玩意
func tmpUser() userModel {
	fmt.Println("请输入一个名字:")

	input, err := inputReader.ReadString('\n')

	if err != nil {
		fmt.Printf("An Error occurred:%s\n", err)
		os.Exit(1)
	}
	name := input[:len(input)-1]
	fmt.Printf("Hello %s! 输入消息后按空格发送\n", name)

	u := userModel{
		Name: name,
	}
	return u
}

// 启动客户端
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

	// 注册账户名
	m := registerMsg(user.Name)
	conn.WriteMessage(websocket.TextMessage, m)

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
		log.Printf("\n \n 接收到消息: %#+v \n \n", msg)

		var n interface{}
		if msg.Ext == nil {
			n = "匿名"
		} else {
			n = msg.Ext["name"]
		}
		name := n.(string)

		switch msg.Type {
		case service.SystemMessage:
			log.Printf("系统广播 %v\n", msg.Content)
		case service.BroadcastMessage:
			log.Printf("%s 说: %s \n", name, msg.Content)
		case service.ConnectedMessage:
			log.Printf("%s  上线\n", name)
		case service.DisconnectedMessage:
			log.Printf("%s  下线\n", name)
		case service.HeartBeatMessage:
			msg := buildMsg("", service.HeartBeatMessage)
			// log.Println("heartBeat ack")
			conn.WriteMessage(websocket.TextMessage, msg)
		default:
			log.Printf("msyType: %d, content: %v\n", msg.Type, msg.Content)
		}
	}
}

// 构建注册类消息(临时注册)
func registerMsg(name string) []byte {
	m := service.Message{
		Content: "",
		Type:    service.RegisterMessage,
		Ext: map[string]interface{}{
			"tmp":  "1",
			"name": name,
		},
	}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return b
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
