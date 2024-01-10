package errors2

type ErrorInfo struct {
	Err error
	Code int
	Msg string
	SubMsg string
	ExtraData interface{}
}
