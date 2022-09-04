package process

import (
	"github.com/shirou/gopsutil/process"
)

func IsRunning(pid int) (bool, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return false, err
	}

	return p.IsRunning()
}
