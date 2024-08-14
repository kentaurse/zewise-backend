/*
Package validers - ZeWise 工具函数包
该文件用于定义密码验证器函数
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package validers

import (
	"regexp"

	"zewise.space/backend/consts"
)

/*
IsValidPassword 验证密码是否合法

参数：
  - password：密码

返回：
  - bool：是否合法
*/
func IsValidPassword(password string) bool {
	length := len(password)
	if length < consts.PASSWORD_MIN_LENGTH || length > consts.PASSWORD_MAX_LENGTH {
		return false
	}

	re := regexp.MustCompile(consts.PASSWORD_REGEX)
	return re.MatchString(password)
}
