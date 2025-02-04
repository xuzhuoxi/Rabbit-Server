// Package server
// Create on 2025/2/2
// @author xuzhuoxi
package server

// IRequestPool
// 请求参数集的对象池接口
type IRequestPool interface {
	// Register 注册创建方法
	Register(newFunc func() IExtensionRequest)
	// GetInstance 获取一个实例
	GetInstance() IExtensionRequest
	// Recycle 回收一个实例
	Recycle(instance IExtensionRequest) bool
}

// IResponsePool
// 响应参数集的对象池接口
type IResponsePool interface {
	// Register 注册创建方法
	Register(newFunc func() IExtensionResponse)
	// GetInstance 获取一个实例
	GetInstance() IExtensionResponse
	// Recycle 回收一个实例
	Recycle(instance IExtensionResponse) bool
}

// INotifyPool
// 通知参数集的对象池接口
type INotifyPool interface {
	// Register 注册创建方法
	Register(newFunc func() IExtensionNotify)
	// GetInstance 获取一个实例
	GetInstance() IExtensionNotify
	// Recycle 回收一个实例
	Recycle(instance IExtensionNotify) bool
}
