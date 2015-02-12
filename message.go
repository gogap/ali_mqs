package ali_mqs

import (
	"encoding/base64"
	"encoding/xml"
)

type MessageResponse struct {
	XMLName   xml.Name `xml:"Message" json:"-"`
	Code      string   `xml:"Code,omitempty" json:"code,omitempty"`
	Message   string   `xml:"Message,omitempty" json:"message,omitempty"`
	RequestId string   `xml:"RequestId,omitempty" json:"request_id,omitempty"`
	HostId    string   `xml:"HostId,omitempty" json:"host_id,omitempty"`
}

type MessageSendRequest struct {
	XMLName      xml.Name `xml:"Message"`
	MessageBody  []byte   `xml:"MessageBody"`
	DelaySeconds int64    `xml:"DelaySeconds"`
	Priority     int64    `xml:"Priority"`
}

type MessageSendResponse struct {
	MessageResponse
	MessageId      string `xml:"MessageId" json:"message_id"`
	MessageBodyMD5 string `xml:"MessageBodyMD5" json:"message_body_md5"`
}

type MessageReceiveResponse struct {
	MessageResponse
	ReceiptHandle    string `xml:"ReceiptHandle" json:"receipt_handle"`
	MessageBodyMD5   string `xml:"MessageBodyMD5" json:"message_body_md5"`
	MessageBody      []byte `xml:"MessageBody" json:"message_body"`
	EnqueueTime      int64  `xml:"EnqueueTime" json:"enqueue_time"`
	NextVisibleTime  int64  `xml:"NextVisibleTime" json:"next_visible_time"`
	FirstDequeueTime int64  `xml:"FirstDequeueTime" json:"first_dequeue_time"`
	DequeueCount     int64  `xml:"DequeueCount" json:"dequeue_count"`
	Priority         int64  `xml:"Priority" json:"priority"`
}

func (p *MessageReceiveResponse) DecodeBody() (body []byte, err error) {
	buf := make([]byte, len(p.MessageBody))
	if length, e := base64.StdEncoding.Decode(buf, p.MessageBody); e != nil {
		return nil, e
	} else {
		return buf[0:length], nil
	}
}

type MessageVisibilityChangeResponse struct {
	XMLName         xml.Name `xml:"ChangeVisibility" json:"-"`
	ReceiptHandle   string   `xml:"ReceiptHandle" json:"receipt_handle"`
	NextVisibleTime int64    `xml:"NextVisibleTime" json:"next_visible_time"`
}
