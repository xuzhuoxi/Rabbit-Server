// Package mysql
// Create on 2023/7/3
// @author xuzhuoxi
package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/eventx"
)

var (
	// QueryTableMeta 用于查询表元数据的Sql语言
	// 在DataSourceManager初始化后，如果配置表上有配置这条Sql语言，则覆盖
	QueryTableMeta = "" +
		"SELECT TABLE_NAME, AVG_ROW_LENGTH, DATA_LENGTH, MAX_DATA_LENGTH, INDEX_LENGTH " +
		"FROM `information_schema`.`TABLES` " +
		"WHERE TABLE_SCHEMA = \"%s\" " +
		"ORDER BY `TABLE_NAME` ASC;"
	// QueryColumnMeta 用于查询表字段元数据的Sql语言
	// 在DataSourceManager初始化后，如果配置表上有配置这条Sql语言，则覆盖
	QueryColumnMeta = "" +
		"SELECT TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION, IS_NULLABLE, DATA_TYPE, COLUMN_KEY " +
		"FROM `information_schema`.`COLUMNS` " +
		"WHERE TABLE_SCHEMA = \"%s\" " +
		"ORDER BY `TABLE_NAME` ASC, `ORDINAL_POSITION` ASC;"
)

type OnQuery func(rows *sql.Rows, err error)
type OnUpdate func(rowLen int64, err error)

func NewIDataSource(config CfgDataSourceItem) IDataSource {
	return NewDataSource(config)
}

func NewDataSource(config CfgDataSourceItem) *DataSource {
	return &DataSource{Config: config}
}

type IDataSource interface {
	eventx.IEventDispatcher
	IsOpen() bool
	Open()
	Close()

	GetMeta() DatabaseMeta
	UpdateMeta()

	Query(query string, onQuery OnQuery)
	Update(query string, onUpdate OnUpdate, args ...interface{})
	//LoadUserData(userId string)
	//SaveUserData(userId string)
	//SaveAllUserData()
}

type DataSource struct {
	eventx.EventDispatcher
	Config CfgDataSourceItem

	Meta DatabaseMeta
	Db   *sql.DB

	open bool
}

func (o *DataSource) IsOpen() bool {
	return o.open
}

func (o *DataSource) Open() {
	db, err := sql.Open(o.Config.Driver, o.Config.DataSourceName())
	if nil != err {
		err = errors.New(fmt.Sprintf("Open mysql failed,%s", err))
		o.DispatchEvent(EventOnDataSourceOpened, o, err)
		return
	}
	o.Db = db
	o.open = true
	o.DispatchEvent(EventOnDataSourceOpened, o, nil)
}

func (o *DataSource) Close() {
	if !o.open || nil == o.Db {
		o.DispatchEvent(EventOnDataSourceClosed, o, nil)
		return
	}
	o.open = false
	err := o.Db.Close()
	if nil != err {
		err = errors.New(fmt.Sprintf("Close mysql failed,%s", err))
		o.DispatchEvent(EventOnDataSourceClosed, o, err)
		return
	}
	o.DispatchEvent(EventOnDataSourceClosed, o, nil)
}

func (o *DataSource) GetMeta() DatabaseMeta {
	return o.Meta
}

func (o *DataSource) UpdateMeta() {
	o.queryTableMeta()
}

func (o *DataSource) Query(query string, onQuery OnQuery) {
	o.query(query, onQuery)
}

func (o *DataSource) Update(query string, onUpdate OnUpdate, args ...interface{}) {
	stmt, err1 := o.Db.Prepare(query)
	if err1 != nil {
		err1 = errors.New(fmt.Sprintf("Prepare failed,%s", err1))
		onUpdate(0, err1)
		return
	}
	defer stmt.Close()
	res, err2 := stmt.Exec(args...) // 执行预编译语句，传入参数
	if err2 != nil {
		err2 = errors.New(fmt.Sprintf("Exec failed,%s", err2))
		onUpdate(0, err2)
		return
	}
	row, err3 := res.RowsAffected() // 获取影响的行数
	if err3 != nil {
		err3 = errors.New(fmt.Sprintf("Rows affected failed,%s", err3))
		onUpdate(0, err3)
		return
	}
	onUpdate(row, nil)
}

func (o *DataSource) query(query string, onQuery OnQuery) {
	rows, err := o.Db.Query(query)
	if nil != rows {
		defer rows.Close()
	}
	onQuery(rows, err)
}

func (o *DataSource) queryTableMeta() {
	query := fmt.Sprintf(QueryTableMeta, o.Config.Schema)
	o.query(query, o.onTableMeta)
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
	o.query(query, o.onColumnMeta)
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
