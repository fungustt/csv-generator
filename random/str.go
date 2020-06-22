package random

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

type StrRandomizer struct {
	ch     chan string
	stop   chan struct{}
	rand   *rand.Rand
	length int
}

func NewStrRandomizer(length int) *StrRandomizer {
	r := &StrRandomizer{
		ch:     make(chan string),
		stop:   make(chan struct{}),
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
		length: length,
	}

	go r.generate()
	return r
}

func (r *StrRandomizer) Get() string {
	return <-r.ch
}

func (r *StrRandomizer) Stop() {
	r.stop <- struct{}{}
}

func (r *StrRandomizer) generate() {
	b := make([]byte, r.length)
	for {
		for i := range b {
			b[i] = charset[r.rand.Intn(len(charset))]
		}

		select {
		case <-r.stop:
			break
		case r.ch <- string(b):
		}
	}
}
