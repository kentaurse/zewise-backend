/*
Package middlewares - ZeWise 后端服务器中间件。
该文件用于定义 Token 认证中间件。
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"zewise.space/backend/stores"
	"zewise.space/backend/types"
	"zewise.space/backend/utils/parsers"
	"zewise.space/backend/utils/serializers"
)

// TokenAuthMiddleware 认证中间件
type TokenAuthMiddleware struct {
	authStorage *stores.AuthStorage
}

/*
NewTokenAuthMiddleware 新建 Token 认证中间件

返回：
  - *TokenAuthMiddleware：Token 认证中间件对象
*/
func (factory *Factory) NewTokenAuthMiddleware() *TokenAuthMiddleware {
	return &TokenAuthMiddleware{factory.storage.AuthStorage}
}

/*
NewMiddleware Token 认证中间件

参数：
  - ctx：Fiber 上下文。

返回：
  - error：错误
*/
func (middleware *TokenAuthMiddleware) NewMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 从请求头中获取 Token
		token, err := parsers.ParseContextTokenString(ctx)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, err.Error())),
			)
		}

		// 验证 Token
		claims, err := parsers.ParseToken(token)

		// 处理 Token 错误
		// Token 过期
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrAuthFailed, "bearer token 已过期")),
			)
		}
		// Token 无效
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrAuthFailed, "bearer token 无效")),
			)
		}

		// 检验 Token 是否可用
		isAvaliable, err := middleware.authStorage.CheckTokenAvailability(claims.UID, token)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(err),
			)
		}
		if !isAvaliable {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrAuthFailed, "bearer token 已失效")),
			)
		}

		// 将 claims 信息存入 ctx.Locals 中
		ctx.Locals("claims", claims)

		return ctx.Next()
	}
}
