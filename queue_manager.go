package ali_mqs

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogap/errors"
)

type MQSLocation string

const (
	Beijing  MQSLocation = "beijing"
	Hangzhou MQSLocation = "hangzhou"
	Qingdao  MQSLocation = "qingdao"
)

type AliQueueManager interface {
	CreateQueue(location MQSLocation, queueName string, delaySeconds int32, maxMessageSize int32, messageRetentionPeriod int32, visibilityTimeout int32, pollingWaitSeconds int32) (err error)
	SetQueueAttributes(location MQSLocation, queueName string, delaySeconds int32, maxMessageSize int32, messageRetentionPeriod int32, visibilityTimeout int32, pollingWaitSeconds int32) (err error)
	GetQueueAttributes(location MQSLocation, queueName string) (attr QueueAttribute, err error)
	DeleteQueue(location MQSLocation, queueName string) (err error)
	ListQueue(location MQSLocation, marker string, retNumber int32, prefix string) (queues Queues, err error)
}

type MQSQueueManager struct {
	ownerId         string
	credential      Credential
	accessKeyId     string
	accessKeySecret string
}

func checkQueueName(queueName string) (err error) {
	if len(queueName) > 256 {
		err = ERR_MQS_QUEUE_NAME_IS_TOO_LONG.New()
		return
	}
	return
}

func checkDelaySeconds(seconds int32) (err error) {
	if seconds > 60480 || seconds < 0 {
		err = ERR_MQS_DELAY_SECONDS_RANGE_ERROR.New()
		return
	}
	return
}

func checkMaxMessageSize(maxSize int32) (err error) {
	if maxSize < 1024 || maxSize > 65536 {
		err = ERR_MQS_MAX_MESSAGE_SIZE_RANGE_ERROR.New()
		return
	}
	return
}

func checkMessageRetentionPeriod(retentionPeriod int32) (err error) {
	if retentionPeriod < 60 || retentionPeriod > 129600 {
		err = ERR_MQS_MSG_RETENTION_PERIOD_RANGE_ERROR.New()
		return
	}
	return
}

func checkVisibilityTimeout(visibilityTimeout int32) (err error) {
	if visibilityTimeout < 1 || visibilityTimeout > 43200 {
		err = ERR_MQS_MSG_VISIBILITY_TIMEOUT_RANGE_ERROR.New()
		return
	}
	return
}

func checkPollingWaitSeconds(pollingWaitSeconds int32) (err error) {
	if pollingWaitSeconds < 0 || pollingWaitSeconds > 30 {
		err = ERR_MQS_MSG_POOLLING_WAIT_SECONDS_RANGE_ERROR.New()
		return
	}
	return
}

func NewMQSQueueManager(ownerId, accessKeyId, accessKeySecret string) AliQueueManager {
	return &MQSQueueManager{
		ownerId:         ownerId,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}
}

func checkAttributes(delaySeconds int32, maxMessageSize int32, messageRetentionPeriod int32, visibilityTimeout int32, pollingWaitSeconds int32) (err error) {
	if err = checkDelaySeconds(delaySeconds); err != nil {
		return
	}
	if err = checkMaxMessageSize(maxMessageSize); err != nil {
		return
	}
	if err = checkMessageRetentionPeriod(messageRetentionPeriod); err != nil {
		return
	}
	if err = checkVisibilityTimeout(visibilityTimeout); err != nil {
		return
	}
	if err = checkPollingWaitSeconds(pollingWaitSeconds); err != nil {
		return
	}
	return
}

