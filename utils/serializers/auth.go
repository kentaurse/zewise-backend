/*
Package serializers - ZeWise 序列化器包
该文件用于序列化认证信息
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package serializers

// AuthLoginResponse 认证登录响应
type AuthLoginResponse struct {
	Token string `json:"token,omitempty"` // Token
}

/*
NewAuthLoginResponse 创建认证登录响应

参数：
  - token：Token

返回：
  - AuthLoginResponse：认证登录响应
*/
func NewAuthLoginResponse(token string) AuthLoginResponse {
	return AuthLoginResponse{
		Token: token,
	}
}
