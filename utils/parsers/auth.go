/*
Package parsers - ZeWise 解析器包
该文件声明了认证相关的解析结构
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package parsers

// UserLoginBody 用户登录请求体
type UserLoginBody struct {
	Email    string `json:"email"`    // 邮箱
	UserName string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}
