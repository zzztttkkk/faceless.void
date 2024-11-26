package fv

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/zzztttkkk/faceless.void/internal"
)

type endpointOptionKey int

type endpointOptionPair struct {
	key endpointOptionKey
	val any
}

const (
	endpointOptionKeyForFilename = endpointOptionKey(iota)
	endpointOptionKeyForMethods
	endpointOptionKeyForPattern
	endpointOptionKeyForDescription
	endpointOptionKeyForInput
	endpointOptionKeyForOutput
	endpointOptionKeyForHandleFunc
	endpointOptionKeyForAnyFunc
)

type endpointOptions struct {
	pairs []endpointOptionPair
}

func Endpoint() *endpointOptions {
	return &endpointOptions{}
}

func (opts *endpointOptions) Filename(filename string) *endpointOptions {
	opts.pairs = append(opts.pairs, endpointOptionPair{key: endpointOptionKeyForFilename, val: filename})
	return opts
}

func (opts *endpointOptions) Methods(methods ...string) *endpointOptions {
	opts.pairs = append(opts.pairs, endpointOptionPair{key: endpointOptionKeyForMethods, val: methods})
	return opts
}

func (opts *endpointOptions) Pattern(pattern string) *endpointOptions {
	opts.pairs = append(opts.pairs, endpointOptionPair{key: endpointOptionKeyForPattern, val: pattern})
	return opts
}

func (opts *endpointOptions) Description(description string) *endpointOptions {
	opts.pairs = append(opts.pairs, endpointOptionPair{key: endpointOptionKeyForDescription, val: description})
	return opts
}

func (opts *endpointOptions) Input(types ...reflect.Type) *endpointOptions {
	opts.pairs = append(opts.pairs, endpointOptionPair{key: endpointOptionKeyForInput, val: types})
	return opts
}

func (opts *endpointOptions) Output(types ...reflect.Type) *endpointOptions {
	opts.pairs = append(opts.pairs, endpointOptionPair{key: endpointOptionKeyForOutput, val: types})
	return opts
}

func (opts *endpointOptions) Func(fnc HttpHandlerFunc) {
	opts.pairs = append(opts.pairs, endpointOptionPair{endpointOptionKeyForHandleFunc, fnc})
}

func (opts *endpointOptions) AnyFunc(fnc any) {
	opts.pairs = append(opts.pairs, endpointOptionPair{endpointOptionKeyForAnyFunc, fnc})
}

type IHttpMarshaler interface {
	ContentType() string
	Marshal(v any, buf io.Writer) error
}

type httpEndpoint struct {
	filename  string
	methods   []string
	pattern   string
	handler   IHttpHandler
	argTypes  []reflect.Type
	outTypes  []reflect.Type
	marshaler IHttpMarshaler
}

// ServeHTTP implements http.Handler.
func (endpoint *httpEndpoint) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		return
	}

	if len(endpoint.methods) > 0 && !slices.Contains(endpoint.methods, req.Method) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var bh _Getter
	ctx := context.WithValue(req.Context(), internal.CtxKeyForHttpRequest, req)
	ctx = bh.init(ctx, req)
	err := endpoint.handler.ServeHTTP(ctx, req, rw)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(nil)
		return
	}
}

var (
	allEndpointsLock        sync.Mutex
	allEndpointsDone        bool
	allEndpoints            = make([]httpEndpoint, 0)
	anonymousFuncNameRegexp = regexp.MustCompile(`^func(\d+)$`)
)

func RegisterHttpEndpoint(opts *endpointOptions) {

	endpoint := httpEndpoint{}

	var fnc HttpHandlerFunc
	var anyfnc any

	for _, opt := range opts.pairs {
		switch opt.key {
		case endpointOptionKeyForFilename:
			{
				endpoint.filename = opt.val.(string)
				break
			}
		case endpointOptionKeyForMethods:
			{
				endpoint.methods = opt.val.([]string)
				break
			}
		case endpointOptionKeyForPattern:
			{
				endpoint.pattern = opt.val.(string)
				break
			}
		case endpointOptionKeyForInput:
			{
				endpoint.argTypes = opt.val.([]reflect.Type)
				break
			}
		case endpointOptionKeyForOutput:
			{
				endpoint.outTypes = opt.val.([]reflect.Type)
				break
			}
		case endpointOptionKeyForHandleFunc:
			{
				fnc = opt.val.(HttpHandlerFunc)
				break
			}
		case endpointOptionKeyForAnyFunc:
			{
				anyfnc = opt.val
				break
			}
		}
	}

	var funcv *runtime.Func
	var filename string
	if anyfnc != nil {
		rv := reflect.ValueOf(anyfnc)
		if rv.Kind() != reflect.Func {
			panic(fmt.Errorf("`%#v` is not a function", fnc))
		}
		funcv := runtime.FuncForPC(rv.Pointer())
		filename, _ = funcv.FileLine(0)
		fnc = endpoint.mkhandler(filename, rv)
	} else {
		funcv = runtime.FuncForPC(reflect.ValueOf(fnc).Pointer())
		filename, _ = funcv.FileLine(0)
	}
	endpoint.handler = fnc
	if endpoint.filename == "" {
		endpoint.filename = filename
	}

	if endpoint.pattern == "" {
		parts := strings.Split(funcv.Name(), ".")
		name := parts[len(parts)-1]
		if anonymousFuncNameRegexp.MatchString(name) {
			panic(fmt.Errorf("`%s` need a pattern option", funcv.Name()))
		}
	}

	allEndpointsLock.Lock()
	if allEndpointsDone {
		allEndpointsLock.Unlock()
		panic(fmt.Errorf("http endpoint register is already done."))
	}

	defer allEndpointsLock.Unlock()
	allEndpoints = append(allEndpoints, endpoint)
}

