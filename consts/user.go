/*
Package consts - ZeWise 常量包
该文件用于定义密码相关常量
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package consts

const (
	// EMAIL_REGEX 邮箱正则表达式
	EMAIL_REGEX = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// USERNAME_MIN_LENGTH 用户名最小长度
	USERNAME_MIN_LENGTH = 4

	// USERNAME_MAX_LENGTH 用户名最大长度
	USERNAME_MAX_LENGTH = 24

	// USERNAME_REGEX 用户名正则表达式
	USERNAME_REGEX = `^[a-z0-9_]+$`

	// PASSWORD_MIN_LENGTH 密码最小长度
	PASSWORD_MIN_LENGTH = 8

	// PASSWORD_MAX_LENGTH 密码最大长度
	PASSWORD_MAX_LENGTH = 32

	// PASSWORD_REGEX 密码正则表达式
	PASSWORD_REGEX = `^[a-zA-Z0-9!@#$%^&*()_+={}\[\]:;'"<>,.?\/|\\~-]+$`
)