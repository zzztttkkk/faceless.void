package fv

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

type HttpContext struct {
	Req *http.Request
}

func (hctx *HttpContext) Deadline() (deadline time.Time, ok bool) {
	return hctx.Req.Context().Deadline()
}

func (hctx *HttpContext) Done() <-chan struct{} {
	return hctx.Req.Context().Done()
}

func (hctx *HttpContext) Err() error {
	return hctx.Req.Context().Err()
}

func (hctx *HttpContext) Value(key any) any {
	return hctx.Req.Context().Value(key)
}

var _ context.Context = (*HttpContext)(nil)

func (hctx *HttpContext) Bind(v any) error {
	return nil
}

type IHttpResponse interface {
	Send(resp http.ResponseWriter)
}

type IHttpHandler interface {
	ServeHTTP(ctx *HttpContext) (IHttpResponse, error)
}

var (
	iHttpHandlerType = reflect.TypeOf((*IHttpHandler)(nil)).Elem()
	ictxType         = reflect.TypeOf((*context.Context)(nil)).Elem()
)

type HttpHandlerFunc func(req *HttpContext) (IHttpResponse, error)

func (fnc HttpHandlerFunc) ServeHTTP(req *HttpContext) (IHttpResponse, error) {
	return fnc(req)
}

func tostd(handler IHttpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := HttpContext{
			Req: r,
		}
		resp, err := handler.ServeHTTP(&ctx)
		if err != nil {
			panic(err)
		}
		if resp == nil {
			w.WriteHeader(200)
			return
		}
		resp.Send(w)
	}
}

type _EndpointOption struct {
	key string
	val any
}

var (
	EndpointOptions _EndpointOption
)

func (*_EndpointOption) Methods(methods ...string) _EndpointOption {
	return _EndpointOption{
		key: "methods",
		val: methods,
	}
}

func (*_EndpointOption) Pattern(pattern string) _EndpointOption {
	return _EndpointOption{
		key: "pattern",
		val: pattern,
	}
}

func (*_EndpointOption) Description(description string) _EndpointOption {
	return _EndpointOption{
		key: "description",
		val: description,
	}
}

type httpEndpoint struct {
	filename string
	methods  []string
	pattern  string
	handler  IHttpHandler
}

var (
	allEndpointsLock   sync.Mutex
	allEndpointsDone   bool
	allEndpoints       = make([]httpEndpoint, 0)
	funcAutoNameRegexp = regexp.MustCompile(`^func(\d+)$`)
)

func RegisterHttpEndpoint(fnc any, opts ..._EndpointOption) {
	rv := reflect.ValueOf(fnc)
	if rv.Kind() != reflect.Func {
		panic(fmt.Errorf("`%#v` is not a function", fnc))
	}
	funcv := runtime.FuncForPC(rv.Pointer())
	filename, _ := funcv.FileLine(0)

	endpoint := httpEndpoint{
		filename: filename,
		handler:  funcToHandler(funcv.Name(), rv),
	}

	for _, opt := range opts {
		switch opt.key {
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

func funcToHandler(funcname string, rv reflect.Value) IHttpHandler {
	hf, ok := rv.Interface().(HttpHandlerFunc)
	if ok {
		return hf
	}
	hfRaw, ok := rv.Interface().(func(req *HttpContext) (IHttpResponse, error))
	if ok {
		return HttpHandlerFunc(hfRaw)
	}

	rt := rv.Type()
	numin := rt.NumIn()

	var argPeeks []func(req *HttpContext) (reflect.Value, error)
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
		if isptr {
			argPeeks = append(argPeeks, func(req *HttpContext) (reflect.Value, error) {
				ptrv := reflect.New(argT)
				return ptrv, req.Bind(ptrv.Interface())
			})
		} else {
			argPeeks = append(argPeeks, func(req *HttpContext) (reflect.Value, error) {
				ptrv := reflect.New(argT)
				err := req.Bind(ptrv.Interface())
				if err != nil {
					return reflect.Value{}, err
				}
				return ptrv.Elem(), nil
			})
		}
	}

	numout := rt.NumOut()

	var mkresp func(outs []reflect.Value) (IHttpResponse, error)

	if numout == 0 {
		mkresp = func(_ []reflect.Value) (IHttpResponse, error) {
			return nil, nil
		}
	} else {
	}

	return HttpHandlerFunc(func(req *HttpContext) (IHttpResponse, error) {
		args := make([]reflect.Value, 0, numin)
		if firstArgIsCtx {
			args = append(args, reflect.ValueOf((context.Context)(req)))
		}
		for _, peek := range argPeeks {
			argv, err := peek(req)
			if err != nil {
				return nil, err
			}
			args = append(args, argv)
		}
		outs := rv.Call(args)
		return mkresp(outs)
	})
}
