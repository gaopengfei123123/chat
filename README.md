### 一个简单的聊天室


#### 依赖文件

```
go get github.com/gorilla/websocket
```

运行: 

main.go
```go
package main

import (
	"chat"
)

func main() {
	chat.ServerStar()
}
```

```bash
$ go run main.go  默认地址 localhost:8080
//或
$ go run main.go -addr 127.0.0.1:8080
```
监听 `ws://localhost:8080/ws`

#### 一个具体例子

```go
// 启动服务端
go run example/server.go

// 启动客户端
go run example/client.go
```
客户端之间消息互相广播


##### TODO

1. 多加几个事件钩子
2. id与conn绑定
3. 建立群组
4. 指定范围的广播
5. socket转http
6. 历史消息/限流
