/*
Package services - ZeWise 服务层
该文件用于声明用户相关服务
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package services

import (
	"bytes"
	"context"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"zewise.space/backend/consts"
	"zewise.space/backend/models"
	"zewise.space/backend/stores"
	"zewise.space/backend/types"
	"zewise.space/backend/utils/encryptors"
	"zewise.space/backend/utils/functools"
	"zewise.space/backend/utils/generators"
	"zewise.space/backend/utils/imagetools"
	"zewise.space/backend/utils/parsers"
	"zewise.space/backend/utils/validers"
)

// UserService 用户服务
type UserService struct {
	Storage *stores.Storage
}

/*
RegisterUser 注册用户

参数：
  - username：用户名
  - email：邮箱
  - password：密码

返回：
  - error：错误信息
*/
func (service *UserService) RegisterUser(username string, email string, password string) error {
	// 验证用户名 密码 邮箱是否合法
	if !validers.IsValidEmail(email) {
		return types.NewError(types.ErrInvalidParams, "不合法的邮箱")
	}
	if !validers.IsValidUsername(username) {
		return types.NewError(types.ErrInvalidParams, "不合法的用户名")
	}
	if !validers.IsValidPassword(password) {
		return types.NewError(types.ErrInvalidParams, "不合法的密码")
	}

	// 创建数据库会话
	ctx := context.Background()
	session, err := service.Storage.NewSession()
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}
	defer session.EndSession(ctx)

	// 开启事务
	_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (any, error) {
		// 检验用户名是否重复
		err := service.Storage.UserStorage.CheckUserExistance(sessionContext, username, email)
		if err != nil {
			return nil, err
		}

		// 生成盐和哈希密码
		salt, err := generators.GenerateSalt(consts.SALT_LENGTH)
		if err != nil {
			return nil, err
		}
		hashedPassword, err := encryptors.HashPassword(password, salt)
		if err != nil {
			return nil, err
		}

		// 注册用户
		err = service.Storage.UserStorage.RegisterUser(sessionContext, username, email, salt, hashedPassword)
		return nil, err
	})

	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}

/*
GetUserProfileByID 获取用户信息

参数：
  - userID：用户ID

返回：
  - models.UserInfo：用户信息
  - error：错误信息
*/
func (service *UserService) GetUserProfileByID(userID string) (models.UserInfo, error) {
	userInfo := models.UserInfo{}

	// 创建数据库会话
	ctx := context.Background()
	session, err := service.Storage.NewSession()
	if err != nil {
		return userInfo, types.NewError(types.ErrServerError, err.Error())
	}
	defer session.EndSession(ctx)

	// 转换用户ID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return userInfo, types.NewError(types.ErrInvalidParams, "不合法的用户ID")
	}

	// 开启事务
	_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (any, error) {
		// 获取用户信息
		userInfo, err = service.Storage.UserStorage.GetUserDataByID(sessionContext, objID)
		return nil, err
	})
	if err != nil {
		return userInfo, err
	}

	return userInfo, nil
}

/*
GetUserProfileByUsername 获取用户信息

参数：
  - username：用户名

返回：
  - models.UserInfo：用户信息
  - error：错误信息
*/
func (service *UserService) GetUserProfileByUsername(username string) (models.UserInfo, error) {
	userInfo := models.UserInfo{}

	// 创建数据库会话
	ctx := context.Background()
	session, err := service.Storage.NewSession()
	if err != nil {
		return userInfo, types.NewError(types.ErrServerError, err.Error())
	}
	defer session.EndSession(ctx)

	// 开启事务
	_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (any, error) {
		// 获取用户信息
		userInfo, err = service.Storage.UserStorage.GetUserDataByUsername(sessionContext, username)
		return nil, err
	})
	if err != nil {
		return userInfo, err
	}

	return userInfo, nil
}

