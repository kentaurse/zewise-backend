/*
Package parsers - ZeWise 解析器包
该文件声明了用户相关的解析结构
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package parsers

// UserLoginBody 用户登录请求体
type UserRegisterBody struct {
	Email    string `json:"email"`    // 邮箱
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// UserUpdateProfileBody 用户更新资料请求体
type UserUpdateProfileBody struct {
	NickName string `json:"nickname"` // 用户名
	Sign     string `json:"sign"`     // 签名
	Birth    int64  `json:"birth"`    // 生日
	Gender   string `json:"gender"`   // 性别
}

// UserUpdatePasswordBody 更新密码请求体
type UserUpdatePasswordBody struct {
	OldPassword string `json:"old_password"` // 旧密码
	NewPassword string `json:"new_password"` // 新密码
}
