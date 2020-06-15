// Package service 这里写注册逻辑 业务id和链接id绑定
package service

import (
	"chat/library"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// UserList 用户列表
var UserList *UList

// UList 当前连线用户列表
type UList struct {
	list      map[int]*UserModel // 存储用户状态
	ActiveCnt int                // 当前连接数
	sync.RWMutex
}

// UserModel 用户模型
type UserModel struct {
	ID          string   // 用户唯一id
	Name        string   // 随便一个名称
	connectList []string // 链接id组, 可能存在多点登录的情况, 因为是运行环境中生成的id, 以私有属性存储
	groupList   []string // 所在分组id
	sync.Mutex
}

func init() {
	// server端缓存一些用户状态, 以及在线连接信息, 目前还是站在单机的角度考虑
	UserList = NewUserList(1024)
}

// NewUserList 新建连接列表
func NewUserList(size int) *UList {
	UserList := new(UList)
	UserList.Init(size)
	return UserList
}

// GetUserList 获取列表
func GetUserList() *UList {
	return UserList
}

// Init 初始化
func (l *UList) Init(size int) {
	if size <= 0 {
		size = 1024
	}
	l.list = make(map[int]*UserModel, size)
}

// AddUser 将用户加入
func (l *UList) AddUser(user *UserModel) error {
	if user.ID == "" {
		return errors.New("需要用户ID")
	}
	k := library.HashStr2Int(user.ID)
	l.Lock()
	defer l.Unlock()
	l.list[k] = user
	l.ActiveCnt = len(l.list)
	return nil
}

// GetUser 获取指定用户信息
func (l *UList) GetUser(userID string) (user *UserModel, err error) {
	l.RLock()
	defer l.RUnlock()

	k := library.HashStr2Int(userID)
	user, ok := l.list[k]
	if !ok {
		err = errors.New("用户ID不存在")
	}
	return
}

// DelUser 移除用户
func (l *UList) DelUser(userID string) (user *UserModel, err error) {
	l.Lock()
	defer l.Unlock()

	k := library.HashStr2Int(userID)
	user, ok := l.list[k]
	if !ok {
		err = errors.New("用户ID不存在")
	}

	delete(l.list, k)
	l.ActiveCnt--

	return
}

// IDList 获取 ID list
func (l *UList) IDList() []string {
	res := make([]string, 0, l.ActiveCnt)
	for _, v := range l.list {
		res = append(res, v.ID)
	}
	return res
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