func (endpoint *httpEndpoint) mkhandler(funcname string, rv reflect.Value) HttpHandlerFunc {
	hf, ok := rv.Interface().(HttpHandlerFunc)
	if ok {
		return hf
	}
	hfRaw, ok := rv.Interface().(func(context.Context, *http.Request, http.ResponseWriter) error)
	if ok {
		return HttpHandlerFunc(hfRaw)
	}

	rt := rv.Type()
	numin := rt.NumIn()

	var argPeeks []func(ctx context.Context, req *http.Request) (reflect.Value, error)
	for i := 0; i < numin; i++ {
		argT := rt.In(i)
		if i == 0 {
			if !argT.Implements(ictxType) {
				panic(fmt.Errorf("`function %s`'s first param type must be `context.Context`", funcname))
			}
			continue
		}

		if !argT.Implements(iBindingType) {
			panic(fmt.Errorf("`function %s`'s param type is not a fv.IBinding, at %d", funcname, i))
		}

		isptr := false
		if argT.Kind() == reflect.Pointer {
			argT = argT.Elem()
			isptr = true
		}
		endpoint.argTypes = append(endpoint.argTypes, argT)

		if isptr {
			argPeeks = append(argPeeks, func(ctx context.Context, req *http.Request) (reflect.Value, error) {
				ptrv := reflect.New(argT)
				err := ptrv.Interface().(IBinding).Binding(ctx)
				if err != nil {
					return reflect.Value{}, err
				}
				return ptrv, nil
			})
		} else {
			argPeeks = append(argPeeks, func(ctx context.Context, req *http.Request) (reflect.Value, error) {
				ptrv := reflect.New(argT)
				err := ptrv.Interface().(IBinding).Binding(ctx)
				if err != nil {
					return reflect.Value{}, err
				}
				return ptrv.Elem(), nil
			})
		}
	}

	var send func(outs []reflect.Value, respw http.ResponseWriter) error
	switch rt.NumOut() {
	case 0:
		{
			send = func(_ []reflect.Value, respw http.ResponseWriter) error {
				respw.WriteHeader(http.StatusOK)
				_, err := respw.Write(nil)
				return err
			}
			break
		}
	case 1:
		{
			outtype := rt.Out(0)
			if outtype == reflect.TypeOf((*int)(nil)).Elem() {
				send = func(outs []reflect.Value, respw http.ResponseWriter) error {
					code := outs[0].Int()
					respw.WriteHeader(int(code))
					_, err := respw.Write(nil)
					return err
				}
				break
			}
		}
	}

	return HttpHandlerFunc(func(ctx context.Context, req *http.Request, respw http.ResponseWriter) error {
		if endpoint.marshaler != nil {
		}

		args := make([]reflect.Value, 0, numin)
		args = append(args, reflect.ValueOf(ctx))

		for _, peek := range argPeeks {
			argv, err := peek(ctx, req)
			if err != nil {
				return err
			}
			args = append(args, argv)
		}
		outs := rv.Call(args)
		return send(outs, respw)
	})
}

func mountEndpoints(mux *http.ServeMux, root string, globs ...string) {
	allEndpointsLock.Lock()
	defer allEndpointsLock.Unlock()

	if allEndpointsDone {
		panic("endpoint register is done")
	}

	for idx := range allEndpoints {
		endpoint := &(allEndpoints[idx])
		matched := false
		for _, glob := range globs {
			if ok, _ := filepath.Match(glob, endpoint.filename); ok {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}

		rel, err := filepath.Rel(root, endpoint.filename)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if strings.HasPrefix(rel, ".") {
			fmt.Println(">>>")
			continue
		}

		fmt.Println(fmt.Sprintf("%s/%s", rel, endpoint.pattern))
		mux.Handle(fmt.Sprintf("%s/%s", rel, endpoint.pattern), endpoint)
	}
}

func readPkgName(fp string) string {
	fv, err := os.Open(fp)
	if err != nil {
		panic(err)
	}
	defer fv.Close()

	reader := bufio.NewReader(fv)
	var linebuf []byte
	for {
		tmp, isPrefix, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		linebuf = append(linebuf, tmp...)
		if isPrefix {
			continue
		}
		line := strings.TrimSpace(string(linebuf))
		if strings.HasPrefix(line, "package ") {
			return line[8:]
		}
	}
}

func loadEndpoints(root string, globs []string) {
	fmt.Println(root)

	pkgs := map[string]struct{}{}

	for _, glob := range globs {
		matchers, _ := filepath.Glob(glob)
		for _, fn := range matchers {
			if !strings.HasSuffix(fn, ".go") {
				panic(fmt.Errorf("glob pattern `%s`, match non-go file: %s", glob, fn))
			}

			fn, err := filepath.Abs(fn)
			if err != nil {
				continue
			}
			pkgs[filepath.Dir(fn)] = struct{}{}
		}
	}

	fmt.Println(pkgs)
}

type HttpSite struct {
	Port          int
	EndpointsGlob string
	HostName      string
	TLSCert       string
	TLSKey        string
}

func RunHTTP(main func(), sites ...HttpSite) {
	mainv := reflect.ValueOf(main)
	mainfunc := runtime.FuncForPC(mainv.Pointer())
	if mainfunc.Name() != "main.main" {
		panic(fmt.Errorf("`%s` is not the main function of main package", mainfunc.Name()))
	}
	maingo, _ := mainfunc.FileLine(0)
	rootpkg := filepath.Dir(maingo)
	fmt.Println(rootpkg)

	allEndpointsLock.Lock()
	allEndpointsDone = true
	allEndpointsLock.Unlock()
}
