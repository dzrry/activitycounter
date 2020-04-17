package api

import (
	"fmt"
	"strings"
)

type ServerError int

const (
	// Full description at https://vk.com/dev/errors
	ErrBadResponseCode ServerError = -1
	ErrBadCode         ServerError = -666
)

type Errors []ExecuteError

func (e Errors) Error() string {
	var s []string
	for _, v := range e {
		s = append(s, v.Error())
	}
	return fmt.Sprintln("Execute errors:", strings.Join(s, ", "))
}

type RequestParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExecuteError struct {
	Method  string      `json:"method"`
	Code    ServerError `json:"error_code"`
	Message string      `json:"error_msg"`
}

// Error contains standard errors.
type Error struct {
	Code       ServerError     `json:"error_code,omitempty"`
	Message    string          `json:"error_msg,omitempty"`
	Params     *[]RequestParam `json:"request_params,omitempty"`
	CaptchaSid string          `json:"captcha_sid,omitempty"`
	CaptchaImg string          `json:"captcha_img,omitempty"`
	Request    Request         `json:"-"`
}

func (e Error) ToError() error {
	return fmt.Errorf("%s (%d)", e.Message, e.Code)
}

// NewError makes *Error from our ServerError and description.
func NewError(code ServerError, description string) (err *Error) {
	err = new(Error)
	err.Code = code
	err.Message = description

	return
}

// setRequest sets Request
func (e *Error) setRequest(r Request) { // nolint: unused
	e.Request = r
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

func (e ExecuteError) Error() string {
	return fmt.Sprintf("%s: %s (%d)", e.Method, e.Message, e.Code)
}

func (s ServerError) String() string {
	return fmt.Sprintf("%d", s)
}

func (s ServerError) Error() string {
	return s.String()
}
