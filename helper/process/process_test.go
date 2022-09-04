package process

import (
	"fmt"
	"testing"
)

func TestIsDone(t *testing.T) {
	done, err := IsRunning(4912)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(done)
}
