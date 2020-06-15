package service

import (
	"fmt"
	"testing"
)

func TestUserBind(t *testing.T) {
	model := NewUserModel("id_test", "测试名称")

	model.BindConnect("connect1")
	model.BindConnect("connect2")
	model.BindConnect("connect3")
	model.UnBindConnect("connect2")

	fmt.Printf("userModel: %v \n", model)
}

func TestUserLoad(t *testing.T) {
	model := UserModel{}

	json := []byte(`{
		"id": "test_id",
		"name": "test_name"
	}`)

	model.Load(json)
	model.BindConnect("connect2")
	fmt.Printf("userModel: %#+v \n", model.Name)
}

func TestUserList(t *testing.T) {

	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("test_id_%d", i)
		user := &UserModel{
			ID:   id,
			Name: fmt.Sprintf("test_name_%d", i),
		}

		err := UserList.AddUser(user)
		if err != nil {
			t.Error(err)
		}

	}
	fmt.Printf("UserList: %#+v \n", UserList.IDList())

	user, err := UserList.GetUser("test_id_0")
	fmt.Printf("User: %v, err: %v \n", user, err)
}
