package main

import (
	"bufio"
	"chat"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
)

var host = flag.String("host", "localhost:8080", "http service address")

func main() {
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

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var msg chat.Message
			json.Unmarshal(message, &msg)
			switch msg.Type {
			case chat.SystemMessage:
				log.Printf("系统广播 %v\n", msg)
			case chat.BroadcastMessage:
				log.Printf("%s 说: %s \n", msg.ID, msg.Content)
			case chat.ConnectedMessage:
				log.Printf("%s  上线\n", msg.ID)
			case chat.DisconnectedMessage:
				log.Printf("%s  下线\n", msg.ID)
			}

		}
	}()

	// 循环接收客户端输入
	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			log.Panic(err)
		}
		conn.WriteMessage(websocket.TextMessage, []byte(input))
	}
}
