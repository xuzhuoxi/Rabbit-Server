// Package mysql
// Create on 2023/7/4
// @author xuzhuoxi
package mysql

import "fmt"

type CfgDataSourceItem struct {
	Name     string `yaml:"name"`     // 数据源配置名称
	Driver   string `yaml:"driver"`   // 数据库驱动名称
	Url      string `yaml:"url"`      // 数据库连接地址
	UserName string `yaml:"username"` // 连接数据库的用户名
	Passwd   string `yaml:"passwd"`   // 连接数据库的密码
	Schema   string `yaml:"schema"`   // 连接数据库的库名
}

func (o *CfgDataSourceItem) DataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", o.UserName, o.Passwd, o.Url, o.Schema)
}

type CfgDataSource struct {
	Default         string              `yaml:"default"`           // 默认的数据源配置项
	DataSources     []CfgDataSourceItem `yaml:"data_sources"`      // 数据源配置项列表
	QueryTableMeta  string              `yaml:"query_table_meta"`  // 用于查询表元数据的Sql语言
	QueryColumnMeta string              `yaml:"query_column_meta"` // 用于查询表字段元数据的Sql语言
}
