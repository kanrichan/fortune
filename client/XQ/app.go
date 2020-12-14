package main

import (
	"encoding/json"

	ft "fortune/fortune"
)

// 插件信息
type AppInfo struct {
	Name   string `json:"name"`   // 插件名字
	Pver   string `json:"pver"`   // 插件版本
	Sver   int    `json:"sver"`   // 框架版本
	Author string `json:"author"` // 作者名字
	Desc   string `json:"desc"`   // 插件说明
}

func newAppInfo() *AppInfo {
	return &AppInfo{
		Name:   "fortune-运势",
		Pver:   "1.0.6",
		Sver:   3,
		Author: "kanri",
		Desc:   "The best of luck! 项目地址 https://github.com/Yiwen-Chan/fortune",
	}
}

// 连接 core

func init() {
	data, _ := json.Marshal(newAppInfo())
	ft.AppInfoJson = string(data)
}

func main() { ft.Main() }
