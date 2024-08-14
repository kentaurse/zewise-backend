/*
Package services - ZeWise 服务层
该文件用于声明认证相关服务
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package services

import (
	"context"
	"os"

	"github.com/mssola/useragent"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"zewise.space/backend/consts"
	"zewise.space/backend/models"
	"zewise.space/backend/stores"
	"zewise.space/backend/types"
	"zewise.space/backend/utils/encryptors"
	"zewise.space/backend/utils/functools"
	"zewise.space/backend/utils/generators"
	"zewise.space/backend/utils/parsers"
	"zewise.space/backend/utils/thirdparty"
)

// AuthService 认证服务
type AuthService struct {
	Storage *stores.Storage
}

/*
AuthLogin 用户登录

参数：
  - email：邮箱
  - username：用户名
  - password：密码

返回：
  - token：JWT 令牌
  - error：错误信息
*/
func (service *AuthService) AuthLogin(email string, username string, password string, ip string, userAgent *useragent.UserAgent) (string, error) {
	var token = ""

	// 创建数据库会话
	ctx := context.Background()
	session, err := service.Storage.NewSession()
	if err != nil {
		return "", types.NewError(types.ErrServerError, err.Error())
	}
	defer session.EndSession(ctx)

	// 开启事务
	_, err = session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		var userAuthInfo models.UserAuthInfo
		// 获取用户信息
		if email != "" {
			userAuthInfo, err = service.Storage.AuthStorage.GetUserAuthInfoByEmail(sessionCtx, email)
		} else {
			userAuthInfo, err = service.Storage.AuthStorage.GetUserAuthInfoByUsername(sessionCtx, username)
		}
		if err != nil {
			return nil, err
		}

		// 校验密码
		err = encryptors.CompareHashPassword(userAuthInfo.PasswordHash, password, userAuthInfo.Salt)
		if err != nil {
			return nil, types.NewError(types.ErrInvalidParams, "邮箱或密码错误")
		}

		// 生成 JWT 令牌
		token, _, err = generators.GenerateToken(userAuthInfo.ID, userAuthInfo.UserName)
		if err != nil {
			return nil, types.NewError(types.ErrServerError, err.Error())
		}

		// 获取令牌数量
		tokens, err := service.Storage.AuthStorage.GetAvailableToken(userAuthInfo.ID.Hex())
		if err != nil {
			return nil, err
		}

		// 超过最大令牌数量
		if len(tokens) >= consts.MAX_TOKENS_PER_USER {
			err = service.Storage.AuthStorage.RemoveEarliestToken(userAuthInfo.ID.Hex(), 1)
			if err != nil {
				return nil, err
			}
		}

		// 保存令牌
		err = service.Storage.AuthStorage.SaveToken(userAuthInfo.ID.Hex(), token)
		if err != nil {
			return nil, err
		}

		// 获取 IP 对应地理位置
		ipInfo, _ := thirdparty.RequestIPInfo(ip, os.Getenv("AMAP_KEY"))

		// 获取浏览器信息
		broswer, broswerVersion := userAgent.Browser()

		// 写入登录日志
		err = service.Storage.AuthStorage.RecordLoginEvent(
			sessionCtx,
			userAuthInfo.ID,
			ip,
			functools.JoinStrings(ipInfo.Province, ipInfo.City),
			userAgent.OSInfo().FullName,
			functools.JoinStrings(broswer, " ", broswerVersion),
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return "", types.NewError(types.ErrServerError, err.Error())
	}

	return token, nil
}

/*
AuthLogout 用户登出

参数：
  - token：JWT 令牌

返回：
  - error：错误信息
*/
func (service *AuthService) AuthLogout(tokenClaim parsers.BearerTokenClaims, token string) error {
	return service.Storage.AuthStorage.RmoveToken(tokenClaim.UID, token)
}

/*
RefreshToken 刷新 Token

参数：
  - userID：用户 ID

返回：
  - newToken：新 Token
  - error：错误信息
*/
func (service *AuthService) RefreshToken(userID primitive.ObjectID, username string, oldToken string) (string, error) {
	// 生成 JWT 令牌
	newToken, _, err := generators.GenerateToken(userID, username)
	if err != nil {
		return "", types.NewError(types.ErrServerError, err.Error())
	}

	// 替换新令牌
	err = service.Storage.AuthStorage.UpdateToken(userID.Hex(), newToken, oldToken)
	if err != nil {
		return "", err
	}

	return newToken, nil
}
