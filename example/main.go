package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gogap/ali_mqs"
	"github.com/gogap/logs"
)

type appConf struct {
	Url             string `json:"url"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

func main() {
	conf := appConf{}

	if bFile, e := ioutil.ReadFile("app.conf"); e != nil {
		panic(e)
	} else {
		if e := json.Unmarshal(bFile, &conf); e != nil {
			panic(e)
		}
	}

	credential := ali_mqs.NewAliMQSCredential()
	client := ali_mqs.NewAliMQSClient(conf.Url,
		conf.AccessKeyId,
		conf.AccessKeySecret,
		credential)

	msg := ali_mqs.MessageSendRequest{
		MessageBody:  []byte("hello gogap/ali_mqs"),
		DelaySeconds: 0,
		Priority:     8}

	queue := ali_mqs.NewMQSQueue("test", client)
	ret, err := queue.SendMessage(msg)

	if err != nil {
		logs.Error(err)
	} else {
		logs.Pretty("response:", ret)
	}

	respChan := make(chan ali_mqs.MessageReceiveResponse)
	errChan := make(chan error)
	go func() {
		for {
			select {
			case resp := <-respChan:
				{
					logs.Pretty("response:", resp)
					logs.Debug("change the visibility: ", resp.ReceiptHandle)
					if ret, e := queue.ChangeMessageVisibility(resp.ReceiptHandle, 5); e != nil {
						logs.Error(e)
					} else {
						logs.Pretty("visibility changed", ret)
					}

					logs.Debug("delete it now: ", resp.ReceiptHandle)
					if e := queue.DeleteMessage(resp.ReceiptHandle); e != nil {
						logs.Error(e)
					}
				}
			case err := <-errChan:
				{
					logs.Error(err)
				}
			}
		}

	}()

	queue.ReceiveMessage(respChan, errChan)
	for {
		time.Sleep(time.Second * 1)
	}

}
