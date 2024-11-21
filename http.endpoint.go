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
)

type endpointOption struct {
	key string
	val any
}

var (
	EndpointOptions endpointOption
)

func (*endpointOption) Filename(filename string) endpointOption {
	return endpointOption{
		key: "filename",
		val: filename,
	}
}

func (*endpointOption) Methods(methods ...string) endpointOption {
	return endpointOption{
		key: "methods",
		val: methods,
	}
}

func (*endpointOption) Pattern(pattern string) endpointOption {
	return endpointOption{
		key: "pattern",
		val: pattern,
	}
}

func (*endpointOption) Description(description string) endpointOption {
	return endpointOption{
		key: "description",
		val: description,
	}
}

func (*endpointOption) Input(types ...reflect.Type) endpointOption {
	return endpointOption{
		key: "input",
		val: types,
	}
}

func (*endpointOption) Output(types ...reflect.Type) endpointOption {
	return endpointOption{
		key: "output",
		val: types,
	}
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
	if len(endpoint.methods) > 0 && !slices.Contains(endpoint.methods, req.Method) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := context.WithValue(req.Context(), ctxKeyForHttpRequest, req)
	resp, err := endpoint.handler.ServeHTTP(ctx)
	if err != nil {
		panic(err)
	}
	if resp == nil {
		rw.WriteHeader(200)
		rw.Write(nil)
		return
	}
	resp.Send(ctx, rw)
}

var (
	allEndpointsLock   sync.Mutex
	allEndpointsDone   bool
	allEndpoints       = make([]httpEndpoint, 0)
	funcAutoNameRegexp = regexp.MustCompile(`^func(\d+)$`)
)

func RegisterHttpEndpoint(fnc any, opts ...endpointOption) {
	rv := reflect.ValueOf(fnc)
	if rv.Kind() != reflect.Func {
		panic(fmt.Errorf("`%#v` is not a function", fnc))
	}
	funcv := runtime.FuncForPC(rv.Pointer())
	filename, _ := funcv.FileLine(0)

	endpoint := httpEndpoint{filename: filename}
	endpoint.handler = endpoint.funcToHandler(funcv.Name(), rv)

	for _, opt := range opts {
		switch opt.key {
		case "filename":
			{
				endpoint.filename = opt.val.(string)
				break
			}
		case "methods":
			{
				endpoint.methods = opt.val.([]string)
				break
			}
		case "pattern":
			{
				endpoint.pattern = opt.val.(string)
				break
			}
		case "input":
			{
				endpoint.argTypes = opt.val.([]reflect.Type)
				break
			}
		case "output":
			{
				endpoint.outTypes = opt.val.([]reflect.Type)
				break
			}
		}
	}
	if endpoint.pattern == "" {
		parts := strings.Split(funcv.Name(), ".")
		name := parts[len(parts)-1]
		if funcAutoNameRegexp.MatchString(name) {
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

type IHttpError interface {
	error
	StatusCode() int
	BodyMessage(ctx context.Context) []byte
}

func (endpoint *httpEndpoint) funcToHandler(funcname string, rv reflect.Value) IHttpHandler {
	hf, ok := rv.Interface().(HttpHandlerFunc)
	if ok {
		return hf
	}
	hfRaw, ok := rv.Interface().(func(req context.Context) (IHttpResponse, error))
	if ok {
		return HttpHandlerFunc(hfRaw)
	}

	rt := rv.Type()
	numin := rt.NumIn()

	var argPeeks []func(req context.Context) (reflect.Value, error)
	var firstArgIsCtx bool
	for i := 0; i < numin; i++ {
		argT := rt.In(i)
		if argT.Implements(ictxType) && i == 0 {
			firstArgIsCtx = true
			continue
		}
		isptr := false
		if argT.Kind() == reflect.Pointer {
			argT = argT.Elem()
			isptr = true
		}
		if argT.Kind() != reflect.Struct {
			panic(fmt.Errorf("`function %s`'s param type is not a struct or a struct pointer, at %d", funcname, i))
		}
		endpoint.argTypes = append(endpoint.argTypes, argT)

		if isptr {
			argPeeks = append(argPeeks, func(ctx context.Context) (reflect.Value, error) {
				ptrv := reflect.New(argT)
				return ptrv, bindHttp(ctx, ptrv.Interface())
			})
		} else {
			argPeeks = append(argPeeks, func(ctx context.Context) (reflect.Value, error) {
				ptrv := reflect.New(argT)
				err := bindHttp(ctx, ptrv.Interface())
				if err != nil {
					return reflect.Value{}, err
				}
				return ptrv.Elem(), nil
			})
		}
	}

	numout := rt.NumOut()
	var mkresp func(outs []reflect.Value) (IHttpResponse, error)
	switch numout {
	case 0:
		{
			mkresp = func(_ []reflect.Value) (IHttpResponse, error) {
				return codeResponse(http.StatusOK), nil
			}
			break
		}
	case 1:
		{
			outType := rt.Out(0)
			if outType == reflect.TypeOf((*int)(nil)).Elem() {
				mkresp = func(outs []reflect.Value) (IHttpResponse, error) {
					return codeResponse(outs[0].Int()), nil
				}
				break
			}
			mkresp = func(outs []reflect.Value) (IHttpResponse, error) {
				return anyResponse{val: outs[0].Interface()}, nil
			}
		}
	case 2:
		{
			firstOutType, secondOutType := rt.Out(0), rt.Out(1)
			fmt.Println(firstOutType, secondOutType, "2121")
		}
	}

	return HttpHandlerFunc(func(ctx context.Context) (IHttpResponse, error) {
		if endpoint.marshaler != nil {
			ctx = context.WithValue(ctx, ctxKeyForHttpMarshaler, endpoint.marshaler)
		}

		args := make([]reflect.Value, 0, numin)
		if firstArgIsCtx {
			args = append(args, reflect.ValueOf(ctx))
		}
		for _, peek := range argPeeks {
			argv, err := peek(ctx)
			if err != nil {
				return nil, err
			}
			args = append(args, argv)
		}
		outs := rv.Call(args)
		return mkresp(outs)
	})
}

func mountEndpoints(mux *http.ServeMux, root string, globs ...string) {
	allEndpointsLock.Lock()
	defer allEndpointsLock.Unlock()

	if allEndpointsDone {
		panic("x")
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

func RunHTTP(port int, main func(), globs []string) {
	mainv := reflect.ValueOf(main)
	mainfunc := runtime.FuncForPC(mainv.Pointer())
	if mainfunc.Name() != "main.main" {
		panic(fmt.Errorf("`%#v` is not the main function of main package", main))
	}
	maingo, _ := mainfunc.FileLine(0)
	rootpkg := filepath.Dir(maingo)
	fmt.Println(rootpkg)

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
