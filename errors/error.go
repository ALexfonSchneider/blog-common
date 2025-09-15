package errors

import (
	"fmt"
)

// ApplicationError — ошибка приложения
type ApplicationError struct {
	code     int    // внутренний код (1000, 1001...)
	httpCode int    // http-code альтернатива ошибки
	message  string // человеко-читаемое сообщение
	detail   string // дополнительная информация
	cause    error  // исходная ошибка
}

func New(code int, httpCode int, message string, detail string, cause error) *ApplicationError {
	return &ApplicationError{code, httpCode, message, detail, cause}
}

func (e *ApplicationError) Error() string {
	if e.detail == "" {
		return e.message
	}
	return fmt.Sprintf("%s: %s", e.message, e.detail)
}

func (e *ApplicationError) Unwrap() error {
	return e.cause
}

func (e *ApplicationError) WithCause(cause error) *ApplicationError {
	err := *e
	err.cause = cause
	return &err
}

func (e *ApplicationError) Wrap(cause error) *ApplicationError {
	return e.WithCause(cause)
}

func (e *ApplicationError) String() string {
	return fmt.Sprintf("AppError{code: %d, Message: %q, Detail: %q, Cause: %v}",
		e.code, e.message, e.detail, e.cause)
}

func (e *ApplicationError) Code() int {
	return e.code
}

func (e *ApplicationError) HttpCode() int {
	return e.httpCode
}

func (e *ApplicationError) Message() string {
	return e.message
}

func (e *ApplicationError) Detail() string {
	return e.detail
}

func (e *ApplicationError) Cause() error {
	return e.cause
}

func (e *ApplicationError) Is(target error) bool {
	switch v := target.(type) {
	case *ApplicationError:
		return e.code == v.code
	default:
		return e.String() == target.Error()
	}
}
