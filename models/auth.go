/*
Package models - ZeWise 数据模型
该文件用于声明认证相关模型
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const REDIS_AVAILABLE_USER_TOKEN_LIST = "AUTH:TOKENS"

// UserLoginLog 用户登录日志
type UserLoginLog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"` // 主键
	UID         primitive.ObjectID `bson:"uid,omitempty"` // 用户ID
	IP          string             `bson:"ip"`            // IP 地址
	Location    string             `bson:"location"`      // 地理位置
	Device      string             `bson:"device"`        // 设备
	Time        time.Time          `bson:"time"`          // 登录时间
	Application string             `bson:"application"`   // 登录应用
	IfSucceed   bool               `bson:"if_succeed"`    // 是否成功
	IfChecked   bool               `bson:"if_checked"`    // 是否验证
}

const USER_LOGIN_LOGS_COLLECTION = "user_login_logs"
