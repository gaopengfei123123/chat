package library

import (
	"math/rand"
	"strings"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq 随机字符串
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Trim 移除两端空格和制表符
func Trim(str string) string {
	str = strings.Trim(str, "\n")
	return strings.Trim(str, " ")
}
