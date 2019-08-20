package errors

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var digitsRegexp = regexp.MustCompile(`:\-?(0|[1-9]\d{0,}).+?:\\?"([\S\s]+)\\?"`)

type customError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// New 初始化一个错误体
func New(code int, msg string) error {
	return &customError{code, msg}
}

func (e *customError) Error() string {
	return fmt.Sprintf(`{"code":%d,"message":"%s"}`, e.Code, e.Message)
}

// Code 获取错误码
func Code(err error) int {
	cus, ok := err.(*customError)
	if ok {
		return cus.Code
	}
	return -99999
}

// Message 获取错误信息
func Message(err error) string {
	cus, ok := err.(*customError)
	if ok {
		return cus.Message
	}
	return err.Error()
}

// CodeWithString 获取自定义错误的错误码
func CodeWithString(err string) int {
	e := &customError{}
	if json.Unmarshal([]byte(err), e) != nil {
		return -99999
	}
	return e.Code
}

// MessageWithString 获取自定义错误的错误信息
func MessageWithString(err string) string {
	e := &customError{}
	if json.Unmarshal([]byte(err), e) != nil {
		return err
	}
	return e.Message
}

func parseError(msg string) (ok bool, err error) {
	msgs := digitsRegexp.FindStringSubmatch(msg)
	if len(msgs) < 3 {
		return false, &customError{
			Code:    -1,
			Message: msg,
		}
	}
	code, _ := strconv.ParseInt(msgs[1], 10, 64)
	return true, &customError{
		Code:    int(code),
		Message: strings.TrimSpace(msgs[2]),
	}
}

// IsCustomerErr 判断是否自定义错误信息
func IsCustomerErr(err error) bool {
	_, ok := err.(*customError)
	return ok
}
