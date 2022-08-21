package ssh

import "errors"

var (
	initFailedError  = errors.New("ssh init failed")
	initTimeoutError = errors.New("ssh init timeout")
	dialTimeoutError = errors.New("ssh dial timeout")
)
