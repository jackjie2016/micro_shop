package defs

type Err struct {
	Error string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErrorResponse struct {
	HttpSe int
	Error Err
}

var (
	ErrorRequestBodyParseFailed = ErrorResponse{HttpSe:400,Error:Err{Error:"Response body is not correct",ErrorCode:"001"}}
    ErrorNotAuthUser = ErrorResponse{HttpSe:401,Error:Err{Error:"user authentication",ErrorCode:"002"}}
)















