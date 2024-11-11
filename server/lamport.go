package main

import (
	"sync/atomic"
)

type LamportClock struct {
	time atomic.Uint64
}

func NewLamportClock() *LamportClock {
	return &LamportClock{
		time: atomic.Uint64{},
	}
}

func (c *LamportClock) Now() uint64 {
	return c.time.Load()
}

func (c *LamportClock) Step() {
	c.time.Add(1)
}

func (c *LamportClock) Max(other uint64) {
	if c.time.Load() < other {
		c.time.Swap(other)
	}
}
