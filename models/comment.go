/*
Package models - NekoBlog backend server database models
This file is for comment related models.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// CommentInfo 评论信息模型
type CommentInfo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`       // 主键
	UID      primitive.ObjectID `bson:"uid,omitempty"`       // 用户ID
	Username string             `bson:"username,omitempty"`  // 用户名
	Content  string             `bson:"content,omitempty"`   // 内容
	IsPublic bool               `bson:"is_public,omitempty"` // 是否公开
}
