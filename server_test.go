package chat

import (
	"sync"
	"testing"
)

var testWg *sync.WaitGroup

func TestServer(t *testing.T) {
	testWg = &sync.WaitGroup{}

	testWg.Add(1)

	go ServerStar()

	t.Log("测试开始")

	template1(t)

	testWg.Wait()
	t.Log("测试结束")
}

func template1(t *testing.T) {
	t.Log("执行测试用例1")

	t.Log("开始结束测试")
	testWg.Done()
}
