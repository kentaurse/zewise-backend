/*
Package models - ZeWise 数据库模型
该文件用于声明用户相关模型
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserInfo 用户信息模型
type UserInfo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`       // 主键
	UserName  string             `bson:"username,omitempty"`  // 用户名
	NickName  string             `bson:"nickname,omitempty"`  // 昵称
	Email     string             `bson:"email,omitempty"`     // 邮箱
	Avatar    string             `bson:"avatar,omitempty"`    // 头像
	Sign      string             `bson:"sign,omitempty"`      // 签名
	Birth     time.Time          `bson:"birth,omitempty"`     // 生日
	Gender    string             `bson:"gender,omitempty"`    // 性别
	Authority uint64             `bson:"authority,omitempty"` // 权限等级
	Level     uint64             `bson:"level,omitempty"`     // 等级
}

const USER_INFO_COLLECTION = "user_info"

// UserAuthInfo 用户认证信息模型
type UserAuthInfo struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`      // 主键
	UserName     string             `bson:"username,omitempty"` // 用户名
	Email        string             `bson:"email,omitempty"`    // 邮箱
	Salt         string             `bson:"salt,omitempty"`     // 盐
	PasswordHash string             `bson:"psw_hash,omitempty"` // 密码哈希值
}

const USER_AUTH_INFO_COLLECTION = "user_auth_info"
