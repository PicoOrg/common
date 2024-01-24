package errors

func New(code int, message string) Error {
	return &implement{
		Code:    code,
		Message: message,
	}
}

type implement struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *implement) GetCode() int {
	return e.Code
}

func (e *implement) Error() string {
	return e.Message
}
