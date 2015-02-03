package ali_mqs

import (
	"fmt"
)

type MQSQueue struct {
	name   string
	client MQSClient
}

func NewMQSQueue(name string, client MQSClient) *MQSQueue {
	if name == "" {
		panic("ali_mqs: queue name could not be empty")
	}

	queue := new(MQSQueue)
	queue.client = client
	queue.name = name
	return queue
}

func (p *MQSQueue) Name() string {
	return p.name
}

func (p *MQSQueue) SendMessage(message interface{}, delaySeconds, priority int32) (resp MessageSendResponse, err error) {
	err = p.client.Send(POST, nil, message, fmt.Sprintf("%s/%s", p.name, "messages"), &resp)
	return
}

func (p *MQSQueue) ReceiveMessage(respChan chan MessageReceiveResponse, errChan chan error, waitseconds ...int64) {
	resource := fmt.Sprintf("%s/%s", p.name, "messages")
	if waitseconds != nil && len(waitseconds) == 1 {
		resource = fmt.Sprintf("%s/%s?waitseconds=%d", p.name, "messages", waitseconds)
	}

	for {
		resp := MessageReceiveResponse{}
		err := p.client.Send(GET, nil, nil, resource, &resp)
		if err != nil {
			errChan <- err
		} else {
			respChan <- resp
		}
	}
	return
}

func (p *MQSQueue) PeekMessage(respChan chan MessageReceiveResponse, errChan chan error, waitseconds int64) {
	for {
		resp := MessageReceiveResponse{}
		err := p.client.Send(GET, nil, nil, fmt.Sprintf("%s/%s?peekonly=true", p.name, "messages"), &resp)
		if err != nil {
			errChan <- err
		} else {
			respChan <- resp
		}
	}
	return
}

func (p *MQSQueue) DeleteMessage(receiptHandle string) (err error) {
	err = p.client.Send(DELETE, nil, nil, fmt.Sprintf("%s/%s?ReceiptHandle=%s", p.name, "messages", receiptHandle), nil)
	return
}

func (p *MQSQueue) ChangeMessageVisibility(receiptHandle string, visibilityTimeout int64) (resp MessageVisibilityChangeResponse, err error) {
	err = p.client.Send(PUT, nil, nil, fmt.Sprintf("%s/%s?ReceiptHandle=%s&VisibilityTimeout=%d", p.name, "messages", receiptHandle, visibilityTimeout), &resp)
	return
}
