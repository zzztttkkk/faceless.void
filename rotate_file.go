package fv

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"

	"github.com/zzztttkkk/faceless.void/internal"
)

type RotateKind int

const (
	RotateKindDaily = RotateKind(iota)
	RotateKindHourly
	RotateKindMinutely
	RotateKindFileSize
)

type _RotateFileBuilder struct {
	pairs []internal.Pair[string]
}

func (builder *_RotateFileBuilder) Filepath(fp string) *_RotateFileBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("filepath", fp))
	return builder
}

func (builder *_RotateFileBuilder) Kind(kind RotateKind) *_RotateFileBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("kind", kind))
	return builder
}

func (builder *_RotateFileBuilder) MaxSize(maxsize int64) *_RotateFileBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("maxsize", maxsize))
	return builder
}

func (builder *_RotateFileBuilder) BufferSize(bufsize int) *_RotateFileBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("bufsize", bufsize))
	return builder
}

func (builder *_RotateFileBuilder) OnRotateError(fnc func(error)) *_RotateFileBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("onerror", fnc))
	return builder
}

type _AutoRotateFile struct {
	lock        sync.Mutex
	fp          string
	fins        *os.File
	w           *bufio.Writer
	currentSize int64

	kind          RotateKind
	maxSize       int64
	bufferSize    int
	onRotateError func(error)

	closed bool

	fname string
	fext  string
	timer *time.Timer
}

func (builder *_RotateFileBuilder) Build() io.WriteCloser {
	ins := &_AutoRotateFile{}

	for _, pair := range builder.pairs {
		switch pair.Key {
		case "filepath":
			{
				fp := pair.Val.(string)
				ext := path.Ext(fp)
				ins.fp = fp
				ins.fext = ext
				ins.fname = fp[0 : len(fp)-len(ext)]
				break
			}
		case "kind":
			{
				ins.kind = pair.Val.(RotateKind)
				break
			}
		case "bufsize":
			{
				ins.bufferSize = pair.Val.(int)
				break
			}
		case "maxsize":
			{
				ins.maxSize = pair.Val.(int64)
				break
			}
		case "onerror":
			{
				ins.onRotateError = pair.Val.(func(error))
				break
			}
		}
	}

	if ins.fp == "" {
		panic(errors.New("empty filepath"))
	}

	if ins.bufferSize < 4096 {
		ins.bufferSize = 4096
	}
	if ins.kind == RotateKindFileSize && ins.maxSize < 1 {
		ins.maxSize = 1024 * 1024 * 64 // 64MB
	}
	if ins.kind != RotateKindFileSize {
		ins.rotateByTime()
	}
	return ins
}

func RotateFileBuilder() *_RotateFileBuilder {
	return &_RotateFileBuilder{}
}

func (arf *_AutoRotateFile) doRotate(newname string) {
	err := arf.doClose()
	if err != nil {
		err = fmt.Errorf("rotate error: when flush, %s", err)
		if arf.onRotateError != nil {
			arf.onRotateError(err)
		} else {
			fmt.Println(err)
		}
	}
	err = os.Rename(arf.fp, newname)
	if err != nil {
		err = fmt.Errorf("rotate error: when rename, %s", err)
		if arf.onRotateError != nil {
			arf.onRotateError(err)
		} else {
			fmt.Println(err)
		}
	}
}

func (arf *_AutoRotateFile) rotateByTime() {
	now := time.Now()
	var end time.Time
	var timepart string
	switch arf.kind {
	case RotateKindDaily:
		{
			end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999, time.Local)
			timepart = end.Format("20060102")
			break
		}
	case RotateKindHourly:
		{
			end = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 59, 59, 999, time.Local)
			timepart = end.Format("2006010215")
			break
		}
	case RotateKindMinutely:
		{
			end = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 59, 999, time.Local)
			timepart = end.Format("200601021504")
			break
		}
	default:
		{
			panic("unreachable code")
		}
	}

	diff := end.Sub(now) + time.Second
	arf.timer = time.AfterFunc(diff, func() {
		arf.lock.Lock()
		defer arf.lock.Unlock()

		newname := fmt.Sprintf(`%s.%s%s`, arf.fname, timepart, arf.fext)
		if exists, _ := internal.FsExists(newname); exists {
			newname = fmt.Sprintf(`%s.%s.%d%s`, arf.fname, timepart, time.Now().Nanosecond(), arf.fext)
		}
		arf.doRotate(newname)
		arf.rotateByTime()
	})
}

func (arf *_AutoRotateFile) doClose() error {
	if arf.fins == nil {
		return nil
	}
	err := arf.w.Flush()
	if err != nil {
		return err
	}
	err = arf.fins.Close()
	if err != nil {
		return err
	}
	arf.fins = nil
	return nil
}

func (arf *_AutoRotateFile) Close() error {
	arf.lock.Lock()
	defer func() {
		arf.closed = true
		arf.lock.Unlock()
	}()
	if arf.timer != nil {
		arf.timer.Stop()
	}
	return arf.doClose()
}

var (
	ErrFileIsClosed = errors.New("file is closed")
)

func (arf *_AutoRotateFile) Write(p []byte) (n int, err error) {
	arf.lock.Lock()
	defer arf.lock.Unlock()

	if arf.closed {
		return 0, ErrFileIsClosed
	}

	if arf.fins == nil {
		f, err := os.OpenFile(arf.fp, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return 0, err
		}
		stat, err := f.Stat()
		if err != nil {
			return 0, err
		}
		arf.currentSize = stat.Size()
		arf.fins = f
		if arf.w == nil {
			arf.w = bufio.NewWriterSize(arf.fins, arf.bufferSize)
		} else {
			arf.w.Reset(arf.fins)
		}
	}

	n, err = arf.w.Write(p)
	if err != nil {
		return n, err
	}

	if arf.kind == RotateKindFileSize {
		arf.currentSize += int64(len(p))
		if arf.currentSize >= arf.maxSize {
			arf.doRotate(fmt.Sprintf(`%s.%s%s`, arf.fname, time.Now().Format("20060102150405.000"), arf.fext))
		}
	}
	return n, nil
}

var _ io.WriteCloser = (*_AutoRotateFile)(nil)
