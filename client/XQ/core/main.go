package core

import "C"
import (
	log "github.com/sirupsen/logrus"

	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type FromDataParm struct {
	Client    string
	Version   string
	Bot       string
	Types     string
	FromGroup string
	FromQQ    string
	Ask       string
	Limit     string
}

type HeaderParm struct {
	Authkey string
	Au_time string
}

type FortuneJson struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Info string `json:"info"`
	Warn string `json:"warn"`
}

type JsonConfig struct {
	Setting []*SettingConfig `json:"运势设置"`
}

type SettingConfig struct {
	Group   string `json:"群号"`
	Trigger string `json:"触发"`
	Reply   string `json:"回复"`
	Types   string `json:"类型"`
	Limit   string `json:"限制"`
}

func Main() {}

func Create() {
	var conf *JsonConfig
	pathConfig := "Config/fortune-运势/config.json"
	if PathExists(pathConfig) {
		conf = Load(pathConfig)
	}
	if conf == nil {
		err := DefaultConfig().Save(pathConfig)
		if err != nil {
			OutPutLog("创建默认配置文件时出现错误")
			return
		}
		OutPutLog("默认配置文件已生成, 请编辑 config.json 后重启程序.")
		time.Sleep(time.Second * 5)
		return
	}
	OutPutLog("[fortune-运势] 项目地址 https://github.com/Yiwen-Chan/fortune")
	OutPutLog("[fortune-运势] 配置文件 XQ/Config/fortune-运势/config.json")
	OutPutLog("[fortune-运势] 特别感谢 fz6m https://github.com/fz6m/nonebot-plugin/tree/master/CQVortune")
	OutPutLog("[fortune-运势] 特别感谢 Lostdegree https://github.com/Lostdegree/Portune")
	OutPutLog("[fortune-运势] 想自定义运势背景并共享可加QQ群 1048452984 ")
}

func Event(botQQ *C.char, msgType C.int, subType C.int, sourceId *C.char, activeQQ *C.char, passiveQQ *C.char, msg *C.char, msgNum *C.char, msgId *C.char, rawMsg *C.char, timeStamp *C.char, retText *C.char) {
	//OutPutLog(strconv.Itoa(int(msgType)))
	var conf *JsonConfig
	pathConfig := "Config/fortune-运势/config.json"
	if PathExists(pathConfig) {
		conf = Load(pathConfig)
	}
	botID := GoString(botQQ)
	groupID := GoString(sourceId)
	userID := GoString(activeQQ)
	text := GoString(msg)

	setting := returnSetting(conf, groupID)

	trigger := setting.Trigger
	reply := setting.Reply
	types := setting.Types
	limit := setting.Limit

	Types := ""

	if strings.Contains(types, "|") {
		list := strings.Split(types, "|")
		length := len(list)
		ran := rand.Intn(length)
		Types = list[ran]
		types = Types
	}

	if limit == "全局" {
		limit = "on"
	} else if limit == "池子" {
		limit = "none"
	} else if limit == "关" {
		limit = "off"
	} else {
		limit = "none"
	}

	if text == trigger && trigger != "关" {
		if reply != "" {
			SendPrivateMsg(botID, 2, groupID, " ", reply, 0)
		}

		client := "xq"
		version := "4"
		ask := text
		limit := limit

		fromDataParm := &FromDataParm{
			client,
			version,
			botID,
			types,
			groupID,
			userID,
			ask,
			limit,
		}

		authkey, au_time := key("test")
		headerParm := &HeaderParm{
			authkey,
			au_time,
		}

		apiFortune := "http://127.0.0.1:8000/fortune"
		apiPic := "http://127.0.0.1:8000/fortune.jpg"
		fortuneJson, code := fortune(apiFortune, fromDataParm, headerParm)

		Message := ""
		path := PathExecute() + "Config/fortune-运势/output.jpg"

		if code != 200 {
			//服务器无响应
			Message = "[fortune-运势] 服务器失联中......"
			SendPrivateMsg(botID, 2, groupID, " ", Message, 0)
		} else if fortuneJson.Code != 200 {
			//服务器状态异常
			Message = fortuneJson.Msg
			SendPrivateMsg(botID, 2, groupID, " ", Message, 0)
		} else if fortuneJson.Msg != "success" {
			Message = fortuneJson.Msg
			SendPrivateMsg(botID, 2, groupID, " ", Message, 0)
		} else if fortuneJson.Warn != "" {
			Message = fortuneJson.Warn
			pic(apiPic, fromDataParm, headerParm)
			Message += "[pic=" + path + "]"
			SendPrivateMsg(botID, 2, groupID, " ", Message, 0)
		} else if fortuneJson.Info != "" {
			if notSend() {
				Message = fortuneJson.Info
			}
			pic(apiPic, fromDataParm, headerParm)
			Message += "[pic=" + path + "]"
			SendPrivateMsg(botID, 2, groupID, " ", Message, 0)
		} else {
			pic(apiPic, fromDataParm, headerParm)
			Message += "[pic=" + path + "]"
			SendPrivateMsg(botID, 2, groupID, " ", Message, 0)
		}
	}
}

