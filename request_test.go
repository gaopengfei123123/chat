package chat

import (
	"fmt"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	sendID := "send_id_123"
	acceptID := "accept_id_456"

	cli := &Client{}
	var msg []byte
	var err error

	msg = []byte(fmt.Sprintf(`{
		"id": "%s",
		"content": "系统消息内容",
		"type": 0,
		"to": ["%s"]
	}`, sendID, acceptID))

	err = HandleRequest(cli, msg)

	if err != nil {
		t.Error(err)
	}

	msg = []byte(`{
		"id": "send_id",
		"content": "广播消息内容",
		"type": 1,
		"to": ["accept_id"]
	}`)

	err = HandleRequest(cli, msg)

	if err != nil {
		t.Error(err)
	}

	msg = []byte(`{
		"id": "send_id",
		"content": "心跳消息内容",
		"type": 2,
		"to": ["accept_id"]
	}`)

	err = HandleRequest(cli, msg)

	if err != nil {
		t.Error(err)
	}

	msg = []byte(`{
		"id": "send_id",
		"content": "上线通知消息内容",
		"type": 3,
		"to": ["accept_id"]
	}`)

	err = HandleRequest(cli, msg)

	if err != nil {
		t.Error(err)
	}

	msg = []byte(fmt.Sprintf(`{
		"id": "%s",
		"content": "下线通知消息内容",
		"type": 4,
		"to": ["%s"]
	}`, sendID, acceptID))

	err = HandleRequest(cli, msg)

	if err != nil {
		t.Error(err)
	}

	msg = []byte(fmt.Sprintf(`{
		"id": "%s",
		"content": "服务断开链接通知",
		"type": 5,
		"to": ["%s"]
	}`, sendID, acceptID))

	err = HandleRequest(cli, msg)

	if err != nil {
		t.Error(err)
	}

	msg = []byte(fmt.Sprintf(`{
		"id": "%s",
		"content": "注册事件消息",
		"type": 6,
		"to": ["%s"]
	}`, sendID, acceptID))

	err = HandleRequest(cli, msg)

	if err != nil {
		t.Error(err)
	}
}
