运行示例:
先开启server

```bash
$ go run server.go
```

再运行client
```bash
$ go run client.go
```

server的默认链接地址是 `0.0.0.0:8080`, 因此,  如果`client`和`server`没有在同一台机器上, 则需要标明`server`的ip地址 运行:
```bash
./client --host=192.168.132.71:8080
```

可以多开几个client, **直接在命令行里输入 回车发送就行**