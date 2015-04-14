package ali_mqs

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gogap/errors"
	"github.com/mreiferson/go-httpclient"
)

const (
	version = "2014-07-08"
)

const (
	DefaultTimeout int64 = 35
)

type Method string

const (
	GET    Method = "GET"
	PUT           = "PUT"
	POST          = "POST"
	DELETE        = "DELETE"
)

type MQSClient interface {
	Send(method Method, headers map[string]string, message interface{}, resource string, v interface{}) (statusCode int, err error)
	SetProxy(url string)
}

type AliMQSClient struct {
	Timeout      int64
	url          string
	credential   Credential
	accessKeyId  string
	clientLocker sync.Mutex
	client       *http.Client
	proxyURL     string
}

func NewAliMQSClient(url, accessKeyId, accessKeySecret string) MQSClient {
	if url == "" {
		panic("ali-mqs: message queue url is empty")
	}

	credential := NewAliMQSCredential(accessKeySecret)

	aliMQSClient := new(AliMQSClient)
	aliMQSClient.credential = credential
	aliMQSClient.accessKeyId = accessKeyId
	aliMQSClient.url = url

	timeoutInt := DefaultTimeout

	if aliMQSClient.Timeout > 0 {
		timeoutInt = aliMQSClient.Timeout
	}

	timeout := time.Second * time.Duration(timeoutInt)

	transport := &httpclient.Transport{
		Proxy:                 aliMQSClient.proxy,
		ConnectTimeout:        time.Second * 3,
		RequestTimeout:        timeout,
		ResponseHeaderTimeout: timeout + time.Second,
	}

	aliMQSClient.client = &http.Client{Transport: transport}

	return aliMQSClient
}

func (p *AliMQSClient) SetProxy(url string) {
	p.url = url
}

func (p *AliMQSClient) proxy(req *http.Request) (*url.URL, error) {
	if p.url != "" {
		return url.Parse(p.url)
	}
	return nil, nil
}

func (p *AliMQSClient) authorization(method Method, headers map[string]string, resource string) (authHeader string, err error) {
	if signature, e := p.credential.Signature(method, headers, resource); e != nil {
		return "", e
	} else {
		authHeader = fmt.Sprintf("MQS %s:%s", p.accessKeyId, signature)
	}

	return
}

func (p *AliMQSClient) Send(method Method, headers map[string]string, message interface{}, resource string, v interface{}) (statusCode int, err error) {
	var xmlContent []byte

	if message == nil {
		xmlContent = []byte{}
	} else {
		if bXml, e := xml.Marshal(message); e != nil {
			err = ERR_MARSHAL_MESSAGE_FAILED.New(errors.Params{"err": e})
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
		err = ERR_GENERAL_AUTH_HEADER_FAILED.New(errors.Params{"err": e})
		return
	} else {
		headers[AUTHORIZATION] = authHeader
	}

	url := p.url + "/" + resource

	postBodyReader := strings.NewReader(string(xmlContent))

	p.clientLocker.Lock()
	defer p.clientLocker.Unlock()

	var req *http.Request
	if req, err = http.NewRequest(string(method), url, postBodyReader); err != nil {
		err = ERR_CREATE_NEW_REQUEST_FAILED.New(errors.Params{"err": err})
		return
	}

	for header, value := range headers {
		req.Header.Add(header, value)
	}

	var resp *http.Response
	if resp, err = p.client.Do(req); err != nil {
		err = ERR_SEND_REQUEST_FAILED.New(errors.Params{"err": err})
		return
	}

	if resp != nil {
		defer resp.Body.Close()
		statusCode = resp.StatusCode
		if bBody, e := ioutil.ReadAll(resp.Body); e != nil {
			err = ERR_READ_RESPONSE_BODY_FAILED.New(errors.Params{"err": e})
			return
		} else if resp.StatusCode != http.StatusCreated &&
			resp.StatusCode != http.StatusOK &&
			resp.StatusCode != http.StatusNoContent {
			errResp := ErrorMessageResponse{}
			if e := xml.Unmarshal(bBody, &errResp); e != nil {
				err = ERR_UNMARSHAL_ERROR_RESPONSE_FAILED.New(errors.Params{"err": e})
				return
			}
			err = to_error(errResp, resource)
			return
		} else if v != nil {
			if e := xml.Unmarshal(bBody, v); e != nil {
				err = ERR_UNMARSHAL_RESPONSE_FAILED.New(errors.Params{"err": e})
				return
			}
		}
	}
	return
}

func to_error(resp ErrorMessageResponse, resource string) (err error) {
	switch resp.Code {
	case "AccessDenied":
		{
			err = ERR_MQS_ACCESS_DENIED.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InvalidAccessKeyId":
		{
			err = ERR_MQS_INVALID_ACCESS_KEY_ID.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InternalError":
		{
			err = ERR_MQS_INTERNAL_ERROR.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InvalidAuthorizationHeader":
		{
			err = ERR_MQS_INVALID_AUTHORIZATION_HEADER.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InvalidDateHeader":
		{
			err = ERR_MQS_INVALID_DATE_HEADER.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InvalidArgument":
		{
			err = ERR_MQS_INVALID_ARGUMENT.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InvalidDegist":
		{
			err = ERR_MQS_INVALID_DEGIST.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InvalidRequestURL":
		{
			err = ERR_MQS_INVALID_REQUEST_URL.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InvalidQueryString":
		{
			err = ERR_MQS_INVALID_QUERY_STRING.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "MalformedXML":
		{
			err = ERR_MQS_MALFORMED_XML.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "MissingAuthorizationHeader":
		{
			err = ERR_MQS_MISSING_AUTHORIZATION_HEADER.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "MissingDateHeader":
		{
			err = ERR_MQS_MISSING_DATE_HEADER.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "MissingVersionHeader":
		{
			err = ERR_MQS_MISSING_VERSION_HEADER.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "MissingReceiptHandle":
		{
			err = ERR_MQS_MISSING_RECEIPT_HANDLE.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "MissingVisibilityTimeout":
		{
			err = ERR_MQS_MISSING_VISIBILITY_TIMEOUT.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "MessageNotExist":
		{
			err = ERR_MQS_MESSAGE_NOT_EXIST.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "QueueAlreadyExist":
		{
			err = ERR_MQS_QUEUE_ALREADY_EXIST.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "QueueDeletedRecently":
		{
			err = ERR_MQS_QUEUE_DELETED_RECENTLY.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InvalidQueueName":
		{
			err = ERR_MQS_INVALID_QUEUE_NAME.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InvalidVersionHeader":
		{
			err = ERR_MQS_INVALID_VERSION_HEADER.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "InvalidContentType":
		{
			err = ERR_MQS_INVALID_CONTENT_TYPE.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "QueueNameLengthError":
		{
			err = ERR_MQS_QUEUE_NAME_LENGTH_ERROR.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "QueueNotExist":
		{
			err = ERR_MQS_QUEUE_NOT_EXIST.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "ReceiptHandleError":
		{
			err = ERR_MQS_RECEIPT_HANDLE_ERROR.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "SignatureDoesNotMatch":
		{
			err = ERR_MQS_SIGNATURE_DOES_NOT_MATCH.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	case "TimeExpired":
		{
			err = ERR_MQS_TIME_EXPIRED.New(errors.Params{"resp": resp, "resource": resource})
			return
		}
	}
	return
}
