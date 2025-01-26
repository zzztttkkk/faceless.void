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

type _EndpointBuilder struct {
	pairs      []internal.Pair[string]
	middleware []HttpMiddlewareFunc
	funced     bool
}

var (
	allEndpointBuilders sync.Map
)

func Endpoint() *_EndpointBuilder {
	obj := &_EndpointBuilder{}
	allEndpointBuilders.Store(obj, true)
	return obj
}

func (builder *_EndpointBuilder) Filename(filename string) *_EndpointBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("filename", filename))
	return builder
}

func (builder *_EndpointBuilder) Methods(methods ...string) *_EndpointBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("methods", methods))
	return builder
}

func (builder *_EndpointBuilder) Pattern(pattern string) *_EndpointBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("pattern", pattern))
	return builder
}

func (builder *_EndpointBuilder) Description(description string) *_EndpointBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("description", description))
	return builder
}

func (builder *_EndpointBuilder) Input(vals ...any) *_EndpointBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("input", vals))
	return builder
}

func (builder *_EndpointBuilder) Output(vals ...any) *_EndpointBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("output", vals))
	return builder
}

func (builder *_EndpointBuilder) Use(fncs ...HttpMiddlewareFunc) *_EndpointBuilder {
	builder.middleware = append(builder.middleware, fncs...)
	return builder
}

func (builder *_EndpointBuilder) Func(fnc HandleFunc) *_EndpointBuilder {
	if builder.funced {
		panic("endpoint already has handle function")
	}
	builder.funced = true
	builder.pairs = append(builder.pairs, internal.PairOf("func", fnc))
	return builder
}

func (builder *_EndpointBuilder) Endable(enable bool) {
	if !enable {
		allEndpointBuilders.Delete(builder)
	}
}

type IHttpMarshaler interface {
	ContentType() string
	Marshal(v any, buf io.Writer) error
}

type _HTTPEndpoint struct {
	filename  string
	methods   []string
	pattern   string
	handler   HandleFunc
	argTypes  []reflect.Type
	outTypes  []reflect.Type
	marshaler IHttpMarshaler
	appscope  *TAppScope
}

// ServeHTTP implements http.Handler.
func (endpoint *_HTTPEndpoint) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		return
	}

	if len(endpoint.methods) > 0 && !slices.Contains(endpoint.methods, req.Method) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := context.WithValue(req.Context(), internal.CtxKeyForAppScope, endpoint.appscope)
	err := endpoint.handler(ctx, req, rw)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(nil)
		return
	}
}

var (
	anonymousFuncNameRegexp = regexp.MustCompile(`^func(\d+)$`)
)

func registerHttpEndpoint(opts *_EndpointBuilder) {
	endpoint := _HTTPEndpoint{}

	var fnc HandleFunc
	for _, opt := range opts.pairs {
		switch opt.Key {
		case "filename":
			{
				endpoint.filename = opt.Val.(string)
				break
			}
		case "methods":
			{
				endpoint.methods = opt.Val.([]string)
				break
			}
		case "pattern":
			{
				endpoint.pattern = opt.Val.(string)
				break
			}
		case "input":
			{
				endpoint.argTypes = opt.Val.([]reflect.Type)
				break
			}
		case "output":
			{
				endpoint.outTypes = opt.Val.([]reflect.Type)
				break
			}
		case "func":
			{
				fnc = opt.Val.(HandleFunc)
				break
			}
		}
	}

	var funcv *runtime.Func
	var filename string
	funcv = runtime.FuncForPC(reflect.ValueOf(fnc).Pointer())
	filename, _ = funcv.FileLine(0)

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
}

func mountEndpoints(endpoint *_HTTPEndpoint, mux *http.ServeMux, root string, globs ...string) {
	matched := false
	for _, glob := range globs {
		if ok, _ := filepath.Match(glob, endpoint.filename); ok {
			matched = true
			break
		}
	}
	if !matched {
		return
	}

	rel, err := filepath.Rel(root, endpoint.filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.HasPrefix(rel, ".") {
		fmt.Println(">>>")
		return
	}

	fmt.Println(fmt.Sprintf("%s/%s", rel, endpoint.pattern))
	mux.Handle(fmt.Sprintf("%s/%s", rel, endpoint.pattern), endpoint)
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
}
