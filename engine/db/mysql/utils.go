// Package mysql
// Create on 2023/8/4
// @author xuzhuoxi
package mysql

import "database/sql"

func InvokeOnQuery(onQuery OnQuery, rows *sql.Rows, err error) {
	if onQuery == nil {
		return
	}
	onQuery(rows, err)
}

func InvokeOnUpdate(onUpdate OnUpdate, rowLen int64, err error) {
	if onUpdate == nil {
		return
	}
	onUpdate(rowLen, err)
}

func InvokeOnTransCtx(onTransCtx OnTransCtx, ctx *SqlCtx, index int) {
	if onTransCtx == nil {
		return
	}
	onTransCtx(ctx, index)
}

func InvokeOnTransCommit(onTransCommit OnTransCommit, err error) {
	if onTransCommit == nil {
		return
	}
	onTransCommit(err)
}
