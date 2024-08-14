/*
Package image tools - ZeWise 图片工具
该文件用于定义图片解码器
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package imagetools

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/webp"
)

// ErrFormatNotSupported 格式不支持错误
var ErrFormatNotSupported = errors.New("不支持该图片格式")

// ImageDecoder 图片解码器
type ImageDecoder interface {
	Decode(imageFile *ImageFile) (image.Image, image.Config, error) // 解码图片
}

// ImageDecodeHandler 图片解码处理器
type ImageDecodeHandler interface {
	Match(fileType string) bool                                     // 匹配文件类型
	Decode(imageFile *ImageFile) (image.Image, image.Config, error) // 解码图片
}

// PNGImageDecoder PNG图片解码器
type PNGImageDecoder struct{}

func (decoder *PNGImageDecoder) Match(fileType string) bool {
	return fileType == "image/png"
}

func (decoder *PNGImageDecoder) Decode(imageFile *ImageFile) (image.Image, image.Config, error) {
	// 解码图片配置
	imageConfig, err := png.DecodeConfig(*imageFile)
	// 解码失败
	if err != nil {
		return nil, image.Config{}, err
	}
	// 重置文件指针
	_, err = (*imageFile).Seek(0, 0)
	if err != nil {
		return nil, image.Config{}, err
	}

	// 解码图片
	imageObject, err := png.Decode(*imageFile)
	// 解码失败
	if err != nil {
		return nil, image.Config{}, err
	}
	// 重置文件指针
	_, err = (*imageFile).Seek(0, 0)
	if err != nil {
		return nil, image.Config{}, err
	}

	// 返回结果
	return imageObject, imageConfig, nil
}

// JPEGImageDecoder JPEG图片解码器
type JPEGImageDecoder struct{}

func (decoder *JPEGImageDecoder) Match(fileType string) bool {
	return fileType == "image/jpeg"
}

func (decoder *JPEGImageDecoder) Decode(imageFile *ImageFile) (image.Image, image.Config, error) {
	// 解码图片配置
	imageConfig, err := jpeg.DecodeConfig(*imageFile)
	// 解码失败
	if err != nil {
		return nil, image.Config{}, err
	}
	// 重置文件指针
	_, err = (*imageFile).Seek(0, 0)
	if err != nil {
		return nil, image.Config{}, err
	}

	// 解码图片
	imageObject, err := jpeg.Decode(*imageFile)
	// 解码失败
	if err != nil {
		return nil, image.Config{}, err
	}
	// 重置文件指针
	_, err = (*imageFile).Seek(0, 0)
	if err != nil {
		return nil, image.Config{}, err
	}

	// 返回结果
	return imageObject, imageConfig, nil
}

// WebpImageDecoder Webp图片解码器
type WebpImageDecoder struct{}

func (decoder *WebpImageDecoder) Match(fileType string) bool {
	return fileType == "image/webp"
}

func (decoder *WebpImageDecoder) Decode(imageFile *ImageFile) (image.Image, image.Config, error) {
	// 解码图片配置
	imageConfig, err := webp.DecodeConfig(*imageFile)
	// 解码失败
	if err != nil {
		return nil, image.Config{}, err
	}
	// 重置文件指针
	_, err = (*imageFile).Seek(0, 0)
	if err != nil {
		return nil, image.Config{}, err
	}

	// 解码图片
	imageObject, err := webp.Decode(*imageFile)
	// 解码失败
	if err != nil {
		return nil, image.Config{}, err
	}
	// 重置文件指针
	_, err = (*imageFile).Seek(0, 0)
	if err != nil {
		return nil, image.Config{}, err
	}

	// 返回结果
	return imageObject, imageConfig, nil
}

// ImageDecoderChain 图片解码器链
type ImageDecoderChain struct {
	ContentType string
	Decoders    []ImageDecodeHandler
}

/*
Decode 解码图片

参数：
  - imageFile：图片文件

返回：
  - image.Image：图片对象
  - image.Config：图片配置
  - error：错误信息
*/
func (chain *ImageDecoderChain) Decode(imageFile *ImageFile) (image.Image, image.Config, error) {
	if chain.ContentType == "" {
		return nil, image.Config{}, errors.New("需要设置内容类型")
	}

	var (
		imageObject image.Image
		imageConfig image.Config
		err         error
	)
	// 解码图片
	for _, decoder := range chain.Decoders {
		// 匹配解码器
		if decoder.Match(chain.ContentType) {
			imageObject, imageConfig, err = decoder.Decode(imageFile)
			if err != nil {
				return nil, image.Config{}, err
			}
			break
		}
	}
	// 未找到解码器
	if imageObject == nil {
		return nil, image.Config{}, ErrFormatNotSupported
	}

	// 返回结果
	return imageObject, imageConfig, nil
}

/*
NewImageDecoderChain 新建图片解码器链

参数：
  - contentType：内容类型
  - decoders：图片解码器

返回：
  - *ImageDecoderChain：图片解码器链
*/
func NewImageDecoderChain(contentType string, decoders ...ImageDecodeHandler) *ImageDecoderChain {
	return &ImageDecoderChain{
		ContentType: contentType,
		Decoders:    decoders,
	}
}

/*
NewDefaultImageDecoderChain 新建默认图片解码器链

返回：
  - *ImageDecoderChain：图片解码器链

默认解码器：
  - JPEGImageDecoder：JPEG图片解码器
  - PNGImageDecoder：PNG图片解码器
  - WebpImageDecoder：Webp图片解码器

注意：使用默认解码器时，需使用 SetContentType 方法设置内容类型
*/
func NewDefaultImageDecoderChain() *ImageDecoderChain {
	return NewImageDecoderChain("", &JPEGImageDecoder{}, &PNGImageDecoder{}, &WebpImageDecoder{})
}

/*
SetContentType 设置内容类型

参数：
  - contentType：内容类型 参考：https://www.iana.org/assignments/media-types/media-types.xhtml#image
*/
func (chain *ImageDecoderChain) SetContentType(contentType string) {
	chain.ContentType = contentType
}
