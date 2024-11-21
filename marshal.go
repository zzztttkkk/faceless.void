package fv

import (
	"context"
	"encoding/json"
	"io"
)

var (
	defaultHttpMarshaler IHttpMarshaler
)

func SetDefaultHTTPMarshaler(v IHttpMarshaler) {
	defaultHttpMarshaler = v
}

func getHttpMarshal(ctx context.Context) IHttpMarshaler {
	av := ctx.Value(ctxKeyForHttpMarshaler)
	if av != nil {
		return av.(IHttpMarshaler)
	}
	return defaultHttpMarshaler
}

type jsonHttpMarshaler struct{}

// ContentType implements IHttpMarshaler.
func (jsonHttpMarshaler) ContentType() string {
	return "application/json"
}

// Marshal implements IHttpMarshaler.
func (jsonHttpMarshaler) Marshal(v any, buf io.Writer) error {
	enc := json.NewEncoder(buf)
	return enc.Encode(v)
}

var _ IHttpMarshaler = jsonHttpMarshaler{}
