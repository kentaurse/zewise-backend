/*
Package services - ZeWise 服务层
该文件用于声明服务对象类
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package services

import "zewise.space/backend/stores"

// Service 服务对象
type Service struct {
	storage     *stores.Storage // 存储对象
	UserService *UserService    // 用户服务
	AuthService *AuthService    // 认证服务
}

/*
NewService 新建服务对象

参数：
  - storage：存储对象

返回：
  - *Service：服务对象
*/
func NewService(storage *stores.Storage) *Service {
	return &Service{
		storage:     storage,
		UserService: &UserService{storage},
		AuthService: &AuthService{storage},
	}
}
