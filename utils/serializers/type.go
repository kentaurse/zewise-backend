/*
Package serializers - ZeWise 序列化器包
该文件用于定义序列化相关类型
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package serializers

// ResponseCode 响应码
type ResponseCode int

const (
	// SUCCESS 成功
	SUCCESS ResponseCode = 0

	// SERVER_ERROR 服务器错误
	SERVER_ERROR ResponseCode = 1

	// PARAMETER_ERROR 参数错误
	PARAMETER_ERROR ResponseCode = 2

	// AUTH_ERROR 认证错误
	AUTH_ERROR ResponseCode = 3

	// NETWORK_ERROR 网络错误
	NETWORK_ERROR ResponseCode = 4

	// UNKNOWN_ERROR 未知错误
	UNKNOWN_ERROR ResponseCode = -1
)
