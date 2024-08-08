package main

import (
	"context"
	"exec/internal"
	"fmt"
	"os"

	fv "github.com/zzztttkkk/faceless.void"
	"sync"
)

func main() {
	fmt.Printf("pid: %d\r\n", os.Getpid())

	fv.RegisterShutdownHook(func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("shutdown")
	})

	fv.Run(func(ctx context.Context) {
		internal.DIC.Run()
		<-ctx.Done()
	})
}
