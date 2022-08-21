package common

import "sync"

func CreateSyncPool(f func() any) sync.Pool {
	return sync.Pool{
		New: f,
	}
}

func CreateByteBufferSyncPool(bufferSize int) sync.Pool {
	return CreateSyncPool(func() any {
		return make([]byte, bufferSize)
	})
}
