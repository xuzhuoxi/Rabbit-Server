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
type OnTransCtx func(ctx *SqlCtx, index int)
type OnTransCommit func(err error)

type SqlCtx struct {
	Sql     string
	Args    []interface{}
	IsQuery bool
	Result  sql.Result
	Rows    *sql.Rows
}

func (o SqlCtx) String() string {
	return fmt.Sprintf("{Sql=%s, Args=[%s]}", o.Sql, fmt.Sprint(o.Args...))
}

func NewIDataSource(config CfgDataSourceItem) IDataSource {
	return NewDataSource(config)
}

func NewDataSource(config CfgDataSourceItem) *DataSource {
	return &DataSource{Config: config}
}

type IDataSource interface {
	eventx.IEventDispatcher

	// IsOpen
	// 判断当前数据源是否已连接
	IsOpen() bool
	// GetMeta
	// 取得当前数据源的元数据
	GetMeta() DatabaseMeta

	// Open
	// 开始启用数据源连接
	Open()
	// Close
	// 开始关闭数据源连接
	Close()

	// UpdateMeta
	// 更新当前数据源的元数据信息
	UpdateMeta()

	// SimpleQuery
	// 执行简单sql语句查询
	SimpleQuery(query string, onQuery OnQuery)
	// Query
	// 执行查询语句
	Query(query string, onQuery OnQuery, args ...interface{})
	// Update
	// 执行更新语句
	Update(query string, onUpdate OnUpdate, args ...interface{})
	// ExecTrans
	// 执行事务
	ExecTrans(sqlCtx []*SqlCtx, onTransCtx OnTransCtx, onTransCommit OnTransCommit) (err error)
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

func (o *DataSource) GetMeta() DatabaseMeta {
	return o.Meta
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

func (o *DataSource) UpdateMeta() {
	o.queryTableMeta()
}

func (o *DataSource) SimpleQuery(query string, onQuery OnQuery) {
	o.simpleQuery(query, onQuery)
}

func (o *DataSource) Query(query string, onQuery OnQuery, args ...interface{}) {
	o.query(query, onQuery, args...)
}

func (o *DataSource) Update(query string, onUpdate OnUpdate, args ...interface{}) {
	o.update(query, onUpdate, args...)
}

func (o *DataSource) ExecTrans(sqlCtx []*SqlCtx, onTransCtx OnTransCtx, onTransCommit OnTransCommit) (err error) {
	return o.execTrans(sqlCtx, onTransCtx, onTransCommit)
}
