/*
Package stores - NekoBlog 后端服务器数据访问层。
该文件用于声明用户存储对象类。
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package stores

import (
	"context"
	"errors"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"

	"zewise.space/backend/models"
	"zewise.space/backend/types"
)

// UserStorage 用户信息数据库
type UserStorage struct {
	redis *redis.Client
	mongo *mongo.Database
	minio *minio.Client
}

/*
CheckUserExistance 检查用户是否存在

参数：
  - sessionContext：数据库会话上下文
  - username：用户名
  - email：邮箱

返回：
  - error：错误信息
*/
func (store *UserStorage) CheckUserExistance(sessionContext mongo.SessionContext, username string, email string) error {
	// 检查用户名和邮箱是否已被注册
	count, err := store.mongo.Collection(models.USER_INFO_COLLECTION).CountDocuments(sessionContext, bson.M{
		"$or": bson.A{
			bson.M{"username": username},
			bson.M{"email": email},
		},
	})
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}
	if count > 0 {
		return types.NewError(types.ErrInvalidParams, "用户名或邮箱已被注册")
	}

	return nil
}

/*
RegisterUser 注册用户

参数：
  - sessionContext：数据库会话上下文
  - username：用户名
  - salt：盐
  - hashedPassword：哈希密码

返回：
  - error：错误信息
*/
func (store *UserStorage) RegisterUser(sessionContext mongo.SessionContext, username string, email string, salt string, hashedPassword string) error {
	// 插入用户信息
	user := models.UserInfo{
		UserName:  username,
		NickName:  username,
		Email:     email,
		Avatar:    "vanilla",
		Sign:      "这个人很懒，什么都没有留下。",
		Authority: 0,
		Level:     1,
	}
	result, err := store.mongo.Collection(models.USER_INFO_COLLECTION).InsertOne(sessionContext, user)
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	// 获取用户ID
	userID := result.InsertedID.(primitive.ObjectID)

	// 插入用户认证信息
	userAuthInfo := models.UserAuthInfo{
		ID:           userID,
		UserName:     username,
		Email:        email,
		Salt:         salt,
		PasswordHash: hashedPassword,
	}
	_, err = store.mongo.Collection(models.USER_AUTH_INFO_COLLECTION).InsertOne(sessionContext, userAuthInfo)
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}

/*
GetUserDataByID 通过用户ID获取用户信息

参数：
  - sessionContext：数据库会话上下文
  - userID：用户ID

返回：
  - models.UserInfo：用户信息
  - error：错误信息
*/
func (store *UserStorage) GetUserDataByID(sessionContext mongo.SessionContext, userID primitive.ObjectID) (models.UserInfo, error) {
	user := models.UserInfo{ID: userID}
	err := store.mongo.Collection(models.USER_INFO_COLLECTION).FindOne(sessionContext, user).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, types.NewError(types.ErrInvalidParams, "用户不存在")
		}
		return user, types.NewError(types.ErrServerError, err.Error())
	}

	return user, nil
}

/*
GetUserDataByEmail 通过邮箱获取用户信息

参数：
  - sessionContext：数据库会话上下文
  - email：邮箱

返回：
  - models.UserInfo：用户信息
  - error：错误信息
*/
func (store *UserStorage) GetUserDataByEmail(sessionContext mongo.SessionContext, email string) (models.UserInfo, error) {
	user := models.UserInfo{Email: email}
	err := store.mongo.Collection(models.USER_INFO_COLLECTION).FindOne(sessionContext, user).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, types.NewError(types.ErrInvalidParams, "用户不存在")
		}
		return user, types.NewError(types.ErrServerError, err.Error())
	}

	return user, nil
}

/*
GetUserDataByUsername 通过用户名获取用户信息

参数：
  - sessionContext：数据库会话上下文
  - username：用户名

返回：
  - models.UserInfo：用户信息
  - error：错误信息
*/
func (store *UserStorage) GetUserDataByUsername(sessionContext mongo.SessionContext, username string) (models.UserInfo, error) {
	user := models.UserInfo{UserName: username}
	err := store.mongo.Collection(models.USER_INFO_COLLECTION).FindOne(sessionContext, user).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, types.NewError(types.ErrInvalidParams, "用户不存在")
		}
		return user, types.NewError(types.ErrServerError, err.Error())
	}

	return user, nil
}

/*
UpdateUserProfile 更新用户信息

参数：
  - sessionContext：数据库会话上下文
  - userInfo：用户信息

返回：
  - error：错误信息
*/
func (store *UserStorage) UpdateUserProfile(sessionContext mongo.SessionContext, userInfo models.UserInfo) error {
	_, err := store.mongo.Collection(models.USER_INFO_COLLECTION).UpdateOne(sessionContext, bson.M{"_id": userInfo.ID}, bson.M{"$set": userInfo})
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}

/*
UploadAvatarFile 上传用户头像文件

参数：
  - ctx 上下文
  - fileName 文件名
  - avatarData 头像数据

返回：
  - error：错误信息
*/
func (store *UserStorage) UploadAvatarFile(ctx context.Context, fileName string, avatarData io.Reader) (minio.UploadInfo, error) {
	info, err := store.minio.PutObject(
		ctx,
		models.USER_AVATAR_BUCKET,
		fileName,
		avatarData,
		-1,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return info, types.NewError(types.ErrServerError, err.Error())
	}

	return info, nil
}

/*
DeleteAvatarFile 删除用户头像文件

参数：
  - ctx 上下文
  - fileName 文件名

返回：
  - error：错误信息
*/
func (store *UserStorage) DeleteAvatarFile(ctx context.Context, fileName string) error {
	return store.minio.RemoveObject(ctx, models.USER_AVATAR_BUCKET, fileName, minio.RemoveObjectOptions{})
}

/*
UpdateUserPassword 更新用户密码

参数：
  - sessionContext：数据库会话上下文
  - userID：用户ID
  - hashedPassword：哈希密码

返回：
  - error：错误信息
*/
func (store *UserStorage) UpdateUserPassword(sessionContext mongo.SessionContext, userID primitive.ObjectID, hashedPassword string) error {
	authInfo := models.UserAuthInfo{PasswordHash: hashedPassword}
	_, err := store.mongo.Collection(models.USER_AUTH_INFO_COLLECTION).UpdateOne(
		sessionContext,
		bson.M{"_id": userID},
		authInfo,
	)
	if err != nil {
		return types.NewError(types.ErrServerError, err.Error())
	}

	return nil
}
