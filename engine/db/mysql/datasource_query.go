// Package mysql
// Create on 2023/8/18
// @author xuzhuoxi
package mysql

import (
	"errors"
	"fmt"
)

func (o *DataSource) simpleQuery(query string, onQuery OnQuery) {
	rows, err := o.Db.Query(query)
	if nil != rows {
		defer rows.Close()
	}
	InvokeOnQuery(onQuery, rows, err)
}

func (o *DataSource) query(query string, onQuery OnQuery, args ...interface{}) {
	stmt, err1 := o.Db.Prepare(query)
	if err1 != nil {
		err1 = errors.New(fmt.Sprintf("Prepare failed,%s", err1))
		InvokeOnQuery(onQuery, nil, err1)
		return
	}
	defer stmt.Close()
	rows, err2 := stmt.Query(args...) // 执行预编译语句，传入参数
	if nil != rows {
		defer rows.Close()
	}
	InvokeOnQuery(onQuery, rows, err2)
}
