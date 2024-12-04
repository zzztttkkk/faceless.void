package accountmodel_test

import (
	"encoding/json"
	"fmt"
	accountmodel "hello/sitea/modules/account/model"
	"testing"

	fv "github.com/zzztttkkk/faceless.void"
	"github.com/zzztttkkk/faceless.void/sqlx/sqltypes"
)

func TestJSON(t *testing.T) {
	var user accountmodel.UserModel
	user.Id.Value = 1
	user.Name.Value = "0.0"

	fmt.Println(string(fv.Must(json.Marshal(&user))))

	sqltypes.Char("name", 12).Primary().Unique().DefaultExpr("xxx")
	sqltypes.VarChar("name", 34).Primary().Comment("0.0")
	sqltypes.BigInt("id").Primary().AutoIncr()
}
