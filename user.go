// Package chat 这里写注册逻辑 业务id和链接id绑定
package chat

import (
	"encoding/json"
	"fmt"
	"sync"
)

// 用户列表暂时放在缓存里面
var modelList map[string]*UserModel

// UserModel 用户模型
type UserModel struct {
	ID          string   // 用户唯一id
	Name        string   // 随便一个名称
	connectList []string // 链接id组, 可能存在多点登录的情况, 因为是运行环境中生成的id, 以私有属性存储
	groupList   []string // 所在分组id
	sync.Mutex
}

func init() {
	modelList = make(map[string]*UserModel, 1)
}

// NewUserModel 初始化用户模型
func NewUserModel(id, name string) *UserModel {
	return &UserModel{
		ID:          id,
		Name:        name,
		connectList: make([]string, 1),
		groupList:   make([]string, 1),
	}
}

// Load 将json转换成实例
func (user *UserModel) Load(jsonData []byte) error {
	return json.Unmarshal(jsonData, user)
}

// Serialize 序列化
func (user *UserModel) Serialize() ([]byte, error) {
	return json.Marshal(user)
}

// BindConnect 用户绑定connectID
func (user *UserModel) BindConnect(connectID string) error {
	user.Lock()
	defer user.Unlock()
	for i := 0; i < len(user.connectList); i++ {
		if user.connectList[i] == connectID {
			return fmt.Errorf("重复的connectID")
		}
	}

	user.connectList = append(user.connectList, connectID)

	return nil
}

// UnBindConnect 解绑链接
func (user *UserModel) UnBindConnect(connectID string) error {
	user.Lock()
	defer user.Unlock()

	for i := 0; i < len(user.connectList); i++ {
		if user.connectList[i] == connectID {
			user.connectList = append(user.connectList[0:i], user.connectList[i+1:len(user.connectList)]...)
			return nil
		}
	}
	return fmt.Errorf("指定ID不存在")
}
