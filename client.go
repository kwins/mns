package mns

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	mnsVersion     = "2015-06-06"
	defaultTimeout = int64(35)
)

// Client client
type Client struct {
	c             *http.Client
	cre           Credential
	consumerQueue string // 服务消费的队列
	accessID      string
	accessSecret  string
	timeOut       int64
	url           string
	consumer      QConsumer
	qps           *QPS
}

// NewClient new client,consumer==nil,则不消费消息
func NewClient(cfg QueueCfg, consumer QConsumer) *Client {
	client := new(Client)
	client.url = cfg.Host
	client.accessID = cfg.AccessID
	client.accessSecret = cfg.AccessSecret
	client.consumerQueue = cfg.QueueName
	client.consumer = consumer
	if client.consumer != nil {
		go client.recvMsgs()
	}
	client.cre = NewCredential(cfg.AccessSecret)
	client.timeOut = defaultTimeout
	client.qps = NewQPS(10)

	// set request timeout
	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, time.Second*5)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		ResponseHeaderTimeout: time.Second * 30,
	}

	client.c = &http.Client{
		Transport: transport,
	}
	return client
}

// Send send mns request
func (client *Client) Send(
	method string, h http.Header,
	message interface{}, resource string) (*http.Response, error) {

	client.qps.CheckQPS()

	var err error
	var body []byte

	if message == nil {
		body = nil
	} else {
		switch m := message.(type) {
		case []byte:
			body = m
		default:
			if body, err = xml.Marshal(message); err != nil {
				return nil, err
			}
		}
	}

	bodyMd5 := md5.Sum(body)
	bodyMd5Str := fmt.Sprintf("%x", bodyMd5)
	if h == nil {
		h = make(http.Header)
	}

	h.Add("x-mns-version", mnsVersion)
	h.Add("Content-Type", "application/xml")
	h.Add("Content-MD5", base64.StdEncoding.EncodeToString([]byte(bodyMd5Str)))
	h.Add("Date", time.Now().UTC().Format(http.TimeFormat))

	client.cre.SetHeader(h)
	client.cre.SetMethod(method)
	client.cre.SetResource(resource)
	signStr, err := client.cre.Sign()
	if err != nil {
		return nil, err
	}

	authSignStr := fmt.Sprintf("MNS %s:%s", client.accessID, signStr)
	h.Add("Authorization", authSignStr)

	url := client.url + "/" + resource
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = h
	return client.c.Do(req)
}

// DeleteMsg 删除消息
func (client *Client) DeleteMsg(receiptHandle string) error {
	resource := fmt.Sprintf("queues/%s/%s?ReceiptHandle=%s", client.consumerQueue, "messages", receiptHandle)
	resp, err := client.Send("DELETE", nil, nil, resource)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		decoder := xml.NewDecoder(resp.Body)
		var errMsg ErrorMessage
		if err := decoder.Decode(&errMsg); err != nil {
			return err
		}
		return errMsg.Error()
	}
	return nil
}

// recvMsgs 批量接受消息接受消息，最大限制16个消息
func (client *Client) recvMsgs() {
	resource := fmt.Sprintf("queues/%s/%s?numOfMessages=%d&waitseconds=%d", client.consumerQueue, "messages", 16, 30)
	for {
		resp, err := client.Send("GET", nil, nil, resource)
		if err != nil {
			if err := client.consumer.Error(err); err != nil {
				log.Println("mns:", err.Error())
			}
			continue
		}

		messages, err := takeResp(resp)
		if err != nil {
			if err := client.consumer.Error(err); err != nil {
				log.Println("mns:", err.Error())
			}
			continue
		}

		for _, message := range messages.Messages {
			if err := client.consumer.HandleMsg(message); err == nil {
				if err := client.DeleteMsg(message.ReceiptHandle); err != nil {
					log.Println("mns:", err.Error())
				}
			}
		}
	}
}

func takeResp(resp *http.Response) (*BatchMessageResp, error) {
	defer resp.Body.Close()
	decoder := xml.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK { // 接受消息失败
		var errMsg ErrorMessage
		if err := decoder.Decode(&errMsg); err != nil {
			return nil, err
		}
		// 消息对列表没有消息
		if strings.Contains(errMsg.Message, "MessageNotExist") {
			return nil, nil
		}
		return nil, errMsg.Error()
	}

	var messages BatchMessageResp
	if err := decoder.Decode(&messages); err != nil {
		return nil, err
	}
	return &messages, nil
}
