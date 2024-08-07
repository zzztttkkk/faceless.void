package modules

import (
	"github.com/gin-gonic/gin"
	_ "github.com/zzztttkkk/faceless.void/tests/exec/modules/account"
	"github.com/zzztttkkk/faceless.void/tests/exec/modules/internal"
)

func init() {
	internal.DIC.Register(func(engine *gin.Engine) {})
}
