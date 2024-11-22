package fv

type HttpError struct {
}

// Error implements error.
func (h *HttpError) Error() string {
	panic("unimplemented")
}

var _ error = (*HttpError)(nil)
