package main

import (
	"context"
	"fmt"
	"os"

	"sync"
	"time"

	fv "github.com/zzztttkkk/faceless.void"
)

func main() {
	fmt.Printf("pid: %d\r\n", os.Getpid())

	fv.RegisterShutdownHook(func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("shutdown")
	})

	fv.Run(func(ctx context.Context) {
		time.AfterFunc(time.Second*50, func() {
			fmt.Println("done")
		})
		<-ctx.Done()
	})
}
