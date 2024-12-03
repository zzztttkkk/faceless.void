package fv

import (
	"context"
	"database/sql"

	"github.com/zzztttkkk/faceless.void/internal"
)

type _DBSource struct {
	main *sql.DB
	subs []*sql.DB
}

var dbsource *_DBSource

type SqlExecutor interface {
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type _txptr struct {
	ptr *sql.Tx
}

func (v _txptr) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return v.ptr.ExecContext(ctx, query, args...)
}

func (v _txptr) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return v.ptr.QueryContext(ctx, query, args...)
}

var _ SqlExecutor = _txptr{}

func WithTx(ctx context.Context) (context.Context, SqlExecutor) {
	return withTxOptions(ctx, nil)
}

func WithReadonlyTx(ctx context.Context) (context.Context, SqlExecutor) {
	var opts sql.TxOptions
	opts.ReadOnly = true
	return withTxOptions(ctx, &opts)
}

func withTxOptions(ctx context.Context, opts *sql.TxOptions) (context.Context, SqlExecutor) {
	txa := ctx.Value(internal.CtxKeyForSqlTx)
	if txa != nil {
		return ctx, txa.(_txptr)
	}
	var tx *sql.Tx
	var err error
	if opts != nil && opts.ReadOnly && len(dbsource.subs) > 0 {
		tx, err = dbsource.subs[0].BeginTx(ctx, opts)
	} else {
		tx, err = dbsource.main.BeginTx(ctx, opts)
	}
	if err != nil {
		panic(err)
	}
	val := _txptr{tx}
	return context.WithValue(ctx, internal.CtxKeyForSqlTx, val), val
}
