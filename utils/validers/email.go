/*
Package validers - ZeWise 工具函数包
该文件用于定义邮箱验证器函数
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package validers

import (
	"regexp"

	"zewise.space/backend/consts"
)

/*
IsValidEmail 验证邮箱是否合法

参数：
  - email：邮箱

返回：
  - bool：是否合法
*/
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(consts.EMAIL_REGEX)
	return re.MatchString(email)
}
