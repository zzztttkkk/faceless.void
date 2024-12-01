package fv

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func LoadEnv(files ...string) (map[string]string, error) {
	var kvs = map[string]string{}
	for _, fn := range files {
		if err := loadOneEnvFile(fn, kvs); err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
	}
	return kvs, nil
}

const (
	multiLinesStringStarts = `"""`
)

func loadOneEnvFile(fp string, kvs map[string]string) error {
	fv, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer fv.Close()

	reader := bufio.NewReader(fv)

	var line []byte

	var inMultiLines = false
	var currentKey string
	var currentValue []string

	online := func() {
		txt := string(line)
		trimed := strings.TrimSpace(txt)
		if inMultiLines {
			if strings.HasSuffix(trimed, multiLinesStringStarts) {
				currentValue = append(currentValue, trimed[:len(trimed)-3])
				inMultiLines = false
				kvs[currentKey] = strings.Join(currentValue, "\n")
				return
			}
			currentValue = append(currentValue, txt)
			return
		}

		if strings.HasPrefix(trimed, "#") {
			return
		}

		idx := strings.Index(txt, "=")
		if idx < 0 {
			kvs[strings.TrimSpace(txt)] = ""
			return
		}

		key := strings.TrimSpace(txt[:idx])
		val := strings.TrimSpace(txt[idx+1:])
		if strings.HasPrefix(val, multiLinesStringStarts) {
			currentKey = key
			currentValue = append(currentValue, val[3:])
			inMultiLines = true
			return
		}
		kvs[key] = val
	}

	for {
		rd, isprefix, err := reader.ReadLine()
		if err != nil {
			return err
		}
		if isprefix {
			line = append(line, rd...)
			continue
		}
		line = rd
		online()
		line = line[:0]
	}
}
