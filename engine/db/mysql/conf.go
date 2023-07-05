// Package mysql
// Create on 2023/7/4
// @author xuzhuoxi
package mysql

import "fmt"

type CfgDataSource struct {
	Name     string `yaml:"name"`
	Driver   string `yaml:"driver"`
	Url      string `yaml:"url"`
	UserName string `yaml:"username"`
	Passwd   string `yaml:"passwd"`
	Schema   string `yaml:"schema"`
}

func (o *CfgDataSource) DataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", o.UserName, o.Passwd, o.Url, o.Schema)
}

type CfgDataSources struct {
	Default         string          `yaml:"default"` // format: url/schema
	DataSources     []CfgDataSource `yaml:"data_sources"`
	QueryTableMeta  string          `yaml:"query_table_meta"`
	QueryColumnMeta string          `yaml:"query_column_meta"`
}
