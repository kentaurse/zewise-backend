/*
Package serializers - ZeWise 序列化器包
该文件用于序列化错误情况下的返回信息
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package serializers

import (
	"errors"

	"github.com/spf13/viper"
	"zewise.space/backend/types"
)

/*
NewErrorResponse 创建错误响应

参数：
  - err：错误

返回：
  - any：序列化后的错误信息
*/
func NewErrorResponse(raw_error error) any {
	// 获取环境类型
	env := viper.GetString("env.type")

	// 提取错误信息
	var code ResponseCode
	var message string

	// 类型断言
	err, ok := raw_error.(types.Error)
	if !ok {
		return NewResponse(UNKNOWN_ERROR, err.Error())
	}

	// 判断错误类型
	if errors.Is(err, types.ErrInvalidParams) {
		code = PARAMETER_ERROR
	}
	if errors.Is(err, types.ErrAuthFailed) {
		code = AUTH_ERROR
	}
	if errors.Is(err, types.ErrNetworkError) {
		code = NETWORK_ERROR
	}
	if errors.Is(err, types.ErrServerError) {
		code = SERVER_ERROR
	}

	// 生产环境下隐藏服务器错误信息
	if env == "production" && code == SERVER_ERROR {
		message = ""
	} else {
		message = err.ErrMessage
	}

	return NewResponse(code, message)
}
