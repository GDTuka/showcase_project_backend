package error_service

type IAppError interface {
	Error() string
	Code() int
}

type AppError struct {
	Err     error
	ErrCode int
}

func (e AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

func (e AppError) Code() int {
	return e.ErrCode
}

func NewAppError(err error, code int) IAppError {
	return AppError{Err: err, ErrCode: code}
}