/*
UpdateUserProfile 更新用户信息

参数：
  - reqBody：请求体

返回：
  - error：错误信息
*/
func (service *UserService) UpdateUserProfile(ID primitive.ObjectID, reqBody parsers.UserUpdateProfileBody) error {
	// 构造用户信息
	userInfo := models.UserInfo{
		ID:       ID,
		NickName: reqBody.NickName,
		Sign:     reqBody.Sign,
		Birth:    time.Unix(reqBody.Birth, 0),
		Gender:   reqBody.Gender,
	}

	// 创建数据库会话
	ctx := context.Background()
	session, err := service.Storage.NewSession()
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}
	defer session.EndSession(ctx)

	// 开启事务
	_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (any, error) {
		// 更新用户信息
		err := service.Storage.UserStorage.UpdateUserProfile(sessionContext, userInfo)
		return nil, err
	})
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}

/*
UpdateUserAvatar 更新用户头像

参数：
  - userID：用户ID
  - avatarFile：头像文件

返回：
  - error：错误信息
*/
func (service *UserService) UpdateUserAvatar(userID primitive.ObjectID, avatarFileHeader *multipart.FileHeader) error {
	// 处理头像文件
	decoder := imagetools.NewDefaultImageDecoderChain()
	decoder.SetContentType(avatarFileHeader.Header.Get("Content-Type"))
	sizeLimiter := imagetools.NewSizeLimiter(consts.AVATAR_MAX_SIZE, consts.AVATAR_MAX_SIZE)
	resizer := imagetools.NewResizeProcessHandler(consts.AVATAR_SIZE, consts.AVATAR_SIZE, &imagetools.ScallingDownProcessor{})
	encoder := imagetools.NewWebpImageEncoder(consts.AVATAR_QUALITY)

	avatarFile, err := avatarFileHeader.Open()
	if err != nil {
		return types.NewError(types.ErrInvalidParams, "不合法的头像文件")
	}
	defer avatarFile.Close()

	imageData, err := imagetools.ProcessImage(avatarFile, decoder, encoder, sizeLimiter, resizer)
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	// 创建数据库会话
	ctx := context.Background()
	session, err := service.Storage.NewSession()
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}
	defer session.EndSession(ctx)

	// 开启事务
	_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (any, error) {
		// 获取用户信息
		userInfo, err := service.Storage.UserStorage.GetUserDataByID(sessionContext, userID)
		if err != nil {
			return nil, err
		}
		if userInfo.Avatar != "vanilla" {
			// 删除原头像
			err = service.Storage.UserStorage.DeleteAvatarFile(
				context.Background(),
				functools.JoinStrings(userInfo.Avatar, ".", encoder.GetFormatFileSuffix()),
			)
			// 如果错误存在且不是文件不存在错误
			if err != nil && minio.ToErrorResponse(err).Code != "NoSuchKey" {
				return nil, err
			}
		}

		// 更新用户头像
		_, err = service.Storage.UserStorage.UploadAvatarFile(
			context.Background(),
			functools.JoinStrings(userID.Hex(), ".", encoder.GetFormatFileSuffix()),
			bytes.NewReader(imageData),
		)
		if err != nil {
			return nil, err
		}

		// 更新用户信息
		err = service.Storage.UserStorage.UpdateUserProfile(sessionContext, models.UserInfo{
			ID:     userID,
			Avatar: userID.Hex(),
		})
		return nil, err
	})

	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}

/*
UpdateUserPassword 更新用户密码

参数：
  - userID：用户ID
  - oldPassword：旧密码
  - newPassword：新密码

返回：
  - error：错误信息
*/
func (service *UserService) UpdateUserPassword(userID primitive.ObjectID, oldPassword string, newPassword string) error {
	// 创建数据库会话
	ctx := context.Background()
	session, err := service.Storage.NewSession()
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}
	defer session.EndSession(ctx)

	// 开启事务
	_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (any, error) {
		// 获取用户信息
		authInfo, err := service.Storage.AuthStorage.GetUserAuthInfoByID(sessionContext, userID)
		if err != nil {
			return nil, err
		}

		// 验证密码
		err = encryptors.CompareHashPassword(authInfo.PasswordHash, oldPassword, authInfo.Salt)
		if err != nil {
			return nil, types.NewError(types.ErrInvalidParams, "旧密码错误")
		}

		// 生成新哈希密码
		hashedPassword, err := encryptors.HashPassword(newPassword, authInfo.Salt)
		if err != nil {
			return nil, err
		}

		// 更新用户密码
		err = service.Storage.UserStorage.UpdateUserPassword(sessionContext, userID, hashedPassword)
		return nil, err
	})

	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}
