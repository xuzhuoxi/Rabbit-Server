// Package mysql
// Create on 2023/8/18
// @author xuzhuoxi
package mysql

import (
	"database/sql"
	"fmt"
)

func (o *DataSource) queryTableMeta() {
	query := fmt.Sprintf(QueryTableMeta, o.Config.Schema)
	o.simpleQuery(query, o.onTableMeta)
}

func (o *DataSource) onTableMeta(rows *sql.Rows, err error) {
	if nil != err {
		o.DispatchEvent(EventOnDataSourceMetaUpdated, o, err)
		return
	}
	var tables []TableMeta
	for rows.Next() {
		meta := TableMeta{TableSchema: o.Config.Schema}
		err1 := rows.Scan(&meta.TableName, &meta.AvgRowLen,
			&meta.DataLen, &meta.MaxDataLen, &meta.IndexLen)
		if nil != err1 {
			o.DispatchEvent(EventOnDataSourceMetaUpdated, o, err1)
			return
		}
		tables = append(tables, meta)
	}
	o.Meta = DatabaseMeta{SchemaName: o.Config.Schema, Tables: tables}
	o.queryColumnMeta()
}

func (o *DataSource) queryColumnMeta() {
	query := fmt.Sprintf(QueryColumnMeta, o.Config.Schema)
	o.simpleQuery(query, o.onColumnMeta)
}

func (o *DataSource) onColumnMeta(rows *sql.Rows, err error) {
	if nil != err {
		o.DispatchEvent(EventOnDataSourceMetaUpdated, o, err)
		return
	}
	var columns []ColMeta
	for rows.Next() {
		meta := ColMeta{TableSchema: o.Config.Schema}
		err := rows.Scan(&meta.TableName, &meta.ColName,
			&meta.Position, &meta.Nullable, &meta.DataType, &meta.ColKey)
		if nil != err {
			o.DispatchEvent(EventOnDataSourceMetaUpdated, o, err)
			return
		}
		columns = append(columns, meta)
	}
	index := 0
	//fmt.Println("Column Size:", len(columns))
	for idxC := range columns {
		if columns[idxC].TableName != o.Meta.Tables[index].TableName {
			index += 1
		}
		o.Meta.Tables[index].Columns = append(o.Meta.Tables[index].Columns, columns[idxC])
	}
	o.DispatchEvent(EventOnDataSourceMetaUpdated, o, nil)
}
