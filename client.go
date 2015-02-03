package ali_mqs

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mreiferson/go-httpclient"
)

const (
	version = "2014-07-08"
)

const (
	DefaultTimeout int64 = 35000
)

type Method string

const (
	GET    Method = "GET"
	PUT           = "PUT"
	POST          = "POST"
	DELETE        = "DELETE"
)

type MQSClient interface {
	Send(method Method, headers map[string]string, message interface{}, resource string, v interface{}) (err error)
}

type AliMQSClient struct {
	Timeout     int64
	url         string
	credential  Credential
	accessKeyId string
}

func NewAliMQSClient(url, accessKeyId, accessKeySecret string, credential Credential) *AliMQSClient {
	if url == "" {
		panic("ali-mqs: message queue url is empty")
	}

	if credential == nil {
		panic("ali-mqs: client credential is nil")
	}

	aliMQSClient := new(AliMQSClient)
	aliMQSClient.credential = credential
	aliMQSClient.accessKeyId = accessKeyId
	aliMQSClient.url = url
	credential.SetSecretKey(accessKeySecret)
	return aliMQSClient
}

func (p *AliMQSClient) authorization(method Method, headers map[string]string, resource string) (authHeader string, err error) {
	if signature, e := p.credential.Signature(method, headers, resource); e != nil {
		err = e
		return
	} else {
		authHeader = fmt.Sprintf("MQS %s:%s", p.accessKeyId, signature)
	}

	return
}

func (p *AliMQSClient) Send(method Method, headers map[string]string, message interface{}, resource string, v interface{}) (err error) {
	var xmlContent []byte
	if message == nil {
		xmlContent = []byte{}
	} else {
		if bXml, e := xml.Marshal(message); e != nil {
			err = e
			return
		} else {
			xmlContent = bXml
		}
	}

	xmlMD5 := md5.Sum(xmlContent)
	strMd5 := fmt.Sprintf("%x", xmlMD5)

	if headers == nil {
		headers = make(map[string]string)
	}

	headers[MQ_VERSION] = version
	headers[CONTENT_TYPE] = "application/xml"
	headers[CONTENT_MD5] = base64.StdEncoding.EncodeToString([]byte(strMd5))
	headers[DATE] = time.Now().UTC().Format(http.TimeFormat)

	if authHeader, e := p.authorization(method, headers, fmt.Sprintf("/%s", resource)); e != nil {
		err = e
		return
	} else {
		headers[AUTHORIZATION] = authHeader
	}

	timeoutInt := DefaultTimeout
	url := p.url + "/" + resource

	if p.Timeout > 0 {
		timeoutInt = p.Timeout
	}

	timeout := time.Second * time.Duration(timeoutInt)

	transport := &httpclient.Transport{
		ConnectTimeout:        timeout,
		RequestTimeout:        timeout,
		ResponseHeaderTimeout: timeout,
	}

	defer transport.Close()

	postBodyReader := strings.NewReader(string(xmlContent))

	client := &http.Client{Transport: transport}

	var req *http.Request
	if req, err = http.NewRequest(string(method), url, postBodyReader); err != nil {
		return
	}

	for header, value := range headers {
		req.Header.Add(header, value)
	}

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	} else if resp != nil {
		if bBody, e := ioutil.ReadAll(resp.Body); e != nil {
			err = e
			return
		} else if resp.StatusCode != http.StatusCreated &&
			resp.StatusCode != http.StatusOK &&
			resp.StatusCode != http.StatusNoContent {
			err = fmt.Errorf("ali-mqs: response status code is %d, body: %s", resp.StatusCode, string(bBody))
		} else if v != nil {
			if e := xml.Unmarshal(bBody, v); e != nil {
				err = e
				return
			}
		}
	}
	return
}
