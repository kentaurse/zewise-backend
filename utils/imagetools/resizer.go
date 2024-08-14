/*
Package image tools - ZeWise 图片工具
该文件用于定义图片格式转换器
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package imagetools

import (
	"image"

	"github.com/KononK/resize"
)

// ImageSize 图片尺寸
type ImageSize struct {
	Width  int // 宽度
	Height int // 高度
}

// ResizeContain 图片缩放处理器
type ResizeProcessor interface {
	Resize(imageObject image.Image, imageConfig image.Config, width int, height int) (image.Image, ImageSize)
}

// ScallingDownProcessor 等比例缩小模式处理器
// 如果图片尺寸大于目标尺寸 则进行等比缩小 否则不进行处理
type ScallingDownProcessor struct{}

func (processor *ScallingDownProcessor) Resize(imageObject image.Image, imageConfig image.Config, width int, height int) (image.Image, ImageSize) {
	// 如果图片宽度大于高度且宽度大于目标宽度
	if imageConfig.Width > imageConfig.Height && imageConfig.Width > width {
		targetHeight := imageConfig.Height * width / imageConfig.Width
		return resize.Resize(uint(width), uint(targetHeight), imageObject, resize.Lanczos3), ImageSize{Width: width, Height: targetHeight}
	}
	// 如果图片高度大于宽度且高度大于目标高度
	if imageConfig.Height > imageConfig.Width && imageConfig.Height > height {
		targetWidth := imageConfig.Width * height / imageConfig.Height
		return resize.Resize(uint(targetWidth), uint(height), imageObject, resize.Lanczos3), ImageSize{Width: targetWidth, Height: height}
	}
	// 如果图片宽度小于高度且宽度小于目标宽度
	return imageObject, ImageSize{Width: imageConfig.Width, Height: imageConfig.Height}
}

// FillProcessor 填充模式处理器
// 直接将图片缩放至目标大小
type FillProcessor struct{}

func (processor *FillProcessor) Resize(imageObject image.Image, imageConfig image.Config, width int, height int) (image.Image, ImageSize) {
	return resize.Resize(uint(width), uint(height), imageObject, resize.Lanczos3), ImageSize{Width: width, Height: height}
}

// ResizeProcessHandler 图片缩放处理器
type ResizeProcessHandler struct {
	Processor    ResizeProcessor
	TargetWidth  int
	TargetHeight int
}

func (handler *ResizeProcessHandler) Process(imageObject image.Image, imageConfig image.Config) (image.Image, image.Config, error) {
	imageObject, size := handler.Processor.Resize(imageObject, imageConfig, handler.TargetWidth, handler.TargetHeight)
	imageConfig.Width = size.Width
	imageConfig.Height = size.Height
	return imageObject, imageConfig, nil
}

// NewResizeProcessHandler 新建图片缩放处理器
//
// 参数：
//   - width：目标宽度
//   - height：目标高度
//   - processor：缩放处理器
//
// 返回：
//   - ImageProcessHandler：图片处理器
func NewResizeProcessHandler(width int, height int, processor ResizeProcessor) ImageProcessHandler {
	return &ResizeProcessHandler{
		Processor:    processor,
		TargetWidth:  width,
		TargetHeight: height,
	}
}
