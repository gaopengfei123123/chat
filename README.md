### 一个简单的聊天室



#### 依赖文件:
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