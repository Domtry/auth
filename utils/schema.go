package utils

type HttpResponse[T any] struct {
	Message   string `json:"message"`
	Success   bool   `json:"success"`
	CodeError int    `json:"code_error"`
	Data      T      `json:"data"`
}

func (h HttpResponse[T]) ToJson() interface{} {
	return map[string]interface{}{
		"message":    h.Message,
		"success":    h.Success,
		"code_error": h.CodeError,
		"data":       h.Data,
	}
}
