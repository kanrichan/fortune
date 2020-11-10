package main

import (
	ft "fortune/fortune"
	"github.com/wdvxdr1123/ZeroBot"
)

func main() {
	zero.Run(zero.Option{
		Host:          ft.Conf.Host,
		Port:          ft.Conf.Port,
		AccessToken:   ft.Conf.AccessToken,
		NickName:      []string{"ft"},
		CommandPrefix: "ft",
		SuperUsers:    []string{ft.Conf.Master},
	})
	select {}
}

func init() {
	ft.Init()
}

/*
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
*/
