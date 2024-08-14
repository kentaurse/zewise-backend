/*
Package controllers - ZeWise 控制器
该文件用于声明认证接口控制器
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mssola/useragent"
	"zewise.space/backend/services"
	"zewise.space/backend/types"
	"zewise.space/backend/utils/parsers"
	"zewise.space/backend/utils/serializers"
)

// AuthController 认证控制器
type AuthController struct {
	service *services.Service // 服务对象
}

/*
NewAuthController 新建认证控制器

返回：
  - *AuthController：认证控制器对象
*/
func (factory *Factory) NewAuthController() *AuthController {
	return &AuthController{factory.service}
}

/*
NewLoginHandler 新建登录接口处理函数

返回：
  - fiber.Handler：Fiber 处理函数
*/
func (controller *AuthController) NewLoginHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody, err := parsers.ParseBody[parsers.UserLoginBody](ctx)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, err.Error())),
			)
		}

		// 校验参数
		if (reqBody.Email == "" && reqBody.UserName == "") || reqBody.Password == "" {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, "邮箱和密码不能为空")),
			)
		}

		// 登录
		token, err := controller.service.AuthService.AuthLogin(
			reqBody.Email,
			reqBody.UserName,
			reqBody.Password,
			ctx.IP(),
			useragent.New(ctx.Get("User-Agent")),
		)

		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, err.Error())),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(serializers.SUCCESS, "", serializers.NewAuthLoginResponse(token)),
		)
	}
}

/*
NewLogoutHandler 新建退出接口处理函数

返回：
  - fiber.Handler：Fiber 处理函数
*/
func (controller *AuthController) NewLogoutHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取 Token
		token, _ := parsers.ParseContextTokenString(ctx)

		// 获取 TokenClaims
		tokenClaims := ctx.Locals("claims").(parsers.BearerTokenClaims)

		// 登出
		err := controller.service.AuthService.AuthLogout(tokenClaims, token)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(err),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(serializers.SUCCESS, ""),
		)
	}
}

/*
NewRefreshTokenHandler 新建刷新 Token 接口处理函数

返回：
  - fiber.Handler：Fiber 处理函数
*/
func (controller *AuthController) NewRefreshTokenHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取 Token
		token, err := parsers.ParseContextTokenString(ctx)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, err.Error())),
			)
		}

		// 解析 Token
		claims, err := parsers.ParseToken(token)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrAuthFailed, "bearer token 无效")),
			)
		}

		userObjectID, err := claims.GetUserObjectID()
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, "不合法的用户ID")),
			)
		}

		// 刷新 Token
		newToken, err := controller.service.AuthService.RefreshToken(userObjectID, claims.UserName, token)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(err),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(serializers.SUCCESS, "", serializers.NewAuthLoginResponse(newToken)),
		)
	}
}
