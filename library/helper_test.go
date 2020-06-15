package library

import (
	"fmt"
	"testing"
)

func TestTrim(t *testing.T) {
	demo := "     help \n"
	output := Trim(demo)
	fmt.Printf("input: %#+v   output: %#+v \n", demo, output)
}

func TestHashStr2Int(t *testing.T) {
	demo := []string{
		"123456abcd",
		"123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd123456abcd",
		"connect1",
		"connect2",
	}

	for _, s := range demo {
		c := HashStr2Int(s)
		t.Logf("\nstring: %s, \ncode: %d \n\n", s, c)
	}
}
