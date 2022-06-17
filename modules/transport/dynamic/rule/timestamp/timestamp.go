package timestamp

import (
	"math/rand"
	"time"
)

type timestamp struct {
	rand *rand.Rand
}

func (t *timestamp) Int() int {
	return t.rand.Int()
}

func (t *timestamp) Intn(n int) int {
	return t.rand.Intn(n)
}

func New() (obj *timestamp, err error) {
	obj = &timestamp{rand: rand.New(rand.NewSource(time.Now().Unix()))}
	return
}
