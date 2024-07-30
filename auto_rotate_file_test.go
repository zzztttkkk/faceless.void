package fv_test

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	fv "github.com/zzztttkkk/faceless.void"
)

func TestAutoRotateFileBySize(t *testing.T) {
	arf := fv.NewAutoRotateFile("./test.log", &fv.RotateOptions{Kind: fv.RotateKindFileSize, MaxSize: 128})
	defer arf.Close()

	for i := 0; i < 10; i++ {
		_, err := arf.Write([]byte(strings.Repeat("a", 256)))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAutoRotateFileByTime(t *testing.T) {
	arf := fv.NewAutoRotateFile("./test.log", &fv.RotateOptions{Kind: fv.RotateKindMinutely})
	defer arf.Close()

	logger := slog.New(slog.NewJSONHandler(arf, nil))
	for i := 0; i < 300; i++ {
		logger.Info("test", slog.Int("i", i))
		fmt.Println("---->", i)
		time.Sleep(time.Second)
	}
}
