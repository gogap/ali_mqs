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
	version = "2015-06-06"
)

const (
	DefaultTimeout int64 = 35
)

type Method string

var (
	errMapping map[string]errors.ErrCodeTemplate
)

func init() {

}

const (
	_GET    Method = "GET"
	_PUT           = "PUT"
	_POST          = "POST"
	_DELETE        = "DELETE"
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

func initMQSErrors() {
	errMapping = map[string]errors.ErrCodeTemplate{
		"AccessDenied":               ERR_MQS_ACCESS_DENIED,
		"InvalidAccessKeyId":         ERR_MQS_INVALID_ACCESS_KEY_ID,
		"InternalError":              ERR_MQS_INTERNAL_ERROR,
		"InvalidAuthorizationHeader": ERR_MQS_INVALID_AUTHORIZATION_HEADER,
		"InvalidDateHeader":          ERR_MQS_INVALID_DATE_HEADER,
		"InvalidArgument":            ERR_MQS_INVALID_ARGUMENT,
		"InvalidDegist":              ERR_MQS_INVALID_DEGIST,
		"InvalidRequestURL":          ERR_MQS_INVALID_REQUEST_URL,
		"InvalidQueryString":         ERR_MQS_INVALID_QUERY_STRING,
		"MalformedXML":               ERR_MQS_MALFORMED_XML,
		"MissingAuthorizationHeader": ERR_MQS_MISSING_AUTHORIZATION_HEADER,
		"MissingDateHeader":          ERR_MQS_MISSING_DATE_HEADER,
		"MissingVersionHeader":       ERR_MQS_MISSING_VERSION_HEADER,
		"MissingReceiptHandle":       ERR_MQS_MISSING_RECEIPT_HANDLE,
		"MissingVisibilityTimeout":   ERR_MQS_MISSING_VISIBILITY_TIMEOUT,
		"MessageNotExist":            ERR_MQS_MESSAGE_NOT_EXIST,
		"QueueAlreadyExist":          ERR_MQS_QUEUE_ALREADY_EXIST,
		"QueueDeletedRecently":       ERR_MQS_QUEUE_DELETED_RECENTLY,
		"InvalidQueueName":           ERR_MQS_INVALID_QUEUE_NAME,
		"QueueNameLengthError":       ERR_MQS_QUEUE_NAME_LENGTH_ERROR,
		"QueueNotExist":              ERR_MQS_QUEUE_NOT_EXIST,
		"ReceiptHandleError":         ERR_MQS_RECEIPT_HANDLE_ERROR,
		"SignatureDoesNotMatch":      ERR_MQS_SIGNATURE_DOES_NOT_MATCH,
		"TimeExpired":                ERR_MQS_TIME_EXPIRED,
		"QpsLimitExceeded":           ERR_MQS_QPS_LIMIT_EXCEEDED,
	}
}

func to_error(resp ErrorMessageResponse, resource string) (err error) {
	if errCodeTemplate, exist := errMapping[resp.Code]; exist {
		err = errCodeTemplate.New(errors.Params{"resp": resp, "resource": resource})
	} else {
		err = ERR_MQS_UNKNOWN_CODE.New(errors.Params{"resp": resp, "resource": resource})
	}
	return
}
