package main

import (
	"encoding/json"

	ft "fortune/fortune"
	"github.com/Yiwen-Chan/xq-go/core"
)

type info struct {
	Name   string `json:"name"`
	Pver   string `json:"pver"`
	Sver   int    `json:"sver"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
}

func appinfo() *info {
	return &info{
		Name:   "fortune-运势",
		Pver:   "1.0.5",
		Sver:   3,
		Author: "kanri",
		Desc:   "项目地址 https://github.com/Yiwen-Chan/fortune",
	}
}

// 连接 core
func main() { core.Main() }
func init() {
	data, _ := json.Marshal(appinfo())
	core.InfoJson = string(data)
	core.OnEnable = onEnable
	core.OnGroupMsg = onGroupMsg
	core.OnPrivateMsg = onPrivateMsg
}

// 插件初始化
func onEnable() {
	ft.Init()
}

// 处理消息
func onGroupMsg(botID int64, messageID int64, groupID int64, userID int64, message string) {
	if message == "xqgo -v" || message == "xqgo -version" {
		core.SendGroupMsg(botID, groupID, "[XQ-GO] Version 1.0.1 By Kanri", 0)
	} else if message == "ft -v" || message == "ft -version" {
		core.SendGroupMsg(botID, groupID, "[fortune-运势] Version 1.0.5 By Kanri", 0)
	} else if message == "ft -r" || message == "ft -reload" {
		ft.Init()
		core.SendGroupMsg(botID, groupID, "[fortune-运势] Fortune Reloaded!", 0)
	}
	ft.App(botID, messageID, groupID, userID, message)
}

func onPrivateMsg(botID int64, messageID int64, userID int64, message string) {
	if message == "xqgo -v" || message == "xqgo -version" {
		core.SendGroupMsg(botID, userID, "[XQ-GO] Version 1.0.1 By Kanri", 0)
	} else if message == "ft -v" || message == "ft -version" {
		core.SendGroupMsg(botID, userID, "[fortune-运势] Version 1.0.5 By Kanri", 0)
	} else if message == "ft -r" || message == "ft -reload" {
		ft.Init()
		core.SendGroupMsg(botID, userID, "[fortune-运势] Fortune Reloaded!", 0)
	}
	//ft.App(botID, messageID, 0, userID, message)
}
