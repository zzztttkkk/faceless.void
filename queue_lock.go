package fv

import (
	"context"
	"errors"
	"sync"
)

type QueueLock struct {
	lock       sync.RWMutex
	busy       bool
	waiters    []chan bool
	maxwaiters int
}

func NewQueueLockWithLimit(maxwaiters int) *QueueLock {
	return &QueueLock{
		maxwaiters: maxwaiters,
	}
}

var (
	ErrQueueLockFull      = errors.New("too many waiters")
	ErrQueueLockIsNotBusy = errors.New("lock is not busy, can not release")
	ErrQueueLockClosed    = errors.New("queue lock is closed")
)

func (ql *QueueLock) SetMaxWaiters(maxwaiters int) {
	ql.lock.Lock()
	defer ql.lock.Unlock()
	ql.maxwaiters = maxwaiters
}

func (ql *QueueLock) Acquire(ctx context.Context) error {
	ql.lock.Lock()

	if !ql.busy {
		ql.busy = true
		ql.lock.Unlock()
		return nil
	}

	if ql.maxwaiters > 0 && len(ql.waiters) >= ql.maxwaiters {
		ql.lock.Unlock()
		return ErrQueueLockFull
	}

	ch := make(chan bool, 1)
	ql.waiters = append(ql.waiters, ch)
	ql.lock.Unlock()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ok := <-ch:
			if ok {
				return nil
			}
			return ErrQueueLockClosed
		}
	}
}

func (ql *QueueLock) Release() {
	ql.lock.Lock()
	defer ql.lock.Unlock()

	if !ql.busy {
		panic(ErrQueueLockIsNotBusy)
	}

	if len(ql.waiters) > 0 {
		ch := ql.waiters[0]
		ql.waiters = ql.waiters[1:]
		ch <- true
	} else {
		ql.busy = false
	}
}

func (ql *QueueLock) Close() {
	ql.lock.Lock()
	defer ql.lock.Unlock()

	for _, ch := range ql.waiters {
		ch <- false
	}
	ql.waiters = nil
}

func (ql *QueueLock) Status() (busy bool, waits int) {
	ql.lock.RLock()
	defer ql.lock.RUnlock()
	return ql.busy, len(ql.waiters)
}

type QueueLockGroup[K comparable] struct {
	lock       sync.RWMutex
	queues     map[K]*QueueLock
	maxwaiters int
}

func (qlg *QueueLockGroup[K]) Get(key K) *QueueLock {
	qlg.lock.RLock()
	ptr := qlg.queues[key]
	if ptr != nil {
		qlg.lock.RUnlock()
		return ptr
	}
	qlg.lock.RUnlock()

	qlg.lock.Lock()
	defer qlg.lock.RLock()

	ptr = NewQueueLockWithLimit(qlg.maxwaiters)
	qlg.queues[key] = ptr
	return ptr
}

func (qlg *QueueLockGroup[K]) Remove(key K) {
	qlg.lock.Lock()

	ptr := qlg.queues[key]
	if ptr != nil {
		qlg.lock.Unlock()
		return
	}
	delete(qlg.queues, key)
	qlg.lock.Unlock()

	ptr.Close()
}
