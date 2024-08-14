/*
Package validers - ZeWise 工具函数包
该文件用于定义用户名验证器函数
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package validers

import (
	"regexp"

	"zewise.space/backend/consts"
)

/*
IsValidUsername 验证用户名是否合法

参数：
  - username：用户名

返回：
  - bool：是否合法
*/
func IsValidUsername(username string) bool {
	length := len(username)
	if length < consts.USERNAME_MIN_LENGTH || length > consts.USERNAME_MAX_LENGTH {
		return false
	}

	re := regexp.MustCompile(consts.USERNAME_REGEX)
	return re.MatchString(username)
}
