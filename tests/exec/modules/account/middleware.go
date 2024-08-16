package account

import "github.com/gin-gonic/gin"

var (
	middleware1 = []gin.HandlerFunc{
		func(ctx *gin.Context) {
		},
	} // go:fv http middleware
)

var (
	Middleware2, Middleware3 []gin.HandlerFunc // go:fv http middleware
	xpk                      int               // ddd
	a, b                     int
)
