package syncUtils

import "sync/atomic"

type AtomicCounter uint64

// Adds 1 and returns the new value
func (c *AtomicCounter) Add() uint64 {
	return atomic.AddUint64((*uint64)(c), 1)
}
