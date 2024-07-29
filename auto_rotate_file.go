package fv

import (
	"bufio"
	"io"
	"os"
	"sync"
)

type RotateKind int

const (
	RotateKindDaily RotateKind = RotateKind(iota)
	RotateKindHourly
	RotateKindMinutely
	RotateKindFileSize
)

type RotateOptions struct {
	Kind       RotateKind
	MaxSize    int64
	BufferSize int
}

type _AutoRotateFile struct {
	lock        sync.Mutex
	fp          string
	fins        *os.File
	w           *bufio.Writer
	currentSize int64
	prevWriteAt int64
	opts        *RotateOptions
}

func NewAutoRotateFile(fp string, opts *RotateOptions) io.WriteCloser {
	if opts == nil {
		opts = &RotateOptions{}
	}
	if opts.BufferSize < 4096 {
		opts.BufferSize = 4096
	}
	return &_AutoRotateFile{fp: fp, opts: opts}
}

func (arf *_AutoRotateFile) rename() {}

func (arf *_AutoRotateFile) rotate(p []byte) (bool, error) {
	if arf.fins == nil {
		return false, nil
	}

	switch arf.opts.Kind {
	case RotateKindFileSize:
		{
			if arf.currentSize+int64(len(p)) > arf.opts.MaxSize {
				_, err := arf.w.Write(p)
				if err != nil {
					return false, err
				}
				err = arf.do_close()
				if err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}
	case RotateKindDaily:
		{
			break
		}
	case RotateKindHourly:
		{
			break
		}
	case RotateKindMinutely:
		{
			break
		}
	}

	return false, nil
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
	defer arf.lock.Unlock()
	return arf.do_close()
}

func (arf *_AutoRotateFile) Write(p []byte) (n int, err error) {
	arf.lock.Lock()
	defer arf.lock.Unlock()

	writen, err := arf.rotate(p)
	if err != nil {
		return 0, err
	}
	if writen {
		return len(p), nil
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
	return n, nil
}

var _ io.WriteCloser = (*_AutoRotateFile)(nil)
