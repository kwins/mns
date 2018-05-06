package mns

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
)

// BaseResponse BaseResponse
type BaseResponse struct {
	XMLName   xml.Name `xml:"Message" json:"-"`
	Code      string   `xml:"Code,omitempty" json:"code,omitempty"`
	Message   string   `xml:"Message,omitempty" json:"message,omitempty"`
	RequestID string   `xml:"RequestId,omitempty" json:"request_id,omitempty"`
	HostID    string   `xml:"HostId,omitempty" json:"host_id,omitempty"`
}

// ErrorMessage ErrorMessage
type ErrorMessage struct {
	XMLName   xml.Name `xml:"Error" json:"-"`
	Code      string   `xml:"Code,omitempty" json:"code,omitempty"`
	Message   string   `xml:"Message,omitempty" json:"message,omitempty"`
	RequestID string   `xml:"RequestId,omitempty" json:"request_id,omitempty"`
	HostID    string   `xml:"HostId,omitempty" json:"host_id,omitempty"`
}

// Error error
func (msg *ErrorMessage) Error() error {
	return fmt.Errorf("HostID:%s HostID:%s Code:%s Message:%s",
		msg.HostID, msg.RequestID, msg.Code, msg.Message)
}

// BatchMessageResp BatchMessageResp
type BatchMessageResp struct {
	XMLName  xml.Name      `xml:"Messages" json:"-"`
	Messages []MessageResp `xml:"Message" json:"messages"`
}

// MessageResp MessageResp
type MessageResp struct {
	BaseResponse
	ReceiptHandle    string      `xml:"ReceiptHandle" json:"receipt_handle"`
	MessageBodyMD5   string      `xml:"MessageBodyMD5" json:"message_body_md5"`
	MessageBody      Base64Bytes `xml:"MessageBody" json:"message_body"`
	EnqueueTime      int64       `xml:"EnqueueTime" json:"enqueue_time"`
	NextVisibleTime  int64       `xml:"NextVisibleTime" json:"next_visible_time"`
	FirstDequeueTime int64       `xml:"FirstDequeueTime" json:"first_dequeue_time"`
	DequeueCount     int64       `xml:"DequeueCount" json:"dequeue_count"`
	Priority         int64       `xml:"Priority" json:"priority"`
}

// Message Mns message
type Message struct {
	XMLName      xml.Name    `xml:"Message"`
	MessageBody  Base64Bytes `xml:"MessageBody"`
	DelaySeconds int64       `xml:"DelaySeconds"`
	Priority     int64       `xml:"Priority"`
}

// BatchMessage batch message
type BatchMessage struct {
	XMLName  xml.Name  `xml:"Messages" json:"-"`
	Messages []Message `xml:"Message" json:"messages"`
}

// Base64Bytes Base64Bytes
type Base64Bytes []byte

// MarshalXML MarshalXML
func (p Base64Bytes) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(base64.StdEncoding.EncodeToString(p), start)
	return nil
}

// UnmarshalXML UnmarshalXML
func (p *Base64Bytes) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var content string
	if e := d.DecodeElement(&content, &start); e != nil {
		err = e
		return
	}

	if buf, e := base64.StdEncoding.DecodeString(content); e != nil {
		err = e
		return
	} else {
		*p = Base64Bytes(buf)
	}

	return nil
}
