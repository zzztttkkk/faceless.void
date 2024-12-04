package accountmodel

import "github.com/zzztttkkk/faceless.void/sqlx"

type UserModel struct {
	Id   sqlx.Field[int64, idmeta]
	Name sqlx.Field[string, namemeta]
}

type idmeta int

func (_ idmeta) FieldMetaInfo() sqlx.FieldMetaInfo {
	return sqlx.FieldMetaInfo{
		Name:       "id",
		SqlType:    "bigint",
		PrimaryKey: true,
		AutoIncr:   true,
	}
}

type namemeta int

func (_ namemeta) FieldMetaInfo() sqlx.FieldMetaInfo {
	return sqlx.FieldMetaInfo{
		Name:     "name",
		SqlType:  "varchar(30)",
		Nullable: false,
		Unique:   true,
	}
}
