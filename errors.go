package ali_mqs

import (
	"github.com/gogap/errors"
)

const (
	ALI_MQS_ERR_NS = "MQS"

	ali_MQS_ERR_TEMPSTR = "ali_mqs response status error,code: {{.resp.Code}}, message: {{.resp.Message}}, resource: {{.resource}} request id: {{.resp.RequestId}}, host id: {{.resp.HostId}}"
)

var (
	ERR_SIGN_MESSAGE_FAILED        = errors.TN(ALI_MQS_ERR_NS, 1, "sign message failed, {{.err}}")
	ERR_MARSHAL_MESSAGE_FAILED     = errors.TN(ALI_MQS_ERR_NS, 2, "marshal message filed, {{.err}}")
	ERR_GENERAL_AUTH_HEADER_FAILED = errors.TN(ALI_MQS_ERR_NS, 3, "general auth header failed, {{.err}}")

	ERR_CREATE_NEW_REQUEST_FAILED = errors.TN(ALI_MQS_ERR_NS, 4, "create new request failed, {{.err}}")
	ERR_SEND_REQUEST_FAILED       = errors.TN(ALI_MQS_ERR_NS, 5, "send request failed, {{.err}}")
	ERR_READ_RESPONSE_BODY_FAILED = errors.TN(ALI_MQS_ERR_NS, 6, "read response body failed, {{.err}}")

	ERR_UNMARSHAL_ERROR_RESPONSE_FAILED = errors.TN(ALI_MQS_ERR_NS, 7, "unmarshal error response failed, {{.err}}")
	ERR_UNMARSHAL_RESPONSE_FAILED       = errors.TN(ALI_MQS_ERR_NS, 8, "unmarshal response failed, {{.err}}")
	ERR_DECODE_BODY_FAILED              = errors.TN(ALI_MQS_ERR_NS, 9, "decode body failed, {{.err}}, body: \"{{.body}}\"")
	ERR_GET_BODY_DECODE_ELEMENT_ERROR   = errors.TN(ALI_MQS_ERR_NS, 10, "get body decode element error, local: {{.local}}, error: {{.err}}")

	ERR_MQS_ACCESS_DENIED                = errors.TN(ALI_MQS_ERR_NS, 100, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INVALID_ACCESS_KEY_ID        = errors.TN(ALI_MQS_ERR_NS, 101, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INTERNAL_ERROR               = errors.TN(ALI_MQS_ERR_NS, 102, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INVALID_AUTHORIZATION_HEADER = errors.TN(ALI_MQS_ERR_NS, 103, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INVALID_DATE_HEADER          = errors.TN(ALI_MQS_ERR_NS, 104, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INVALID_ARGUMENT             = errors.TN(ALI_MQS_ERR_NS, 105, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INVALID_DEGIST               = errors.TN(ALI_MQS_ERR_NS, 106, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INVALID_REQUEST_URL          = errors.TN(ALI_MQS_ERR_NS, 107, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INVALID_QUERY_STRING         = errors.TN(ALI_MQS_ERR_NS, 108, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_MALFORMED_XML                = errors.TN(ALI_MQS_ERR_NS, 109, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_MISSING_AUTHORIZATION_HEADER = errors.TN(ALI_MQS_ERR_NS, 110, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_MISSING_DATE_HEADER          = errors.TN(ALI_MQS_ERR_NS, 111, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_MISSING_VERSION_HEADER       = errors.TN(ALI_MQS_ERR_NS, 112, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_MISSING_RECEIPT_HANDLE       = errors.TN(ALI_MQS_ERR_NS, 113, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_MISSING_VISIBILITY_TIMEOUT   = errors.TN(ALI_MQS_ERR_NS, 114, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_MESSAGE_NOT_EXIST            = errors.TN(ALI_MQS_ERR_NS, 115, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_QUEUE_ALREADY_EXIST          = errors.TN(ALI_MQS_ERR_NS, 116, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_QUEUE_DELETED_RECENTLY       = errors.TN(ALI_MQS_ERR_NS, 117, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INVALID_QUEUE_NAME           = errors.TN(ALI_MQS_ERR_NS, 118, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INVALID_VERSION_HEADER       = errors.TN(ALI_MQS_ERR_NS, 119, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_INVALID_CONTENT_TYPE         = errors.TN(ALI_MQS_ERR_NS, 120, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_QUEUE_NAME_LENGTH_ERROR      = errors.TN(ALI_MQS_ERR_NS, 121, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_QUEUE_NOT_EXIST              = errors.TN(ALI_MQS_ERR_NS, 122, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_RECEIPT_HANDLE_ERROR         = errors.TN(ALI_MQS_ERR_NS, 123, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_SIGNATURE_DOES_NOT_MATCH     = errors.TN(ALI_MQS_ERR_NS, 124, ali_MQS_ERR_TEMPSTR)
	ERR_MQS_TIME_EXPIRED                 = errors.TN(ALI_MQS_ERR_NS, 125, ali_MQS_ERR_TEMPSTR)

	ERR_MQS_QUEUE_NAME_IS_TOO_LONG                 = errors.TN(ALI_MQS_ERR_NS, 126, "queue name is too long, the max length is 256")
	ERR_MQS_DELAY_SECONDS_RANGE_ERROR              = errors.TN(ALI_MQS_ERR_NS, 127, "queue delay seconds is not in range of (0~60480)")
	ERR_MQS_MAX_MESSAGE_SIZE_RANGE_ERROR           = errors.TN(ALI_MQS_ERR_NS, 128, "max message size is not in range of (1024~65536)")
	ERR_MQS_MSG_RETENTION_PERIOD_RANGE_ERROR       = errors.TN(ALI_MQS_ERR_NS, 129, "message retention period is not in range of (60~129600)")
	ERR_MQS_MSG_VISIBILITY_TIMEOUT_RANGE_ERROR     = errors.TN(ALI_MQS_ERR_NS, 130, "message visibility timeout is not in range of (1~43200)")
	ERR_MQS_MSG_POOLLING_WAIT_SECONDS_RANGE_ERROR  = errors.TN(ALI_MQS_ERR_NS, 131, "message poolling wait seconds is not in range of (0~30)")
	REE_MQS_GET_QUEUE_RET_NUMBER_RANGE_ERROR       = errors.TN(ALI_MQS_ERR_NS, 132, "get queue list param of ret number is not in range of (1~1000)")
	ERR_MQS_QUEUE_ALREADY_EXIST_AND_HAVE_SAME_ATTR = errors.TN(ALI_MQS_ERR_NS, 133, "mqs queue already exist, and the attribute is the same, queue name: {{.name}}")
)
