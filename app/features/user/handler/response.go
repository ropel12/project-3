package handler

type WebResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func CreateWebResponse(code int, message string, data any) any {
	return WebResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
