package errors

type Error interface {
	GetCode() (code int)
	Error() (message string)
}
