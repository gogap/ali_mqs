package ali_mqs

import (
	"fmt"
	"sync"
	"time"
)

var (
	RECEIVER_COUNT = 10
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
	name      string
	client    MQSClient
	stopChans []chan bool
}

func NewMQSQueue(name string, client MQSClient) AliMQSQueue {
	if name == "" {
		panic("ali_mqs: queue name could not be empty")
	}

	queue := new(MQSQueue)
	queue.client = client
	queue.name = name
	queue.stopChans = make([]chan bool, RECEIVER_COUNT)
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
	for i := 0; i < RECEIVER_COUNT; i++ {
		select {
		case p.stopChans[i] <- true:
		case <-time.After(time.Second * 30):
		}
		close(p.stopChans[i])
	}
	p.stopChans = make([]chan bool, RECEIVER_COUNT)
}

func (p *MQSQueue) ReceiveMessage(respChan chan MessageReceiveResponse, errChan chan error, waitseconds ...int64) {
	resource := fmt.Sprintf("%s/%s", p.name, "messages")
	if waitseconds != nil && len(waitseconds) == 1 {
		resource = fmt.Sprintf("%s/%s?waitseconds=%d", p.name, "messages", waitseconds[0])
	}

	//mqs's http pool is active by send while no message exist, so more sender will get back fast
	//ali-mqs bug:	error code of 499, while client disconnet the request,
	//				the mqs server did not drop the sleeping recv connection,
	//				so the others recv connection could not recv message
	//				until the sleeping recv released

	var wg sync.WaitGroup

	funcSend := func(respChan chan MessageReceiveResponse, errChan chan error, stopChan chan bool) {
		defer wg.Done()
		for {
			resp := MessageReceiveResponse{}
			_, err := p.client.Send(GET, nil, nil, resource, &resp)
			if err != nil {
				errChan <- err
			} else {
				respChan <- resp
			}

			select {
			case _ = <-stopChan:
				{
					return
				}
			default:
			}
		}
	}

	for i := 0; i < RECEIVER_COUNT; i++ {
		wg.Add(1)
		stopChan := make(chan bool)
		p.stopChans[i] = stopChan
		go funcSend(respChan, errChan, stopChan)
	}

	wg.Wait()

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
