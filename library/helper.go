package library

import (
	"hash/crc32"
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

// HashStr2Int 将字符串转换成一串唯一编码
func HashStr2Int(s string) (code int) {
	code = int(crc32.ChecksumIEEE([]byte(s)))
	if code >= 0 {
		return
	}
	if -code >= 0 {
		return -code
	}
	// v == MinInt
	return
}
