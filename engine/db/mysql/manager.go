// Package mysql
// Create on 2023/7/4
// @author xuzhuoxi
package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xuzhuoxi/Rabbit-Server/engine/utils"
	"github.com/xuzhuoxi/infra-go/eventx"
)

func NewIDataSourceManager() IDataSourceManager {
	return NewDataSourceManager()
}

func NewDataSourceManager() *DataSourceManager {
	return &DataSourceManager{}
}

type IDataSourceManager interface {
	eventx.IEventDispatcher
	Init(cfgPath string) error
	OpenAll()
	UpdateMeta()
	CloseAll()

	List() []string
	GetDataSource(dbName string) IDataSource
	GetDefaultDataSource() IDataSource
}

type DataSourceManager struct {
	eventx.EventDispatcher
	Default     string
	DataSources []*DataSource
	Index       int
}

func (o *DataSourceManager) Init(cfgPath string) error {
	config := &CfgDataSources{}
	cfgPath = utils.FixFilePath(cfgPath)
	err := utils.UnmarshalFromYaml(cfgPath, config)
	if nil != err {
		o.DispatchEvent(EventOnManagerInited, o, err)
		return err
	}
	o.DataSources = nil
	for index := range config.DataSources {
		o.DataSources = append(o.DataSources, NewDataSource(config.DataSources[index]))
	}
	o.Default = config.Default
	o.DispatchEvent(EventOnManagerInited, o, nil)
	return nil
}

func (o *DataSourceManager) OpenAll() {
	o.Index = 0
	o.open()
}

func (o *DataSourceManager) open() {
	if o.Index >= len(o.DataSources) {
		o.DispatchEvent(EventOnManagerOpened, o, nil)
		return
	}
	ds := o.DataSources[o.Index]
	ds.OnceEventListener(EventOnDataSourceOpened, o.onOpened)
	ds.Open()
}

func (o *DataSourceManager) onOpened(evd *eventx.EventData) {
	err, isErr := evd.Data.(error)
	if isErr {
		o.DispatchEvent(EventOnManagerOpened, o, err)
		return
	}
	o.Index += 1
	o.open()
}

func (o *DataSourceManager) CloseAll() {
	o.Index = len(o.DataSources) - 1
	o.close()
}

func (o *DataSourceManager) close() {
	if o.Index < 0 {
		o.DispatchEvent(EventOnManagerClosed, o, nil)
		return
	}
	ds := o.DataSources[o.Index]
	ds.OnceEventListener(EventOnDataSourceClosed, o.onClosed)
	ds.Close()
}

func (o *DataSourceManager) onClosed(evd *eventx.EventData) {
	err, isErr := evd.Data.(error)
	if isErr {
		o.DispatchEvent(EventOnManagerClosed, o, err)
		return
	}
	o.Index -= 1
	o.close()
}

func (o *DataSourceManager) UpdateMeta() {
	o.Index = 0
	o.updateMeta()
}

func (o *DataSourceManager) updateMeta() {
	if o.Index >= len(o.DataSources) {
		o.DispatchEvent(EventOnManagerMetaUpdated, o, nil)
		return
	}
	ds := o.DataSources[o.Index]
	ds.OnceEventListener(EventOnDataSourceMetaUpdated, o.onMetaUpdated)
	ds.UpdateMeta()
}

func (o *DataSourceManager) onMetaUpdated(evd *eventx.EventData) {
	err, isErr := evd.Data.(error)
	if isErr {
		o.DispatchEvent(EventOnManagerMetaUpdated, o, err)
		return
	}
	o.Index += 1
	o.updateMeta()
}

func (o *DataSourceManager) List() []string {
	if len(o.DataSources) == 0 {
		return nil
	}
	rs := make([]string, len(o.DataSources))
	for index := range o.DataSources {
		rs[index] = o.DataSources[index].Config.Name
	}
	return rs
}

func (o *DataSourceManager) GetDataSource(dbName string) IDataSource {
	for index := range o.DataSources {
		if o.DataSources[index].Config.Name == dbName {
			return o.DataSources[index]
		}
	}
	return nil
}

func (o *DataSourceManager) GetDefaultDataSource() IDataSource {
	return o.GetDataSource(o.Default)
}
