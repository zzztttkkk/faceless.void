package fv

import (
	"encoding/json"
	"net/http"
)

type IMarshaler interface {
	Marshal(w http.ResponseWriter, val any) error
}

type _StdJSONMarshaler struct{}

func (_ _StdJSONMarshaler) Marshal(w http.ResponseWriter, val any) error {
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	return enc.Encode(val)
}

var StdJSONMarshaler IMarshaler = _StdJSONMarshaler{}
