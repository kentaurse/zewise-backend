/*
Package serializers - ZeWise 序列化器包
该文件用于序列化用户信息
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package serializers

import (
	"zewise.space/backend/consts"
	"zewise.space/backend/models"
	"zewise.space/backend/utils/functools"
)

// UserProfileResponse 用户信息响应
type UserProfileResponse struct {
	ID       string `json:"id,omitempty"`       // 用户ID
	Username string `json:"username,omitempty"` // 用户名
	Nickname string `json:"nickname,omitempty"` // 昵称
	Email    string `json:"email,omitempty"`    // 邮箱
	Avatar   string `json:"avatar,omitempty"`   // 头像
	Sign     string `json:"sign,omitempty"`     // 签名
	Birth    int64  `json:"birth,omitempty"`    // 生日
	Gender   string `json:"gender,omitempty"`   // 性别
	Level    uint64 `json:"level,omitempty"`    // 等级
}

/*
NewUserProfileResponse 创建用户信息响应
*/
func NewUserProfileResponse(data models.UserInfo) UserProfileResponse {
	return UserProfileResponse{
		ID:       data.ID.Hex(),
		Username: data.UserName,
		Nickname: data.NickName,
		Email:    data.Email,
		Avatar:   functools.JoinStrings(consts.AVATAR_URL_PREFIX, data.Avatar, ".webp"),
		Sign:     data.Sign,
		Birth:    data.Birth.Unix(),
		Gender:   data.Gender,
		Level:    data.Level,
	}
}
