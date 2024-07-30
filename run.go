package fv

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

type OnShutdownFunc func(wg *sync.WaitGroup)

var (
	_hooks_lock   = sync.Mutex{}
	_hooks        = make([]OnShutdownFunc, 0)
	_once         = atomic.Bool{}
	_app_finished = atomic.Bool{}
)

func RegisterShutdownHook(fnc OnShutdownFunc) {
	if _app_finished.Load() {
		panic("process already finished, can not register shutdown hook anymore")
	}

	_hooks_lock.Lock()
	defer _hooks_lock.Unlock()
	_hooks = append(_hooks, fnc)
}

func Run(launch func(ctx context.Context)) {
	if !_once.CompareAndSwap(false, true) {
		panic("Run already called")
	}

	// https://github.com/ElyDotDev/windows-kill
	// Send signal to process by PID in Windows, like POSIX kill
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer cancel()
		launch(ctx)
	}()
	<-ctx.Done()
	_app_finished.Store(true)

	_hooks_lock.Lock()
	defer _hooks_lock.Unlock()

	wg := &sync.WaitGroup{}
	wg.Add(len(_hooks))
	for _, fnc := range _hooks {
		go fnc(wg)
	}
	wg.Wait()
}
