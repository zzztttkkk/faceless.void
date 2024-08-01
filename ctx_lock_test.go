package fv_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	fv "github.com/zzztttkkk/faceless.void"
)

func TestLock(t *testing.T) {
	var lock = &sync.RWMutex{}
	lock.Lock()

	ctx, cancel := context.WithTimeout(fv.CtxValForLockAcquireSleepStep(context.Background(), time.Millisecond*200), time.Second)
	defer cancel()

	err := fv.AcquireLock(ctx, lock)
	if err != nil {
		fmt.Println(err)
	}
}
