package ali_mqs

import (
	"github.com/gogap/errors"
)

const (
	ALI_MQS_ERR_NS = "MQS"

	ali_MQS_ERR_TEMPSTR = "ali_mqs response status error,code: {{.resp.Code}}, message: {{.resp.Message}}, request id: {{.resp.RequestId}}, host id: {{.resp.HostId}}"
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
)
