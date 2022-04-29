package utils

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	s := Md5("000000000")
	fmt.Println(s)
}

func TestPasswd(t *testing.T) {
	s := Password("000000", "123455")
	fmt.Println(s)
}
