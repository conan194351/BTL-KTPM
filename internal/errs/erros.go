package errs

type Error struct {
	Message        string `json:"Message"`
	HttpStatusCode int    `json:"HttpStatusCode"`
	ErrorCode      string `json:"ErrorCode"`
}

type ErrorCode int

//go:generate go run github.com/dmarkham/enumer -type=ErrorCode

const (
	ErrNone ErrorCode = iota
	ErrBadRequest
	ErrUnauthorized
	ErrForbidden
	ErrNotFound
	ErrMethodNotAllowed
	ErrRequestTimeout
	ErrConflict
	ErrInternalServer
	ErrServiceUnavailable
	ErrGatewayTimeout
	ErrMalformedJSON
	ErrEmbedding
	ErrChatCompletion
)

type ErrorCodeMap map[ErrorCode]Error

var ErrorCodes = ErrorCodeMap{
	ErrNone:               {HttpStatusCode: 0, Message: "Success", ErrorCode: ErrNone.String()},
	ErrBadRequest:         {HttpStatusCode: 400, Message: "Bad Request", ErrorCode: ErrBadRequest.String()},
	ErrUnauthorized:       {HttpStatusCode: 401, Message: "Unauthorized", ErrorCode: ErrUnauthorized.String()},
	ErrForbidden:          {HttpStatusCode: 403, Message: "Forbidden", ErrorCode: ErrForbidden.String()},
	ErrNotFound:           {HttpStatusCode: 404, Message: "Not Found", ErrorCode: ErrNotFound.String()},
	ErrMethodNotAllowed:   {HttpStatusCode: 405, Message: "Method Not Allowed", ErrorCode: ErrMethodNotAllowed.String()},
	ErrRequestTimeout:     {HttpStatusCode: 408, Message: "Request Timeout", ErrorCode: ErrRequestTimeout.String()},
	ErrConflict:           {HttpStatusCode: 409, Message: "Conflict", ErrorCode: ErrConflict.String()},
	ErrInternalServer:     {HttpStatusCode: 500, Message: "Internal Server Error", ErrorCode: ErrInternalServer.String()},
	ErrServiceUnavailable: {HttpStatusCode: 503, Message: "Service Unavailable", ErrorCode: ErrServiceUnavailable.String()},
	ErrGatewayTimeout:     {HttpStatusCode: 504, Message: "Gateway Timeout", ErrorCode: ErrGatewayTimeout.String()},
	ErrMalformedJSON:      {HttpStatusCode: 400, Message: "Malformed JSON", ErrorCode: ErrMalformedJSON.String()},
	ErrEmbedding:          {HttpStatusCode: 500, Message: "Embedding is not supported", ErrorCode: ErrEmbedding.String()},
	ErrChatCompletion:     {HttpStatusCode: 500, Message: "Chat completion error", ErrorCode: ErrChatCompletion.String()},
}

func ResponseError(code ErrorCode) Error {
	return ErrorCodes[code]
}