func (p *MQSQueueManager) CreateQueue(location MQSLocation, queueName string, delaySeconds int32, maxMessageSize int32, messageRetentionPeriod int32, visibilityTimeout int32, pollingWaitSeconds int32) (err error) {
	queueName = strings.TrimSpace(queueName)

	if err = checkQueueName(queueName); err != nil {
		return
	}

	if err = checkAttributes(delaySeconds,
		maxMessageSize,
		messageRetentionPeriod,
		visibilityTimeout,
		pollingWaitSeconds); err != nil {
		return
	}

	message := CreateQueueRequest{
		DelaySeconds:           delaySeconds,
		MaxMessageSize:         maxMessageSize,
		MessageRetentionPeriod: messageRetentionPeriod,
		VisibilityTimeout:      visibilityTimeout,
		PollingWaitSeconds:     pollingWaitSeconds,
	}

	url := fmt.Sprintf("http://%s.mqs-cn-%s.aliyuncs.com", p.ownerId, string(location))

	cli := NewAliMQSClient(url, p.accessKeyId, p.accessKeySecret)

	var code int
	code, err = cli.Send(_PUT, nil, &message, queueName, nil)

	if code == http.StatusNoContent {
		err = ERR_MQS_QUEUE_ALREADY_EXIST_AND_HAVE_SAME_ATTR.New(errors.Params{"name": queueName})
		return
	}

	return
}

func (p *MQSQueueManager) SetQueueAttributes(location MQSLocation, queueName string, delaySeconds int32, maxMessageSize int32, messageRetentionPeriod int32, visibilityTimeout int32, pollingWaitSeconds int32) (err error) {
	queueName = strings.TrimSpace(queueName)

	if err = checkQueueName(queueName); err != nil {
		return
	}

	if err = checkAttributes(delaySeconds,
		maxMessageSize,
		messageRetentionPeriod,
		visibilityTimeout,
		pollingWaitSeconds); err != nil {
		return
	}

	message := CreateQueueRequest{
		DelaySeconds:           delaySeconds,
		MaxMessageSize:         maxMessageSize,
		MessageRetentionPeriod: messageRetentionPeriod,
		VisibilityTimeout:      visibilityTimeout,
		PollingWaitSeconds:     pollingWaitSeconds,
	}

	url := fmt.Sprintf("http://%s.mqs-cn-%s.aliyuncs.com", p.ownerId, string(location))

	cli := NewAliMQSClient(url, p.accessKeyId, p.accessKeySecret)

	_, err = cli.Send(_PUT, nil, &message, fmt.Sprintf("%s?metaoverride=true", queueName), nil)
	return
}

func (p *MQSQueueManager) GetQueueAttributes(location MQSLocation, queueName string) (attr QueueAttribute, err error) {
	queueName = strings.TrimSpace(queueName)

	if err = checkQueueName(queueName); err != nil {
		return
	}

	url := fmt.Sprintf("http://%s.mqs-cn-%s.aliyuncs.com", p.ownerId, string(location))

	cli := NewAliMQSClient(url, p.accessKeyId, p.accessKeySecret)

	_, err = cli.Send(_GET, nil, nil, queueName, &attr)

	return
}

func (p *MQSQueueManager) DeleteQueue(location MQSLocation, queueName string) (err error) {
	queueName = strings.TrimSpace(queueName)

	if err = checkQueueName(queueName); err != nil {
		return
	}

	url := fmt.Sprintf("http://%s.mqs-cn-%s.aliyuncs.com", p.ownerId, string(location))

	cli := NewAliMQSClient(url, p.accessKeyId, p.accessKeySecret)

	_, err = cli.Send(_DELETE, nil, nil, queueName, nil)

	return
}

func (p *MQSQueueManager) ListQueue(location MQSLocation, marker string, retNumber int32, prefix string) (queues Queues, err error) {

	url := fmt.Sprintf("http://%s.mqs-cn-%s.aliyuncs.com", p.ownerId, string(location))

	cli := NewAliMQSClient(url, p.accessKeyId, p.accessKeySecret)

	header := map[string]string{}

	marker = strings.TrimSpace(marker)
	if marker != "" {
		header["x-mqs-marker"] = marker
	}

	if retNumber > 0 {
		if retNumber >= 1 && retNumber <= 1000 {
			header["x-mqs-ret-number"] = strconv.Itoa(int(retNumber))
		} else {
			err = REE_MQS_GET_QUEUE_RET_NUMBER_RANGE_ERROR.New()
			return
		}
	}

	prefix = strings.TrimSpace(prefix)
	if prefix != "" {
		header["x-mqs-prefix"] = prefix
	}

	_, err = cli.Send(_GET, header, nil, "", &queues)

	return
}
