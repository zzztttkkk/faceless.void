package fv_test

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	fv "github.com/zzztttkkk/faceless.void"
)

func TestRotateFileBySize(t *testing.T) {
	rf := fv.RotateFileBuilder().Filepath("./test.rotate_by_size.log").Kind(fv.RotateKindFileSize).MaxSize(1024).Build()
	defer rf.Close()

	for i := 0; i < 10; i++ {
		_, err := rf.Write([]byte(strings.Repeat("a", 256)))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestRotateFileByTime(t *testing.T) {
	rf := fv.RotateFileBuilder().Filepath("./test.rotate_by_time.log").Kind(fv.RotateKindMinutely).Build()
	defer rf.Close()

	logger := slog.New(slog.NewJSONHandler(rf, nil))
	for i := 0; i < 300; i++ {
		logger.Info("test", slog.Int("i", i))
		fmt.Println("---->", i)
		time.Sleep(time.Second)
	}
}
