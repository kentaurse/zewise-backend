/*
Package models - NekoBlog backend server database models
This file is for post related models.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// PostInfo 博文信息模型
type PostInfo struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`            // 主键
	ParentPostID primitive.ObjectID   `bson:"parent_post_id,omitempty"` // 父博文ID
	UID          primitive.ObjectID   `bson:"uid,omitempty"`            // 用户ID
	IpAddrress   string               `bson:"ip_address,omitempty"`     // IP地址
	Title        string               `bson:"title,omitempty"`          // 标题
	Content      string               `bson:"content,omitempty"`        // 内容
	MediaIDs     []primitive.ObjectID `bson:"media_ids,omitempty"`      // 媒体ID
	IsPublic     bool                 `bson:"is_public,omitempty"`      // 是否公开
}
