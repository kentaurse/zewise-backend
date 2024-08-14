/*
Package consts - ZeWise 常量包
该文件用于声明令牌相关常量
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package consts

const (
	// TOKEN_EXPIRE_DURATION 有效期
	TOKEN_EXPIRE_DURATION = 7 * 24 * 60 * 60

	// TOKEN_SECRET 令牌密钥
	TOKEN_SECRET = "ZEWISE_BACKEND_EXAMPLE_SECRET"

	// TOKEN_ISSUER 令牌签发者
	TOKEN_ISSUER = "space.zewise.auth"

	// MAX_TOKENS_PER_USER 最大令牌数量
	MAX_TOKENS_PER_USER = 5
)
