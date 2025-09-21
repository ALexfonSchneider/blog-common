package errors

import (
	"errors"
	"fmt"
)

// ApplicationError — ошибка приложения
type ApplicationError struct {
	code    int    // внутренний код (1000, 1001...)
	message string // Человеко-читаемое сообщение
	detail  string // Дополнительная информация
	cause   error  // Исходная ошибка
}

func New(code int, message string, detail string) *ApplicationError {
	return &ApplicationError{code: code, message: message, detail: detail}
}

func (e *ApplicationError) Error() string {
	return e.String()
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
	return fmt.Sprintf(
		"AppError{code: %d, Message: %q, Detail: %q, Cause: %v}", e.code, e.message, e.detail, e.cause,
	)
}

func (e *ApplicationError) Code() int {
	return e.code
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
	if targetErr, ok := target.(*ApplicationError); ok {
		return e.code == targetErr.code
	}
	return errors.Is(e.cause, target)
}
