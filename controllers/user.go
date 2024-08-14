/*
Package controllers - ZeWise 控制器
该文件用于声明用户接口控制器
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package controllers

import (
	"github.com/gofiber/fiber/v2"

	"zewise.space/backend/models"
	"zewise.space/backend/services"
	"zewise.space/backend/types"
	"zewise.space/backend/utils/parsers"
	"zewise.space/backend/utils/serializers"
)

// UserController 用户控制器
type UserController struct {
	service *services.Service // 服务对象
}

/*
NewUserController 新建用户控制器

返回：
  - *UserController：用户控制器对象
*/
func (factory *Factory) NewUserController() *UserController {
	return &UserController{factory.service}
}

/*
NewRegisterHandler 新建注册接口处理函数

返回：
  - fiber.Handler：Fiber 处理函数
*/
func (controller *UserController) NewRegisterHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody, err := parsers.ParseBody[parsers.UserRegisterBody](ctx)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, err.Error())),
			)
		}

		// 校验参数
		if reqBody.Username == "" || reqBody.Email == "" || reqBody.Password == "" {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, "用户名、密码和邮箱不能为空")),
			)
		}

		// 注册用户
		err = controller.service.UserService.RegisterUser(reqBody.Username, reqBody.Email, reqBody.Password)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(err),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(serializers.SUCCESS, ""),
		)
	}
}

/*
NewProfileHandler 新建用户信息接口处理函数

返回：
  - fiber.Handler：Fiber 处理函数
*/
func (controller *UserController) NewProfileHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取请求参数
		userID := ctx.Query("id")
		username := ctx.Query("username")

		if userID == "" && username == "" {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, "需要提供用户ID或用户名")),
			)
		}

		// 获取用户信息
		var userInfo models.UserInfo
		var err error
		if userID != "" {
			userInfo, err = controller.service.UserService.GetUserProfileByID(userID)
		} else {
			userInfo, err = controller.service.UserService.GetUserProfileByUsername(username)
		}
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(err),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(serializers.SUCCESS, "", serializers.NewUserProfileResponse(userInfo)),
		)
	}
}

/*
NewUpdateProfileHandler 新建更新用户信息接口处理函数

返回：
  - fiber.Handler：Fiber 处理函数
*/
func (controller *UserController) NewUpdateProfileHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody, err := parsers.ParseBody[parsers.UserUpdateProfileBody](ctx)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, err.Error())),
			)
		}

		// 获取用户ID
		claims := ctx.Locals("claims").(parsers.BearerTokenClaims)
		userID, err := claims.GetUserObjectID()
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, "不合法的用户ID")),
			)
		}

		// 更新用户信息
		err = controller.service.UserService.UpdateUserProfile(userID, reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(err),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(serializers.SUCCESS, ""),
		)
	}
}

/*
NewUpdateAvatarHandler 新建更新用户头像接口处理函数

返回：
  - fiber.Handler：Fiber 处理函数
*/
func (controller *UserController) NewUpdateAvatarHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取用户ID
		claims := ctx.Locals("claims").(parsers.BearerTokenClaims)
		userID, err := claims.GetUserObjectID()
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, "不合法的用户ID")),
			)
		}

		// 获取表单文件
		form, err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, "无法获取头像文件")),
			)
		}
		files := form.File["avatar"]
		if len(files) == 0 || files[0] == nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, "头像文件不能为空")),
			)
		}
		if len(files) > 1 {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, "只能上传一个头像文件")),
			)
		}
		fileHeader := files[0]

		// 更新用户头像
		err = controller.service.UserService.UpdateUserAvatar(userID, fileHeader)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(err),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(serializers.SUCCESS, ""),
		)
	}
}

/*
NewUpdatePasswordHandler 新建更新用户密码接口处理函数

返回：
  - fiber.Handler：Fiber 处理函数
*/
func (controller *UserController) NewUpdatePasswordHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取用户ID
		claims := ctx.Locals("claims").(parsers.BearerTokenClaims)
		userID, err := claims.GetUserObjectID()
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, "不合法的用户ID")),
			)
		}

		// 解析请求体
		reqBody, err := parsers.ParseBody[parsers.UserUpdatePasswordBody](ctx)

		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(types.NewError(types.ErrInvalidParams, err.Error())),
			)
		}

		// 更新用户密码
		err = controller.service.UserService.UpdateUserPassword(userID, reqBody.OldPassword, reqBody.NewPassword)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewErrorResponse(err),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(serializers.SUCCESS, ""),
		)
	}
}
