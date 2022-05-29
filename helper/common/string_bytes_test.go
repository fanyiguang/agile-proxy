package common

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
		fmt.Println(bytes, changeInt(bytes))
	}
	fmt.Println("successful")
}

func changeInt(bPort []byte) int {
	if len(bPort) < 2 {
		return 0
	}
	return int(bPort[0])<<8 | int(bPort[1])
}
