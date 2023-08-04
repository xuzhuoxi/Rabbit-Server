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

func InvokeOnTrans(onTrans OnTrans, err error) {
	if onTrans == nil {
		return
	}
	onTrans(err)
}
