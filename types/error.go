/*
Package types ZeWise 类型
该文件用于定义 ZeWise 后端服务的错误类型
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package types

import (
	"errors"
	"fmt"
)

// ErrorType 错误类型
type ErrorType error

var (
	// ErrServerError 服务器错误
	ErrServerError ErrorType = errors.New("ServerError")

	// ErrInvalidParams 参数错误
	ErrInvalidParams ErrorType = errors.New("InvalidParams")

	// ErrAuthFailed 认证错误
	ErrAuthFailed ErrorType = errors.New("AuthFailed")

	// ErrNetworkError 网络错误
	ErrNetworkError ErrorType = errors.New("NetworkError")

	// ErrUnknownError 未知错误
	ErrUnknownError ErrorType = errors.New("UnknownError")
)

// Error 错误
type Error struct {
	ErrType    ErrorType // 错误类型
	ErrMessage string    // 错误信息
}

/*
Error 错误信息

返回：
  - string：错误信息
*/
func (err Error) Error() string {
	return fmt.Sprintf("%v: %s", err.ErrType, err.ErrMessage)
}

/*
Is 判断错误类型是否相同

参数：
  - another：另一个错误

返回：
  - bool：是否相同
*/
func (err Error) Is(another error) bool {
	return errors.Is(err.ErrType, another)
}

/*
NewError 新建错误

参数：
  - errType：错误类型
  - errMessage：错误信息

返回：
  - Error：错误
*/
func NewError(errType ErrorType, errMessage string) Error {
	return Error{
		ErrType:    errType,
		ErrMessage: errMessage,
	}
}
