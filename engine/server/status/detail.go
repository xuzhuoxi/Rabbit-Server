// Package status
// Create on 2023/6/14
// @author xuzhuoxi
package status

import (
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"runtime"
	"sync"
	"time"
)

func NewServerStatusDetail(id string, statsInterval int64) *ServerStatusDetail {
	return &ServerStatusDetail{
		UpdateDetailInfo: core.UpdateDetailInfo{
			Id:            id,
			StatsInterval: statsInterval,
		},
	}
}

type ServerStatusDetail struct {
	core.UpdateDetailInfo
	lock sync.RWMutex
}

// GetPassNano 启动时间
func (s *ServerStatusDetail) GetPassNano() int64 {
	return time.Now().UnixNano() - s.StartTimestamp
}

// StatsWeight
// 当前统计的服务权重(连接数 + 统计响应时间 / 统计时间 )
// 越大代表压力越大
func (s *ServerStatusDetail) StatsWeight() float64 {
	if 0 == s.StatsRespUnixNano {
		return 0
	} else {
		return s.StatsRespCoefficient()
	}
}

// RespCoefficient
// 响应系数(响应总时间/统计总时间),
// 注意：结果正常设置下为[0,1]
func (s *ServerStatusDetail) RespCoefficient() float64 {
	return float64(s.TotalRespTime) / (float64(s.GetPassNano()) * float64(runtime.NumCPU()))
}

// RespAvgTime
// 平均响应时间(响应总时间/响应次数)
func (s *ServerStatusDetail) RespAvgTime() float64 {
	return float64(s.TotalRespTime) / float64(s.TotalReqCount)
}

// ReqDensityTime
// 时间请求密度(次数/秒)
func (s *ServerStatusDetail) ReqDensityTime() int {
	pass := s.GetPassNano()
	return int(int64(time.Second) * s.TotalReqCount / pass)
}

// StatsRespCoefficient
// 区间响应系数(响应总时间/统计总时间),
// 注意：结果正常设置下为[0,1]
func (s *ServerStatusDetail) StatsRespCoefficient() float64 {
	return float64(s.StatsRespUnixNano) / (float64(s.getStatsPass()) * float64(runtime.NumCPU()))
}

// StatsRespAvgTime
// 区间平均响应时间(响应总时间/响应次数)
func (s *ServerStatusDetail) StatsRespAvgTime() float64 {
	return float64(s.StatsRespUnixNano) / float64(s.StatsReqCount)
}

// StatsReqDensityTime
// 区间时间请求密度(次数/秒)
func (s *ServerStatusDetail) StatsReqDensityTime() int {
	pass := s.getStatsPass()
	return int(int64(time.Second) * s.StatsReqCount / pass)
}

//-------------------------------

// Start 启动
func (s *ServerStatusDetail) Start() {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now().UnixNano()
	s.StartTimestamp = now
	s.StatsTimestamp = now
}

// AddLinkCount 增加一个连接
func (s *ServerStatusDetail) AddLinkCount() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Links++
	if s.Links > s.MaxLinks { //更新最大连接数
		s.MaxLinks = s.Links
	}
}

// RemoveLinkCount 减少一个连接
func (s *ServerStatusDetail) RemoveLinkCount() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Links--
}

//----------------------

// AddReqCount 增加一个请求
func (s *ServerStatusDetail) AddReqCount() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.statsFull() {
		s.statsReset()
	}
	s.TotalReqCount++
	s.StatsReqCount++
}

// AddRespUnixNano 增加响应时间量
func (s *ServerStatusDetail) AddRespUnixNano(unixNano int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.StatsRespUnixNano += unixNano
	s.TotalRespTime += unixNano
	if unixNano > s.MaxRespTime { //更新最大响应时间量
		s.MaxRespTime = unixNano
	}
}

// ReStats 重新统计
func (s *ServerStatusDetail) ReStats() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.statsReset()
}

//----------------------

//重置统计数据
func (s *ServerStatusDetail) statsReset() {
	s.StatsReqCount = 0
	s.StatsTimestamp = time.Now().UnixNano()
	s.StatsRespUnixNano = 0
}

func (s *ServerStatusDetail) getStatsPass() int64 {
	return time.Now().UnixNano() - s.StatsTimestamp
}

func (s *ServerStatusDetail) statsFull() bool {
	return s.getStatsPass() >= s.StatsInterval
}
