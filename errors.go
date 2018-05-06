package mns

import "errors"

var (
	errMnsQueueNnameIsTooLong             = errors.New("queue name is too long, the max length is 256")
	errMnsDelaySecondsRangeError          = errors.New("queue delay seconds is not in range of (0~60480)")
	errMnsMaxMessageSizeRangeError        = errors.New("max message size is not in range of (1024~65536)")
	errMnsMsgRetentionPeriodRangeError    = errors.New("message retention period is not in range of (60~129600)")
	errMnsMsgVisibilityTimeoutRangeError  = errors.New("message visibility timeout is not in range of (1~43200)")
	errMnsMsgPoolingWaitSecondsRangeError = errors.New("message poolling wait seconds is not in range of (0~30)")
	errMnsGetQueueRetNumberRangeError     = errors.New("get queue list param of ret number is not in range of (1~1000)")
)
