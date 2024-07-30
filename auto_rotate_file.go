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
	RotateKindDaily RotateKind = RotateKind(iota)
	RotateKindHourly
	RotateKindMinutely
	RotateKindFileSize
)

type RotateOptions struct {
	Kind          RotateKind
	MaxSize       int64
	BufferSize    int
	OnRotateError func(error)
}

type _AutoRotateFile struct {
	lock        sync.Mutex
	fp          string
	fins        *os.File
	w           *bufio.Writer
	currentSize int64
	opts        *RotateOptions
	closed      bool

	fname string
	fext  string
	timer *time.Timer
}

func NewAutoRotateFile(fp string, opts *RotateOptions) io.WriteCloser {
	if opts == nil {
		opts = &RotateOptions{}
	}
	if opts.BufferSize < 4096 {
		opts.BufferSize = 4096
	}

	ext := path.Ext(fp)

	arf := &_AutoRotateFile{fp: fp, opts: opts, fext: ext, fname: fp[0 : len(fp)-len(ext)]}
	if arf.opts.Kind != RotateKindFileSize {
		arf.rotate_by_time()
	}
	return arf
}

func (arf *_AutoRotateFile) do_rotate(newname string) {
	err := arf.do_close()
	if err != nil {
		err = fmt.Errorf("rotate error: when flush current file, %s", err)
		if arf.opts.OnRotateError != nil {
			arf.opts.OnRotateError(err)
		} else {
			fmt.Println(err)
		}
	}
	err = os.Rename(arf.fp, newname)
	if err != nil {
		err = fmt.Errorf("rotate error: when rename, %s", err)
		if arf.opts.OnRotateError != nil {
			arf.opts.OnRotateError(err)
		} else {
			fmt.Println(err)
		}
	}
}

func (arf *_AutoRotateFile) rotate_by_time() {
	now := time.Now()
	var end time.Time
	var timepart string
	switch arf.opts.Kind {
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
		arf.do_rotate(newname)
		arf.rotate_by_time()
	})
}

func (arf *_AutoRotateFile) do_close() error {
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
	return arf.do_close()
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
		arf.w = bufio.NewWriterSize(arf.fins, arf.opts.BufferSize)
	}

	n, err = arf.w.Write(p)
	if err != nil {
		return n, err
	}

	arf.currentSize += int64(len(p))
	if arf.opts.Kind == RotateKindFileSize && arf.currentSize >= arf.opts.MaxSize {
		arf.do_rotate(fmt.Sprintf(`%s.%d%s`, arf.fname, time.Now().UnixNano(), arf.fext))
	}
	return n, nil
}

var _ io.WriteCloser = (*_AutoRotateFile)(nil)
