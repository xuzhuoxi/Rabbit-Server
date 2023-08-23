// Package mysql
// Create on 2023/8/18
// @author xuzhuoxi
package mysql

import (
	"database/sql"
	"errors"
)

var errNoneCtx = errors.New("No SqlCtx Error! ")

func (o *DataSource) execTrans(sqlCtx []*SqlCtx, onTransCtx OnTransCtx, onTransCommit OnTransCommit) (err error) {
	if len(sqlCtx) == 0 {
		InvokeOnTransCommit(onTransCommit, errNoneCtx)
		return nil
	}
	tx, err1 := o.Db.Begin()
	if err1 != nil {
		InvokeOnTransCommit(onTransCommit, err1)
		return err1
	}
	for idx1 := range sqlCtx {
		err2 := o.execTransCtx(tx, sqlCtx[idx1])
		if nil != err2 {
			tx.Rollback()
			InvokeOnTransCommit(onTransCommit, err1)
			return err1
		}
		InvokeOnTransCtx(onTransCtx, sqlCtx[idx1], idx1)
		o.closeTransCtx(sqlCtx[idx1])
	}
	err1 = tx.Commit()
	if nil == err1 || err1 == sql.ErrTxDone {
		return nil
	}
	tx.Rollback()
	InvokeOnTransCommit(onTransCommit, err1)
	return err1
}

func (o *DataSource) closeTransCtx(ctx *SqlCtx) {
	if ctx.IsQuery && ctx.Rows != nil {
		_ = ctx.Rows.Close()
	}
}

func (o *DataSource) execTransCtx(tx *sql.Tx, ctx *SqlCtx) (err error) {
	if ctx.IsQuery {
		return o.execTransCtxRows(tx, ctx)
	} else {
		return o.execTransCtxResult(tx, ctx)
	}
}

func (o *DataSource) execTransCtxRows(tx *sql.Tx, ctx *SqlCtx) (err error) {
	if !ctx.IsQuery {
		err = o.execTransCtxResult(tx, ctx)
		return
	}
	stmt, err1 := tx.Prepare(ctx.Sql)
	if err1 != nil {
		err = err1
		return
	}
	rows, err2 := stmt.Query(ctx.Args...)
	if err2 != nil {
		err = err2
		return
	}
	ctx.Rows = rows
	return
}

func (o *DataSource) execTransCtxResult(tx *sql.Tx, ctx *SqlCtx) (err error) {
	stmt, err1 := tx.Prepare(ctx.Sql)
	if err1 != nil {
		err = err1
		return
	}
	rs, err3 := stmt.Exec(ctx.Args...)
	if err3 != nil {
		err = err3
		return
	}
	ctx.Result = rs
	return
}
