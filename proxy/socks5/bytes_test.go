package socks5

import (
	"fmt"
	"testing"
)

func TestIntToBytes(t *testing.T) {
	ports := []int{443, 80}
	for _, port := range ports {
		bytes, err := IntToBytes(port, 2)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(bytes, ChangeStrPort(bytes))
	}
	fmt.Println("successful")
}
