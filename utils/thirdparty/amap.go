/*
Package thirdparty - ZeWise 后端服务器第三方服务包
该文件用于请求高德地图 API
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package thirdparty

import (
	"github.com/gofiber/fiber/v2"
	"zewise.space/backend/utils/functools"
)

type AMapResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	InfoCode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Adcode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

/*
RequestIPInfo 请求 IP 地理位置

参数：
  - ip：IP 地址
  - key：高德地图 API 密钥

返回：
  - string：地理位置
  - error：错误信息
*/
func RequestIPInfo(ip string, key string) (AMapResponse, error) {
	// 请求API
	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.SetRequestURI("https://restapi.amap.com/v3/ip")
	req.URI().SetQueryString(functools.JoinStrings("key=", key, "&ip=", ip))

	if err := agent.Parse(); err != nil {
		return AMapResponse{}, err
	}

	// 解析响应
	var response AMapResponse
	code, _, errs := agent.Struct(&response)
	if len(errs) > 0 {
		return AMapResponse{}, errs[0]
	}
	if code != 200 {
		return AMapResponse{}, fiber.ErrBadRequest
	}

	return response, nil
}
