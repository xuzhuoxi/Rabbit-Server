// Package verify
// Create on 2023/8/23
// @author xuzhuoxi
package verify

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/config"
	"github.com/xuzhuoxi/infra-go/extendx/protox"
	"sync"
	"time"
)

func newReqLog() *reqLog {
	rs := &reqLog{origin: make([]int64, 0, 32)}
	rs.ReqStamps = rs.origin[0:0]
	return rs
}

type reqLog struct {
	Name      string
	PId       string
	ReqStamps []int64
	origin    []int64
}

func (o *reqLog) ReplaceStamp(newStamp int64) {
	o.ReqStamps = append(o.origin, o.ReqStamps...)
	o.ReqStamps = append(o.ReqStamps, newStamp)
}

func (o *reqLog) AppendStamp(newStamp int64) {
	o.ReqStamps = append(o.ReqStamps, newStamp)
}

func (o *reqLog) SetStamp(newStamp int64) {
	o.ReqStamps = append(o.origin, newStamp)
}

func NewRabbitVerify(cfg config.CfgVerifyRoot) *RabbitVerify {
	return &RabbitVerify{
		CfgVerifyRoot: cfg,
		Logs:          make(map[string]*reqLog, 2048),
	}
}

type RabbitVerify struct {
	CfgVerifyRoot config.CfgVerifyRoot
	Logs          map[string]*reqLog
	Lock          sync.RWMutex
}

func (o *RabbitVerify) Verify(name string, pid string, uid string) (rsCode int32) {
	found := o.CfgVerifyRoot.FindVerify(name, pid)
	log := o.findLog(uid)
	nowStamp := time.Now().UnixNano()
	if log.Name != name || log.PId != pid {
		log.Name, log.PId = name, pid
		log.SetStamp(nowStamp)
		return protox.CodeSuc
	}
	//fmt.Println("RabbitVerify.Verify", name, pid, uid, found, log.ReqStamps)
	if found.MinFreqVal > 0 && (nowStamp-log.ReqStamps[len(log.ReqStamps)-1]) < int64(found.MinFreqVal) {
		log.ReplaceStamp(nowStamp)
		return protox.CodeFreq
	}
	if found.MaxPerSec == 0 || len(log.ReqStamps) < found.MaxPerSec {
		log.AppendStamp(nowStamp)
		return protox.CodeSuc
	}
	index := len(log.ReqStamps) - found.MaxPerSec
	if nowStamp-log.ReqStamps[index] < int64(time.Second) {
		log.ReplaceStamp(nowStamp)
		return protox.CodeFreq
	}
	log.AppendStamp(nowStamp)
	return protox.CodeSuc
}

func (o *RabbitVerify) findLog(uid string) *reqLog {
	o.Lock.RLock()
	if l, ok := o.Logs[uid]; ok {
		o.Lock.RUnlock()
		return l
	}
	o.Lock.RUnlock()
	o.Lock.Lock()
	defer o.Lock.Unlock()
	rs := newReqLog()
	o.Logs[uid] = rs
	return rs
}
