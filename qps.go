package mns

import (
	"log"
	"sync/atomic"
	"time"
)

// QPS 限速
type QPS struct {
	qpsLimit       int64
	currentIndex   int64
	perSecondQuery int64
}

// NewQPS NewQPS
func NewQPS(qpsLimit int64) *QPS {
	qps := QPS{
		qpsLimit:       qpsLimit,
		perSecondQuery: 0,
	}
	return &qps
}

// Pulse Pulse
func (qps *QPS) pulse() {
	index := time.Now().Unix() % 5
	if index != qps.currentIndex {
		atomic.StoreInt64(&qps.currentIndex, index)
		atomic.StoreInt64(&qps.perSecondQuery, 0)
	}
	atomic.AddInt64(&qps.perSecondQuery, 1)
}

// CheckQPS CheckQPS
func (qps *QPS) CheckQPS() {
	qps.pulse()
	if qps.qpsLimit > 0 {
		for qps.perSecondQuery > qps.qpsLimit {
			qps.pulse()
			log.Printf("wait 10 millisecond,qps:%d", qps.perSecondQuery)
			time.Sleep(time.Millisecond * 10)
		}
	}
}
