/*
Package parsers - ZeWise 解析器包
该文件用于解析请求体
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package parsers

import (
	"github.com/gofiber/fiber/v2"
	"zewise.space/backend/types"
)

/*
ParseBody 解析请求体

参数：
  - T：请求体类型
  - ctx：Fiber 上下文

返回：
  - T：请求体
  - error：错误信息
*/
func ParseBody[T any](ctx *fiber.Ctx) (T, error) {
    reqBody := new(T)
	err := ctx.BodyParser(reqBody)
	if err != nil {
		return *reqBody, types.NewError(types.ErrInvalidParams, err.Error())
	}
	return *reqBody, nil
}
