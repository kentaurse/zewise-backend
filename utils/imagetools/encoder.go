/*
Package image tools - ZeWise 图片工具
该文件用于定义图片编码器
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package imagetools

import (
	"image"

	"github.com/chai2010/webp"
)

// ImageEncoder 图片编码器
type ImageEncoder interface {
	GetFormatFileSuffix() string
	Encode(imageObject image.Image, imageConfig image.Config) ([]byte, error)
}

// WebpImageEncoder Webp图片编码器
type WebpImageEncoder struct {
	Quality float32 // 图片质量
}

func (encoder *WebpImageEncoder) GetFormatFileSuffix() string {
	return "webp"
}

func (encoder *WebpImageEncoder) Encode(imageObject image.Image, imageConfig image.Config) ([]byte, error) {
	// 编码图片
	return webp.EncodeRGBA(imageObject, encoder.Quality)
}

/*
NewWebpImageEncoder 新建Webp图片编码器

参数：
  - quality：图片质量

返回：
  - *WebpImageEncoder：Webp图片编码器
*/
func NewWebpImageEncoder(quality float32) *WebpImageEncoder {
	return &WebpImageEncoder{
		Quality: quality,
	}
}
