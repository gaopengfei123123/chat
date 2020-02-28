客户端发送定向消息:

to字段, 其中是cli.ID的数组
```json
{"id":"MuE/CmLxOy4lLCx20xOdfA==","content":"测试定向消息5","type":7,"to":["zzKqVQCKp5bt7/x1bd6zxA=="]}
```

心跳检测消息:
```json
{"ID":"","Content":"","sent_at":1582885834,"Type":2,"From":"","To":null,"from_user_id":"","to_user_id":"","ext":null}
```
server端发送这条消息, 三次没回应则会中断链接,  期间client任何一次回复, 重试计数归零