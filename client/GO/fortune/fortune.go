package fortune

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/wdvxdr1123/ZeroBot"
)

// 全局变量初始化
var ApiHost = "api.kanri.top"
var ApiPort = "10086"
var ApiFortune = "http://" + ApiHost + ":" + ApiPort + "/fortune"
var ApiPic = "http://" + ApiHost + ":" + ApiPort + "/fortune.jpg"

var ClientKey = "233666"
var ClientName = "xq"
var ClientVer = "5"

var AppPath = strings.Replace(PathExecute(), ":\\", ":\\\\", -1) + "data\\"
var ConfPath = AppPath + "config.yml"
var PicPath = AppPath + "output.jpg"
var ResultPath = AppPath + "output.txt"

var Conf = &YamlConfig{}

func init() {
	a := testPlugin{}
	zero.RegisterPlugin(a) // 注册插件
}

type testPlugin struct{}

func (testPlugin) GetPluginInfo() zero.PluginInfo { // 返回插件信息
	return zero.PluginInfo{
		Author:     "kanri",
		PluginName: "fortune-运势",
		Version:    "1.0.5",
		Details:    "项目地址 https://github.com/Yiwen-Chan/fortune",
	}
}

func (testPlugin) Start() { // 插件主体
	zero.OnMessage().
		Handle(
			func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
				App(0, 0, event.GroupID, event.UserID, event.RawMessage)
				return zero.SuccessResponse
			},
		)
}

// 增加发信息函数
func AllSendMsg(botID int64, groupID int64, userID int64, text string) {
	if groupID != 0 {
		zero.SendGroupMessage(groupID, text)
	} else if userID != 0 {
		zero.SendPrivateMessage(userID, text)
	}
}
func OutPutLog(str string) {
	fmt.Println(str)
}

func Init() {
	OutPutLog("本插件基于 ZeroBot 开发")
	err := CreatePath(AppPath)
	if err != nil {
		OutPutLog("创建应用文件夹时出现错误")
	}
	if PathExists(ConfPath) {
		Conf = Load(ConfPath)
	} else {
		err = DefaultConfig().Save(ConfPath)
		if err != nil {
			OutPutLog("创建默认配置文件时出现错误")
			return
		}
		Conf = Load(ConfPath)
		OutPutLog("[fortune-运势] 检测到初次运行本插件，已生成默认配置文件")
		OutPutLog("[fortune-运势] 特别感谢 fz6m https://github.com/fz6m/nonebot-plugin/tree/master/CQVortune")
		OutPutLog("[fortune-运势] 特别感谢 Lostdegree https://github.com/Lostdegree/Portune")
		OutPutLog("[fortune-运势] 有需要请按GitHub项目上描述的方法修改配置文件")
	}
	OutPutLog("[fortune-运势] 项目地址 https://github.com/Yiwen-Chan/fortune")
	OutPutLog("[fortune-运势] 配置文件 XQ/data/app/fortune/config.yml")
	OutPutLog("[fortune-运势] 想自定义运势背景并共享可加QQ群 1048452984 ")
}

// 应用函数
func App(botID int64, messageID int64, groupID int64, userID int64, message string) {
	// 获取本群运势配置
	setting := getSetting(Conf, Int2Str(groupID))

	// 如果匹配到关键词则开始处理信息
	if message == setting.Trigger && setting.Trigger != "关" {

		// 获得POST提交表单数据
		fromDataStruct := &FromDataStruct{
			client:    ClientName,
			version:   ClientVer,
			bot:       Int2Str(botID),
			types:     getTypes(setting.Types),
			fromGroup: Int2Str(groupID),
			fromQQ:    Int2Str(userID),
			ask:       message,
			limit:     setting.Limit,
		}

		// 获得POST提交的Header
		headerParm := getHeader(ClientKey)

		// 获得服务器数据
		fortuneJson, code := fortune(ApiFortune, fromDataStruct, headerParm)

		// 开始处理信息提交方式
		text := ""

		if code == 200 && fortuneJson.Msg != "fortuned" {
			if setting.Reply != "" {
				AllSendMsg(botID, groupID, userID, setting.Reply)
			}
		}

		if code != 200 {
			//服务器无响应
			text = "[fortune-运势] 服务器失联中......"
			AllSendMsg(botID, groupID, userID, text)
		} else if fortuneJson.Code != 200 {
			//服务器状态异常
			text = fortuneJson.Msg
			AllSendMsg(botID, groupID, userID, text)
		} else if fortuneJson.Msg == "fortuned" {
			text = setting.Warm
			AllSendMsg(botID, groupID, userID, text)
		} else if fortuneJson.Msg != "success" {
			text = fortuneJson.Msg
			AllSendMsg(botID, groupID, userID, text)
		} else if fortuneJson.Warn != "" {
			text = fortuneJson.Warn
			pic(ApiPic, fromDataStruct, headerParm)
			text += "[CQ:image,file=file:///" + PicPath + "]"
			AllSendMsg(botID, groupID, userID, text)
		} else if fortuneJson.Info != "" {
			if notSend() {
				text = fortuneJson.Info
			}
			pic(ApiPic, fromDataStruct, headerParm)
			text += "[CQ:image,file=file:///" + PicPath + "]"
			AllSendMsg(botID, groupID, userID, text)
		} else {
			pic(ApiPic, fromDataStruct, headerParm)
			text += "[CQ:image,file=file:///" + PicPath + "]"
			AllSendMsg(botID, groupID, userID, text)
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
