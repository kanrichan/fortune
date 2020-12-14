package fortune

import (
	"yaya/core"
)

func init() {
	core.Create = XQCreate
	core.Event = XQEvent
	core.DestroyPlugin = XQDestroyPlugin
	core.SetUp = XQSetUp
}

func XQCreate(version string) string {
	return AppInfoJson
}

func XQEvent(selfID int64, mseeageType int64, subType int64, groupID int64, userID int64, noticID int64, message string, messageNum int64, messageID int64, rawMessage []byte, time int64, ret int64) int64 {
	switch mseeageType {
	case 12001:
		go ProtectRun(func() { onStart() }, "onStart()")
	// 消息事件
	// 0：临时会话 1：好友会话 4：群临时会话 7：好友验证会话
	case 0, 1, 4, 5, 7:
		go ProtectRun(func() { onPrivateMessage(selfID, mseeageType, groupID, userID, message) }, "onPrivateMessage()")
	// 2：群聊信息
	case 2, 3:
		go ProtectRun(func() { onGroupMessage(selfID, mseeageType, groupID, userID, message) }, "onGroupMessage()")
	default:
		//
	}
	return 0
}

func XQDestroyPlugin() int64 {
	return 0
}

func XQSetUp() int64 {
	return 0
}

func onPrivateMessage(botID int64, type_ int64, groupID int64, userID int64, message string) {
	switch message {
	case "ft -v", "ft -version":
		sendMessage(botID, type_, groupID, userID, "Fortune-运势 Version 1.0.6 BY Kanri")
	case "ft -r", "ft -reload":
		Conf = Load(AppPath + "config.yml")
		if Conf != nil {
			sendMessage(botID, type_, groupID, userID, "Setting is already reload!")
		} else {
			sendMessage(botID, type_, groupID, userID, "Setting ERROR!")
		}
	case "ft -c", "ft -core":
		sendMessage(botID, type_, groupID, userID, "OneBot-YaYa Version 1.1.3 Beta BY Kanri  Using only core")
	default:
	}
}

func onGroupMessage(botID int64, type_ int64, groupID int64, userID int64, message string) {
	App(botID, type_, groupID, userID, message)
	switch message {
	case "ft -v", "ft -version":
		sendMessage(botID, type_, groupID, userID, "Fortune-运势 Version 1.0.6 BY Kanri")
	case "ft -r", "ft -reload":
		Conf = Load(AppPath + "config.yml")
		if Conf != nil {
			sendMessage(botID, type_, groupID, userID, "Setting is already reload!")
		} else {
			sendMessage(botID, type_, groupID, userID, "Setting ERROR!")
		}
	case "ft -c", "ft -core":
		sendMessage(botID, type_, groupID, userID, "OneBot-YaYa Version 1.1.3 Beta BY Kanri  Using only core")
	default:
	}

}
