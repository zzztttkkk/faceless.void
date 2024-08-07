package account

import (
	"github.com/gin-gonic/gin"
	fv "github.com/zzztttkkk/faceless.void"
	"github.com/zzztttkkk/faceless.void/tests/exec/modules/internal"
)

func HttpRegister(ctx *gin.Context) {
}

func HttpLogin(ctx *gin.Context) {
}

func init() {
	internal.DIC.Register(func(groupGetter fv.TokenValueGetter[*gin.RouterGroup]) {
		group := groupGetter.Get("account")
		group.POST("/login", HttpLogin)
		group.POST("/create", HttpLogin)
	})
}
