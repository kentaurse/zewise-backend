/*
Package models - ZeWise 数据库模型
该文件用于声明回复相关模型
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ReplyInfo 评论信息模型
type ReplyInfo struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`             // 主键
	ParentReplyID primitive.ObjectID `bson:"parent_reply_id,omitempty"` // 父回复ID
	UID           primitive.ObjectID `bson:"uid,omitempty"`             // 用户ID
	Content       string             `bson:"content,omitempty"`         // 内容
	IsPublic      bool               `bson:"is_public,omitempty"`       // 是否公开
}
