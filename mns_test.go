package mns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

var c *Client
var q *QProducer

type ConsumerDemo struct {
}

// HandleMsg HandleMsg
func (cd *ConsumerDemo) HandleMsg(msg MessageResp) error {
	log.Println("msg:", string(msg.MessageBody))
	return nil
}

// HandleMsg HandleMsg
func (cd *ConsumerDemo) Error(err error) error {
	log.Println("error:", err.Error())
	return nil
}

func Init() {
	var cfg QueueCfg
	b, err := ioutil.ReadFile("./mns.json")
	if err != nil {
		panic(err.Error())
	}
	json.Unmarshal(b, &cfg)
	c = NewClient(cfg, &ConsumerDemo{})
	q = NewQProducer(c)
}

func TestSendMsg(t *testing.T) {
	Init()
	var msg Message
	msg.Priority = 0
	msg.MessageBody = []byte("hello world")
	if err := q.SendMsg("A-UPDATE-TEST", msg); err != nil {
		t.Error(err.Error())
	}
}

func TestBatchSendMsg(t *testing.T) {
	Init()
	var msg Message
	msg.MessageBody = []byte("hello world")

	var msg1 Message
	msg1.MessageBody = []byte("hello world1")
	if err := q.BatchSendMsg("A-UPDATE-TEST", msg, msg1); err != nil {
		t.Error(err.Error())
	}
}

func TestRecvMsg(t *testing.T) {
	Init()
	go func() {
		k := time.NewTicker(time.Second * 5)
		var i int
		for {
			<-k.C
			var msg Message
			msg.Priority = 0
			msg.MessageBody = []byte(fmt.Sprintf("hello world-%d", i))
			if err := q.SendMsg("A-UPDATE-TEST", msg); err != nil {
				t.Error(err.Error())
			}
			i++
		}
	}()
	time.Sleep(time.Hour)
}

func TestQPS(t *testing.T) {
	qps := NewQPS(10)
	for i := 0; i < 50; i++ {
		qps.CheckQPS()
		log.Println("====now:", time.Now().Unix(), " currentIndex:", qps.currentIndex, " perSecondQuery", qps.perSecondQuery)
		// time.Sleep(time.Millisecond * 10)
	}
}
