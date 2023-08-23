// Package mysql
// Create on 2023/8/18
// @author xuzhuoxi
package mysql

import (
	"errors"
	"fmt"
)

func (o *DataSource) update(query string, onUpdate OnUpdate, args ...interface{}) {
	stmt, err1 := o.Db.Prepare(query)
	if err1 != nil {
		err1 = errors.New(fmt.Sprintf("Prepare failed,%s", err1))
		InvokeOnUpdate(onUpdate, 0, err1)
		return
	}
	defer stmt.Close()
	res, err2 := stmt.Exec(args...) // 执行预编译语句，传入参数
	if err2 != nil {
		err2 = errors.New(fmt.Sprintf("Exec failed,%s", err2))
		InvokeOnUpdate(onUpdate, 0, err2)
		return
	}
	row, err3 := res.RowsAffected() // 获取影响的行数
	if err3 != nil {
		err3 = errors.New(fmt.Sprintf("Rows affected failed,%s", err3))
		InvokeOnUpdate(onUpdate, 0, err3)
		return
	}
	InvokeOnUpdate(onUpdate, row, nil)
}
