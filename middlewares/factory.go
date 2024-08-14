/*
Package middlewares - NekoBlog backend server middlewares.
This file is for middlewares factory.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package middlewares

import "zewise.space/backend/stores"

// Factory 中间件工厂
type Factory struct {
	storage *stores.Storage
}

/*
NewFactory 新建中间件工厂

参数：
  - storage：存储对象

返回：
  - *Factory：中间件工厂对象
*/
func NewFactory(storage *stores.Storage) *Factory {
	return &Factory{storage}
}
