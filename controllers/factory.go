/*
Package controllers - ZeWise 控制器
该文件用于声明控制器工厂
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package controllers

import "zewise.space/backend/services"

// Factory 控制器工厂
type Factory struct {
	service *services.Service // 服务对象
}

/*
NewFactory 创建控制器工厂

参数：
  - service：服务对象

返回：
  - *Factory：控制器工厂对象
*/
func NewFactory(service *services.Service) *Factory {
	return &Factory{service}
}
