/*
Package image tools - ZeWise 图片工具
该文件用于定义图片限制器
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package imagetools

import (
	"errors"
	"image"
)

var ErrImageSizeExceed = errors.New("图片尺寸超出限制")

// SizeLimiter 图片尺寸限制器
type SizeLimiter struct {
	MaxWidth  int // 最大宽度
	MaxHeight int // 最大高度
}

func (limiter *SizeLimiter) Process(imageObject image.Image, imageConfig image.Config) (image.Image, image.Config, error) {
	if imageConfig.Width > limiter.MaxWidth || imageConfig.Height > limiter.MaxHeight {
		return nil, image.Config{}, ErrImageSizeExceed
	}
	return imageObject, imageConfig, nil
}

/*
NewSizeLimiter 新建图片尺寸限制器

参数：
  - maxWidth：最大宽度
  - maxHeight：最大高度

返回：
  - *SizeLimiter：图片尺寸限制器
*/
func NewSizeLimiter(maxWidth int, maxHeight int) *SizeLimiter {
	return &SizeLimiter{
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
	}
}
