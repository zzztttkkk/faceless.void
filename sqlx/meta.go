package sqlx

import (
	"database/sql"
	"encoding/json"
)

type TableMetaInfo struct {
	Name string
}

type ITableMate interface {
	TableMetaInfo() TableMetaInfo
}

type FieldMetaInfo struct {
	Name       string
	SqlType    string
	AutoIncr   bool
	PrimaryKey bool
	Unique     bool
	Nullable   bool
	Default    sql.NullString
	Check      string
	Comment    string
}

type IFieldMeta interface {
	FieldMetaInfo() FieldMetaInfo
}

type Field[T any, M IFieldMeta] struct {
	meta  [0]M
	Value T
}

func (field *Field[T, M]) Ptr() *T {
	return &field.Value
}

func (field *Field[T, M]) Metainfo() *FieldMetaInfo {
	return nil
}

func (field *Field[T, M]) MarshalJSON() ([]byte, error) {
	return json.Marshal(field.Value)
}

func (field *Field[T, M]) UnmarshalJSON(bs []byte) error {
	return json.Unmarshal(bs, &field.Value)
}

type opKind int

const (
	opKindEq = opKind(iota)
	opKindLt
	opKindGt
)

type Op struct {
	left  string
	op    string
	right any
	args  []any
}

func (field *Field[T, M]) Eq(val any) Op {
	return Op{
		left:  field.Metainfo().Name,
		op:    "==",
		right: val,
	}
}

func (field *Field[T, M]) Lt(val any) Op {
	return Op{
		left:  field.Metainfo().Name,
		op:    "<",
		right: val,
	}
}

func (field *Field[T, M]) Gt(val any) Op {
	return Op{
		left:  field.Metainfo().Name,
		op:    ">",
		right: val,
	}
}

func (field *Field[T, M]) Not(val any) Op {
	return Op{
		op:    "!",
		right: val,
	}
}
