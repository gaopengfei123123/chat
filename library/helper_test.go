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
