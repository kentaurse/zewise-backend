/*
Package image tools - ZeWise 图片工具
该文件用于定义图片处理器
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package imagetools

import (
	"image"
	"io"
)

// ImageFile 图片文件接口
type ImageFile interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

// ImageProcessHandler 图片处理器
type ImageProcessHandler interface {
	Process(imageFile image.Image, imageConfig image.Config) (image.Image, image.Config, error)
}

/*
ProcessImage 处理图片

参数：
  - imageFile：图片文件
  - decoder：图片解码器
  - encoder：图片编码器
  - processHandlers：图片处理器

返回：
  - []byte：图片数据
  - error：错误
*/
func ProcessImage(imageFile ImageFile, decoder ImageDecoder, encoder ImageEncoder, processHandlers ...ImageProcessHandler) ([]byte, error) {
	// 解码图片
	imageObject, imageConfig, err := decoder.Decode(&imageFile)
	if err != nil {
		return nil, err
	}
	// 处理图片
	for _, handler := range processHandlers {
		imageObject, imageConfig, err = handler.Process(imageObject, imageConfig)
		if err != nil {
			return nil, err
		}
	}
	// 编码图片
	return encoder.Encode(imageObject, imageConfig)
}
