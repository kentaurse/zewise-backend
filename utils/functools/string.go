/*
Package functools - ZeWise 工具函数包
该文件用于定义字符串相关工具函数
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package functools

import "strings"

/*
JoinStrings 连接字符串

参数：
  	- strs：字符串列表

返回：
	- string：连接后的字符串
*/
func JoinStrings(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
