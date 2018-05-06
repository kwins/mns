package mns

// QueueCfg queue cfg
type QueueCfg struct {
	Host         string `json:"host"`
	AccessID     string `json:"access_id"`
	AccessSecret string `json:"access_secret"`
	QueueName    string `json:"queue_name"`
	QPS          int64  `json:"qps"` // 限制每秒发起请求数量 ，如果为0则不限制
}

// NewQueueCfg new queue cfg
func NewQueueCfg(host, accessID, accessSecret, queueName string, qps int64) QueueCfg {
	return QueueCfg{
		Host:         host,
		AccessID:     accessID,
		AccessSecret: accessSecret,
		QueueName:    queueName,
		QPS:          qps,
	}
}
