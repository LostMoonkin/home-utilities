package response

type APIResponse struct {
	BizCode int64  `json:"biz_code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(data any) *APIResponse {
	return NewAPIResponse(data, BizSuccess, "")
}

func Fail(code int64, msg string) *APIResponse {
	return &APIResponse{
		BizCode: code,
		Message: msg,
	}
}

func NewAPIResponse(data any, code int64, message string) *APIResponse {
	return &APIResponse{
		code, message, data,
	}
}
