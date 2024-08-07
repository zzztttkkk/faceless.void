package fv

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func loadOneEnvFile(fp string) error {
	file, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}

		idx := strings.Index(line, "=")
		if idx < 0 {
			return fmt.Errorf("invalid env line: %s", line)
		}
		key, val := strings.ToUpper(strings.TrimSpace(line[:idx])), strings.TrimSpace(line[idx+1:])
		err = os.Setenv(key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func Load[T any](typehint T, fps ...string) (*T, error) {
	return nil, nil
}

var bytestype = reflect.TypeOf((*[]byte)(nil)).Elem()

func parseOneEnvValue(ctx context.Context, rtype reflect.Type, txt string, readUri bool) (reflect.Value, error) {
	getrealbytes := func() ([]byte, error) {
		if readUri {
			if strings.HasPrefix(txt, "file://") {
				return os.ReadFile(txt[7:])
			}
			//goland:noinspection HttpUrlsUsage
			if strings.HasPrefix(txt, "http://") || strings.HasPrefix(txt, "https://") {
				req, err := http.NewRequestWithContext(ctx, "get", txt, nil)
				if err != nil {
					return nil, err
				}
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()
				return io.ReadAll(resp.Body)
			}
		}
		return []byte(txt), nil
	}

	getrealtxt := func() (string, error) {
		bs, err := getrealbytes()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(bs)), nil
	}

	rrtype := rtype

	isptr := false
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
		isptr = true
	}
	switch rtype.Kind() {
	case reflect.Func, reflect.Ptr, reflect.Chan:
		{
			panic(fmt.Errorf("%s", rrtype))
		}
	default:
		{
		}
	}

	var vv reflect.Value

	switch rtype.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			var n int64
			var e error

			txt, e = getrealtxt()
			if e != nil {
				return reflect.Value{}, e
			}
			if strings.HasPrefix(txt, "0x") {
				n, e = strconv.ParseInt(txt[2:], 16, 64)
			} else {
				n, e = strconv.ParseInt(txt, 10, 64)
			}
			if e != nil {
				return reflect.ValueOf(nil), e
			}
			vv = reflect.New(rtype).Elem()
			vv.SetInt(n)
			break
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			var n uint64
			var e error

			txt, e = getrealtxt()
			if e != nil {
				return reflect.Value{}, e
			}
			if strings.HasPrefix(txt, "0x") {
				n, e = strconv.ParseUint(txt[2:], 16, 64)
			} else {
				n, e = strconv.ParseUint(txt, 10, 64)
			}
			if e != nil {
				return reflect.ValueOf(nil), e
			}
			vv = reflect.New(rtype).Elem()
			vv.SetUint(n)
			break
		}
	case reflect.Float32, reflect.Float64:
		{
			txt, e := getrealtxt()
			if e != nil {
				return reflect.Value{}, e
			}
			n, e := strconv.ParseFloat(txt, 64)
			if e != nil {
				return reflect.ValueOf(nil), e
			}
			vv = reflect.New(rtype).Elem()
			vv.SetFloat(n)
			break
		}
	case reflect.String:
		{
			txt, e := getrealtxt()
			if e != nil {
				return reflect.Value{}, e
			}
			vv = reflect.New(rtype).Elem()
			vv.SetString(txt)
			break
		}
	case reflect.Bool:
		{
			txt, e := getrealtxt()
			if e != nil {
				return reflect.Value{}, e
			}
			n, e := strconv.ParseBool(txt)
			if e != nil {
				return reflect.ValueOf(nil), e
			}
			vv = reflect.New(rtype).Elem()
			vv.SetBool(n)
			break
		}
	default:
		{
			bs, e := getrealbytes()
			if e != nil {
				return reflect.Value{}, e
			}
			if rtype.AssignableTo(bytestype) {
				vv = reflect.New(rtype).Elem()
				vv.SetBytes(bs)
			} else {
				vv = reflect.New(rtype)
				e = json.Unmarshal(bs, vv.Interface())
				if e != nil {
					return reflect.ValueOf(nil), e
				}
				vv = vv.Elem()
			}
			break
		}
	}
	if isptr {
		return vv.Addr(), nil
	}
	return vv, nil
}
