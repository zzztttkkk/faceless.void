package account

import (
	"context"
	"encoding/json"
	"net/http"

	fv "github.com/zzztttkkk/faceless.void"
)

func init() {
	fv.Endpoint().
		Func(
			func(ctx context.Context, req *http.Request, respw http.ResponseWriter) error {
				var params RegisterParams
				result, err := Register(ctx, &params)
				if err != nil {
					return err
				}
				enc := json.NewEncoder(respw)
				return enc.Encode(result)
			},
		).
		Register()
}
