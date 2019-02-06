package errorStatus

type AppErrorRt interface {
	error
	State() int
}

type AppError struct {
	Code int   `json:"code"`
	Err  error `json:"error"`
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) State() int {
	return e.Code
}