func returnSetting(conf *JsonConfig, groupID string) *SettingConfig {
	setting := conf.Setting[0]
	for index, _ := range conf.Setting {
		settingCache := conf.Setting[index]
		if groupID == settingCache.Group {
			setting := settingCache
			return setting
		}
	}
	return setting
}

func key(au_key string) (string, string) {
	au_time := fmt.Sprintf("%v", time.Now().Unix())

	au_key_time := fmt.Sprintf("%v|%v", au_key, au_time)

	m := md5.New()
	m.Write([]byte(au_key_time))
	authkey := hex.EncodeToString(m.Sum(nil))
	return authkey, au_time
}

func fortune(api string, fromDataParm *FromDataParm, headerParm *HeaderParm) (FortuneJson, int) {
	data := url.Values{}
	data.Set("client", fromDataParm.Client)
	data.Set("version", fromDataParm.Version)
	data.Set("bot", fromDataParm.Bot)
	data.Set("types", fromDataParm.Types)
	data.Set("fromGroup", fromDataParm.FromGroup)
	data.Set("fromQQ", fromDataParm.FromQQ)
	data.Set("ask", fromDataParm.Ask)
	data.Set("limit", fromDataParm.Limit)

	fromdata := strings.NewReader(data.Encode())

	transport := http.Transport{
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: &transport,
	}

	req, err := http.NewRequest("POST", api, fromdata)
	if err != nil {
		log.Error(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authkey", headerParm.Authkey)
	req.Header.Set("autime", headerParm.Au_time)

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}

	defer resp.Body.Close()

	code := resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	fortuneJson := FortuneJson{}
	err = json.Unmarshal(body, &fortuneJson)
	if err != nil {
		log.Error(err)
	}

	return fortuneJson, code
}

func pic(api string, fromDataParm *FromDataParm, headerParm *HeaderParm) {
	data := url.Values{}
	data.Set("client", fromDataParm.Client)
	data.Set("version", fromDataParm.Version)
	data.Set("bot", fromDataParm.Bot)
	data.Set("types", fromDataParm.Types)
	data.Set("fromGroup", fromDataParm.FromGroup)
	data.Set("fromQQ", fromDataParm.FromQQ)
	data.Set("ask", fromDataParm.Ask)
	data.Set("limit", fromDataParm.Limit)

	fromdata := strings.NewReader(data.Encode())

	transport := http.Transport{
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: &transport,
	}

	req, err := http.NewRequest("POST", api, fromdata)
	if err != nil {
		log.Error(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authkey", headerParm.Authkey)
	req.Header.Set("autime", headerParm.Au_time)

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	pathPic := "Config/fortune-运势/output.jpg"
	f, err := os.OpenFile(pathPic, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write([]byte(string(body)))
	}
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

func PathExecute() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dir)

	return dir + "/"
}

func DefaultConfig() *JsonConfig {
	return &JsonConfig{
		Setting: []*SettingConfig{
			{
				Group:   "默认",
				Trigger: "运势",
				Reply:   "少女祈祷中......",
				Types:   "李清歌|碧蓝幻想|公主连结",
				Limit:   "全局",
			},
			{
				Group:   "单个群设置填群号",
				Trigger: "这里填触发关键词",
				Reply:   "这里填收到关键词的回复",
				Types:   "这里是池子类型，多个池子用|分开",
				Limit:   "每天限制一张，可填 全局 池子 关",
			},
			{
				Group:   "00000000",
				Trigger: "抽签",
				Reply:   "少女折寿中......",
				Types:   "车万",
				Limit:   "池子",
			},
			{
				Group:   "1048452984",
				Trigger: "运势测试",
				Reply:   "收到命令！",
				Types:   "李清歌",
				Limit:   "关",
			},
		},
	}
}

func Load(p string) *JsonConfig {
	if !PathExists(p) {
		log.Warnf("尝试加载配置文件 %v 失败: 文件不存在", p)
		return nil
	}
	c := JsonConfig{}
	err := json.Unmarshal([]byte(ReadAllText(p)), &c)
	if err != nil {
		log.Warnf("尝试加载配置文件 %v 时出现错误: %v", p, err)
		log.Infoln("原文件已备份")
		os.Rename(p, p+".backup"+strconv.FormatInt(time.Now().Unix(), 10))
		return nil
	}
	return &c
}

func (c *JsonConfig) Save(p string) error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	WriteAllText(p, string(data))
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ReadAllText(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func WriteAllText(path, text string) {
	_ = ioutil.WriteFile(path, []byte(text), 0644)
}
