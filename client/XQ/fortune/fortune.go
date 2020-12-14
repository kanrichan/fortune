package fortune

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"yaya/core"
)

// 应用函数
func App(botID int64, type_ int64, groupID int64, userID int64, message string) {
	// 获取本群运势配置
	setting := getSetting(Conf, core.Int2Str(groupID))

	// 如果匹配到关键词则开始处理信息
	if message == setting.Trigger && setting.Trigger != "关" {

		// 获得POST提交表单数据
		fromDataStruct := &FromDataStruct{
			client:    ClientName,
			version:   ClientVer,
			bot:       core.Int2Str(botID),
			types:     getTypes(setting.Types),
			fromGroup: core.Int2Str(groupID),
			fromQQ:    core.Int2Str(userID),
			ask:       message,
			limit:     setting.Limit,
		}

		// 获得POST提交的Header
		headerParm := getHeader(ClientKey)

		// 获得服务器数据
		fortuneJson, code := fortune(ApiFortune, fromDataStruct, headerParm)

		// 开始处理信息提交方式
		if code == 200 && fortuneJson.Msg == "fortuned" {
			if setting.Warm != "" {
				sendMessage(botID, type_, groupID, userID, setting.Warm)
				return
			}
		} else {
			if setting.Reply != "" {
				sendMessage(botID, type_, groupID, userID, setting.Reply)
			}
		}

		switch {
		case code != 200:
			// 服务器无响应
			sendMessage(botID, type_, groupID, userID, "[fortune-运势] 服务器失联中......")
			break
		case fortuneJson.Code != 200:
			// 服务器状态异常
			sendMessage(botID, type_, groupID, userID, fortuneJson.Msg)
			break
		case fortuneJson.Msg != "success":
			sendMessage(botID, type_, groupID, userID, fortuneJson.Msg)
			break
		case fortuneJson.Warn != "":
			// 有警告信息
			sendMessage(botID, type_, groupID, userID, fortuneJson.Warn)
			pic(ApiPic, fromDataStruct, headerParm)
			sendPicture(botID, type_, groupID, userID, PicPath)
			break
		case fortuneJson.Info != "":
			// 有提示消息
			if notSend() {
				sendMessage(botID, type_, groupID, userID, fortuneJson.Info)
			}
			pic(ApiPic, fromDataStruct, headerParm)
			sendPicture(botID, type_, groupID, userID, PicPath)
		default:
			// 成功
			pic(ApiPic, fromDataStruct, headerParm)
			sendPicture(botID, type_, groupID, userID, PicPath)
		}
	}
}

func getTypes(types string) string {
	if strings.Contains(types, "|") {
		list := strings.Split(types, "|")
		length := len(list)
		ran := rand.Intn(length)
		types = list[ran]
	}
	return types
}

var Note string = ""

func notSend() bool {
	todayTime := time.Now().Day()
	today := fmt.Sprintf("%v", todayTime)
	if Note == "" {
		Note = today
		return true
	} else if Note != today {
		Note = today
		return true
	} else {
		return false
	}
}
