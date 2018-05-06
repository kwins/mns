package mns

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// QConsumer queue consumer
// 实现此接口，如果返回error不等与nil，则会重发消息
// 直到MNS消息保留的最大时长
// 否则代表消费成功，删除消息
type QConsumer interface {
	HandleMsg(MessageResp) error
	Error(error) error
}

// QProducer queue producer
type QProducer struct {
	c *Client
}

// NewQProducer NewQProducer
func NewQProducer(c *Client) *QProducer {
	return &QProducer{c: c}
}

// SendMsg 发送消息
func (producer *QProducer) SendMsg(queueName string, message Message) error {
	resource := fmt.Sprintf("queues/%s/%s", queueName, "messages")
	if message.Priority == 0 {
		message.Priority = 8
	}
	resp, err := producer.c.Send("POST", nil, &message, resource)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return toErr(resp)
}

// BatchSendMsg 批量发送消息,返回error代表不成功，error为错误信息
func (producer *QProducer) BatchSendMsg(queueName string, message ...Message) error {
	resource := fmt.Sprintf("queues/%s/%s", queueName, "messages")
	var batchMsg BatchMessage
	batchMsg.Messages = make([]Message, len(message))
	for i, v := range message {
		if batchMsg.Messages[i].Priority == 0 {
			batchMsg.Messages[i].Priority = 8
		}
		batchMsg.Messages[i] = v
	}
	resp, err := producer.c.Send("POST", nil, &batchMsg, resource)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return toErr(resp)
}

func toErr(resp *http.Response) error {
	if resp.StatusCode != http.StatusCreated { // 创建失败
		decoder := xml.NewDecoder(resp.Body)
		var errMsg ErrorMessage
		if err := decoder.Decode(&errMsg); err != nil {
			return err
		}
		return errMsg.Error()
	}
	return nil
}
