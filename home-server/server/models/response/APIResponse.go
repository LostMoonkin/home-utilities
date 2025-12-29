package response

type APIResponse[T any] struct {
	BizCode int64  `json:"biz_code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func Success[T any](data T) *APIResponse[T] {
	return NewAPIResponse(data, BizSuccess, "")
}

func NewAPIResponse[T any](data T, code int64, message string) *APIResponse[T] {
	return &APIResponse[T]{
		code, message, data,
	}
}
