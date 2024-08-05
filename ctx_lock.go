package fv

import (
	"context"
	"sync"
	"time"
)

type tryLocker interface {
	TryLock() bool
}

var _ tryLocker = (*sync.Mutex)(nil)

func CtxValForLockAcquireSleepStep(ctx context.Context, duration time.Duration) context.Context {
	return context.WithValue(ctx, _ctx_key_for_lock_sleep_step, duration)
}

func getLockSleepStep(ctx context.Context) time.Duration {
	tmp := ctx.Value(_ctx_key_for_lock_sleep_step)
	var duration = time.Millisecond * 30
	if tmp != nil {
		v, ok := tmp.(time.Duration)
		if ok {
			duration = v
		}
	}
	return duration
}

func AcquireLock(ctx context.Context, v tryLocker) error {
	duration := getLockSleepStep(ctx)
	for {
		if v.TryLock() {
			return nil
		}
		select {
		case <-ctx.Done():
			{
				return ctx.Err()
			}
		default:
			{
				time.Sleep(duration)
			}
		}
	}
}

type tryRLocker interface {
	TryRLock() bool
}

var _ tryRLocker = (*sync.RWMutex)(nil)

func AcquireRLock(ctx context.Context, v tryRLocker) error {
	duration := getLockSleepStep(ctx)
	for {
		if v.TryRLock() {
			return nil
		}
		select {
		case <-ctx.Done():
			{
				return ctx.Err()
			}
		default:
			{
				time.Sleep(duration)
			}
		}
	}
}
