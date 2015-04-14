package ali_mqs

import (
	"fmt"
	"os"
	"strings"
)

var (
	RECEIVER_COUNT = 10
)

const (
	PROXY_PREFIX = "MQS_PROXY_"
	GLOBAL_PROXY = "MQS_GLOBAL_PROXY"
)

type AliMQSQueue interface {
	Name() string
	SendMessage(message MessageSendRequest) (resp MessageSendResponse, err error)
	ReceiveMessage(respChan chan MessageReceiveResponse, errChan chan error, waitseconds ...int64)
	PeekMessage(respChan chan MessageReceiveResponse, errChan chan error)
	DeleteMessage(receiptHandle string) (err error)
	ChangeMessageVisibility(receiptHandle string, visibilityTimeout int64) (resp MessageVisibilityChangeResponse, err error)
	Stop()
}

type MQSQueue struct {
	name     string
	client   MQSClient
	stopChan chan bool
}

func NewMQSQueue(name string, client MQSClient) AliMQSQueue {
	if name == "" {
		panic("ali_mqs: queue name could not be empty")
	}

	queue := new(MQSQueue)
	queue.client = client
	queue.name = name
	queue.stopChan = make(chan bool)

	proxyURL := ""
	queueProxyEnvKey := PROXY_PREFIX + strings.Replace(strings.ToUpper(name), "-", "_", -1)
	if url := os.Getenv(queueProxyEnvKey); url != "" {
		proxyURL = url
	} else if globalurl := os.Getenv(GLOBAL_PROXY); globalurl != "" {
		proxyURL = globalurl
	}

	if proxyURL != "" {
		queue.client.SetProxy(proxyURL)
	}

	return queue
}

func (p *MQSQueue) Name() string {
	return p.name
}

func (p *MQSQueue) SendMessage(message MessageSendRequest) (resp MessageSendResponse, err error) {
	_, err = p.client.Send(POST, nil, message, fmt.Sprintf("%s/%s", p.name, "messages"), &resp)
	return
}

func (p *MQSQueue) Stop() {
	p.stopChan <- true
}

func (p *MQSQueue) ReceiveMessage(respChan chan MessageReceiveResponse, errChan chan error, waitseconds ...int64) {
	resource := fmt.Sprintf("%s/%s", p.name, "messages")
	if waitseconds != nil && len(waitseconds) == 1 {
		resource = fmt.Sprintf("%s/%s?waitseconds=%d", p.name, "messages", waitseconds[0])
	}

	for {
		resp := MessageReceiveResponse{}
		_, err := p.client.Send(GET, nil, nil, resource, &resp)
		if err != nil {
			errChan <- err
		} else {
			respChan <- resp
		}

		select {
		case _ = <-p.stopChan:
			{
				return
			}
		default:
		}
	}

	return
}

func (p *MQSQueue) PeekMessage(respChan chan MessageReceiveResponse, errChan chan error) {
	for {
		resp := MessageReceiveResponse{}
		_, err := p.client.Send(GET, nil, nil, fmt.Sprintf("%s/%s?peekonly=true", p.name, "messages"), &resp)
		if err != nil {
			errChan <- err
		} else {
			respChan <- resp
		}
	}
	return
}

func (p *MQSQueue) DeleteMessage(receiptHandle string) (err error) {
	_, err = p.client.Send(DELETE, nil, nil, fmt.Sprintf("%s/%s?ReceiptHandle=%s", p.name, "messages", receiptHandle), nil)
	return
}

func (p *MQSQueue) ChangeMessageVisibility(receiptHandle string, visibilityTimeout int64) (resp MessageVisibilityChangeResponse, err error) {
	_, err = p.client.Send(PUT, nil, nil, fmt.Sprintf("%s/%s?ReceiptHandle=%s&VisibilityTimeout=%d", p.name, "messages", receiptHandle, visibilityTimeout), &resp)
	return
}
