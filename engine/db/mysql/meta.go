// Package mysql
// Create on 2023/7/3
// @author xuzhuoxi
package mysql

import "fmt"

type ColKeyValue = string

const (
	// ColKeyNull 没有被索引，或者只是一个多列非唯一索引的次要列
	ColKeyNull ColKeyValue = ""
	// ColKeyPri 主键，或者是一个多列主键的一部分。
	ColKeyPri ColKeyValue = "PRI"
	// ColKeyUni 唯一索引的第一列。
	// （唯一索引允许多个 NULL 值，但你可以通过检查 IS_NULLABLE 列来判断列是否允许 NULL。）
	ColKeyUni ColKeyValue = "UNI"
	// ColKeyMul 一个非唯一索引的第一列，且该列允许多个相同的值。
	ColKeyMul ColKeyValue = "MUL"
)

type ColMeta struct {
	TableSchema string      `db:"TABLE_SCHEMA"`
	TableName   string      `db:"TABLE_NAME"`
	ColName     string      `db:"COLUMN_NAME"`
	Position    int64       `db:"ORDINAL_POSITION"`
	Nullable    string      `db:"IS_NULLABLE"`
	DataType    string      `db:"DATA_TYPE"`
	ColKey      ColKeyValue `db:"COLUMN_KEY"`
}

func (o ColMeta) IsNullable() bool {
	return o.Nullable != "NO"
}

func (o ColMeta) String() string {
	return fmt.Sprintf("{%s, %d, %v, %s, %s}",
		o.ColName, o.Position, o.IsNullable(), o.DataType, o.ColKey)
}

type TableMeta struct {
	TableSchema string `db:"TABLE_SCHEMA"`
	TableName   string `db:"TABLE_NAME"`
	AvgRowLen   int64  `db:"AVG_ROW_LENGTH"`
	DataLen     int64  `db:"DATA_LENGTH"`
	MaxDataLen  int64  `db:"MAX_DATA_LENGTH"`
	IndexLen    int64  `db:"INDEX_LENGTH"`
	Columns     []ColMeta
}

func (o *TableMeta) String() string {
	return fmt.Sprintf("{%s, %d, %d, %d, %d, %v}",
		o.TableName, o.AvgRowLen, o.DataLen, o.MaxDataLen, o.IndexLen, o.Columns)
}

func (o *TableMeta) GetPriKeys() []string {
	var rs []string
	for index := range o.Columns {
		if o.Columns[index].ColKey == ColKeyPri {
			rs = append(rs, o.Columns[index].ColName)
		}
	}
	return rs
}

type DatabaseMeta struct {
	SchemaName string
	Tables     []TableMeta
}

func (o *DatabaseMeta) String() string {
	return fmt.Sprintf("{%s, %v}", o.SchemaName, o.Tables)
}

func (o *DatabaseMeta) GetTables() []string {
	if len(o.Tables) == 0 {
		return nil
	}
	rs := make([]string, len(o.Tables))
	for index := range o.Tables {
		rs[index] = o.Tables[index].TableName
	}
	return rs
}

func (o *DatabaseMeta) GetTableMetas() []TableMeta {
	if len(o.Tables) == 0 {
		return nil
	}
	rs := make([]TableMeta, len(o.Tables))
	copy(rs, o.Tables)
	return rs
}

func (o *DatabaseMeta) GetTableMeta(tableName string) (meta TableMeta, ok bool) {
	if len(o.Tables) == 0 {
		return
	}
	for index := range o.Tables {
		if o.Tables[index].TableName == tableName {
			meta, ok = o.Tables[index], true
		}
	}
	return
}
