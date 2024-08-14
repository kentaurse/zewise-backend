/*
Package stores - NekoBlog 后端服务器数据访问层
该文件用于实现认证信息存储对象类
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package stores

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"zewise.space/backend/models"
	"zewise.space/backend/types"
	"zewise.space/backend/utils/functools"

	"github.com/redis/go-redis/v9"
)

// AuthStorage 用户信息数据库
type AuthStorage struct {
	redis *redis.Client
	mongo *mongo.Database
}

/*
GetUserAuthInfoByID 通过 ID 获取用户认证信息

参数：
  - sessionContext：数据库会话上下文
  - id：用户 ID

返回：
  - models.UserAuthInfo：用户认证信息
  - error：错误信息
*/
func (store *AuthStorage) GetUserAuthInfoByID(sessionContext mongo.SessionContext, id primitive.ObjectID) (models.UserAuthInfo, error) {
	userAuthInfo := models.UserAuthInfo{ID: id}
	err := store.mongo.Collection(models.USER_AUTH_INFO_COLLECTION).FindOne(sessionContext, userAuthInfo).Decode(&userAuthInfo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return userAuthInfo, types.NewError(types.ErrInvalidParams, "用户不存在")
		}
		return userAuthInfo, types.NewError(types.ErrServerError, err.Error())
	}

	return userAuthInfo, nil
}

/*
GetUserAuthInfoByEmail 通过邮箱获取用户认证信息

参数：
  - sessionContext：数据库会话上下文
  - email：邮箱

返回：
  - models.UserAuthInfo：用户认证信息
  - error：错误信息
*/
func (store *AuthStorage) GetUserAuthInfoByEmail(sessionContext mongo.SessionContext, email string) (models.UserAuthInfo, error) {
	userAuthInfo := models.UserAuthInfo{Email: email}
	err := store.mongo.Collection(models.USER_AUTH_INFO_COLLECTION).FindOne(sessionContext, userAuthInfo).Decode(&userAuthInfo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return userAuthInfo, types.NewError(types.ErrInvalidParams, "用户不存在")
		}
		return userAuthInfo, types.NewError(types.ErrServerError, err.Error())
	}

	return userAuthInfo, nil
}

/*
GetUserAuthInfoByUsername 通过用户名获取用户认证信息

参数：
  - sessionContext：数据库会话上下文
  - username：用户名

返回：
  - models.UserAuthInfo：用户认证信息
  - error：错误信息
*/
func (store *AuthStorage) GetUserAuthInfoByUsername(sessionContext mongo.SessionContext, username string) (models.UserAuthInfo, error) {
	userAuthInfo := models.UserAuthInfo{UserName: username}
	err := store.mongo.Collection(models.USER_AUTH_INFO_COLLECTION).FindOne(sessionContext, userAuthInfo).Decode(&userAuthInfo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return userAuthInfo, types.NewError(types.ErrInvalidParams, "用户不存在")
		}
		return userAuthInfo, types.NewError(types.ErrServerError, err.Error())
	}

	return userAuthInfo, nil
}

/*
CheckTokenAvailability 检查令牌是否可用

参数：
  - token：令牌

返回：
  - bool：令牌是否可用
  - error：错误信息
*/
func (store *AuthStorage) CheckTokenAvailability(userID string, token string) (bool, error) {
	// 获取用户令牌列表
	list, err := store.GetAvailableToken(userID)
	if err != nil {
		return false, types.NewError(types.ErrServerError, err.Error())
	}
	if len(list) == 0 {
		return false, nil
	}

	// 检查令牌是否在列表中
	for _, t := range list {
		if t == token {
			return true, nil
		}
	}

	return false, nil
}

/*
GetAvailableToken 获取可用令牌

参数：
  - userID：用户 ID

返回：
  - []token：令牌列表
  - error：错误信息
*/
func (store *AuthStorage) GetAvailableToken(userID string) ([]string, error) {
	ctx := context.Background()

	// 获取用户令牌列表
	list, err := store.redis.LRange(
		ctx, functools.JoinStrings(models.REDIS_AVAILABLE_USER_TOKEN_LIST, ":", userID), 0, -1,
	).Result()
	if err != nil {
		return nil, types.NewError(types.ErrServerError, err.Error())
	}

	return list, nil
}

/*
RemoveEarliestToken 移除最早的令牌

参数：
  - userID：用户 ID
  - count：移除数量

返回：
  - error：错误信息
*/
func (store *AuthStorage) RemoveEarliestToken(userID string, count int) error {
	ctx := context.Background()

	// 移除最早的令牌
	_, err := store.redis.LTrim(
		ctx, functools.JoinStrings(models.REDIS_AVAILABLE_USER_TOKEN_LIST, ":", userID), int64(count), -1,
	).Result()
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}

/*
SaveToken 保存令牌

参数：
  - userID：用户 ID
  - token：令牌

返回：
  - error：错误信息
*/
func (store *AuthStorage) SaveToken(userID string, token string) error {
	ctx := context.Background()

	// 保存令牌
	_, err := store.redis.RPush(
		ctx, functools.JoinStrings(models.REDIS_AVAILABLE_USER_TOKEN_LIST, ":", userID), token,
	).Result()
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}

/*
UpdateToken 更新令牌

参数：
  - userID：用户 ID
  - token：令牌

返回：
  - error：错误信息
*/
func (store *AuthStorage) UpdateToken(userID string, newToken string, oldToken string) error {
	ctx := context.Background()

	// 查找旧令牌索引
	tokens, err := store.GetAvailableToken(userID)
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	index := -1
	for i, t := range tokens {
		if t == oldToken {
			index = i
			break
		}
	}
	if index == -1 {
		return types.NewError(types.ErrServerError, "未找到旧令牌")
	}

	// 更新令牌
	_, err = store.redis.LSet(
		ctx, functools.JoinStrings(models.REDIS_AVAILABLE_USER_TOKEN_LIST, ":", userID), int64(index), newToken,
	).Result()
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}

/*
RmoveToken 移除令牌

参数：
  - userID：用户 ID
  - token：令牌

返回：
  - error：错误信息
*/
func (store *AuthStorage) RmoveToken(userID string, token string) error {
	ctx := context.Background()

	// 移除令牌
	_, err := store.redis.LRem(
		ctx, functools.JoinStrings(models.REDIS_AVAILABLE_USER_TOKEN_LIST, ":", userID), 0, token,
	).Result()
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}

/*
RecordLoginEvent 记录登录事件

参数：
  - userID：用户 ID
  - ip：IP 地址
  - Location：位置
  - Device：设备
  - Application：应用

返回：
  - error：错误信息
*/
func (store *AuthStorage) RecordLoginEvent(sessionContext mongo.SessionContext, userID primitive.ObjectID, ip string, location string, device string, application string) error {
	// 创建登录日志
	loginLog := models.UserLoginLog{
		UID:         userID,
		IP:          ip,
		Location:    location,
		Device:      device,
		Time:        time.Now(),
		Application: application,
		IfSucceed:   false,
		IfChecked:   false,
	}

	// 保存登录日志
	_, err := store.mongo.Collection(models.USER_LOGIN_LOGS_COLLECTION).InsertOne(sessionContext, loginLog)
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}
