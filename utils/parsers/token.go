/*
Package parsers - ZeWise 解析器包
该文件用于实现令牌解析器
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package parsers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"zewise.space/backend/consts"
)

var (
	// ErrNoBearerToken 未提供 Bearer Token
	ErrNoBearerToken = errors.New("未提供 Bearer Token")
	// ErrInvalidBearerTokenFormat Bearer Token 格式错误
	ErrInvalidBearerTokenFormat = errors.New("Bearer Token 格式错误")
)

// BeaerTokenClaims Bearer Token 声明
type BearerTokenClaims struct {
	jwt.RegisteredClaims        // JWT 注册声明
	UID                  string `json:"uid"`      // 用户 ID
	UserName             string `json:"username"` // 用户名
}

/*
GetUID 获取用户 ID

返回：
  - primitive.ObjectID：用户 ID
  - error：错误信息
*/
func (claims *BearerTokenClaims) GetUserObjectID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(claims.UID)
}

/*
ParseContextTokenString 解析上下文中的 Token 字符串

参数：
  - ctx：上下文

返回：
  - string：Token 字符串
*/
func ParseContextTokenString(ctx *fiber.Ctx) (string, error) {
	// 从请求头中获取 Token
	token := ctx.Get("Authorization")
	if token == "" {
		return "", ErrNoBearerToken
	}

	// 检查 Token 是否为 Bearer Token
	if len(token) < 7 || token[:7] != "Bearer " {
		return "", ErrInvalidBearerTokenFormat
	}

	// 提取 Token
	return token[7:], nil
}

/*
ParseToken 解析令牌

参数：
  - token：令牌字符串

返回：
  - BearerTokenClaims：令牌中的声明
  - error：错误信息
*/
func ParseToken(token string) (BearerTokenClaims, error) {
	// 解析令牌
	claims := BearerTokenClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.TOKEN_SECRET), nil
	})

	// 返回结果
	return claims, err
}
