/*
Package generators - ZeWise 后端服务器生成器包
该文件用于生成 Token
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package generators

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"zewise.space/backend/consts"
	"zewise.space/backend/utils/parsers"
)

/*
GenerateToken 生成 Token

参数：
  - uid：用户 ID
  - username：用户名

返回：
  - string：Token
  - *BearerTokenClaims：Token Claims
  - error：错误信息
*/
func GenerateToken(uid primitive.ObjectID, username string) (string, parsers.BearerTokenClaims, error) {
	// 构造 Token 的 Claims
	claims := parsers.BearerTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.TOKEN_EXPIRE_DURATION * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    consts.TOKEN_ISSUER,
			Subject:   "BearerToken",
			ID:        uuid.New().String(),
		},
		UID:      uid.Hex(),
		UserName: username,
	}

	// 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 Token
	tokenString, err := token.SignedString([]byte(consts.TOKEN_SECRET))

	// 返回 Token 和 Claims
	return tokenString, claims, err
}
