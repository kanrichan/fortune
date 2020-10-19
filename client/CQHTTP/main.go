package main

import (
	"github.com/Yiwen-Chan/qq-bot-api"
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
	Host    string           `json:"WS服务器"`
	Port    uint16           `json:"WS端口"`
	Bot     int64            `json:"机器人QQ"`
	Setting []*SettingConfig `json:"运势设置"`
}

type SettingConfig struct {
	Group   string `json:"群号"`
	Trigger string `json:"触发"`
	Reply   string `json:"回复"`
	Types   string `json:"类型"`
	Limit   string `json:"限制"`
}

func main() {
	log.Printf("Fortune-运势 正在启动")
	log.Printf("项目地址：https://github.com/Yiwen-Chan/fortune")
	var conf *JsonConfig
	if PathExists("config.json") {
		conf = Load("config.json")
	}
	if conf == nil {
		err := DefaultConfig().Save("config.json")
		if err != nil {
			log.Fatalf("创建默认配置文件时出现错误: %v", err)
			return
		}
		log.Infof("默认配置文件已生成, 请编辑 config.json 后重启程序.")
		time.Sleep(time.Second * 5)
		return
	}
	log.Printf("Fortune-运势 加载配置完毕")

	// Whether to use WebSocket or LongPolling depends on the address.
	// To use WebSocket, the address should be something like "ws://localhost:6700"
	url := fmt.Sprintf("ws://%v:%v", conf.Host, conf.Port)
	bot, err := qqbotapi.NewBotAPI("", url, "")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	u := qqbotapi.NewUpdate(0)
	u.PreloadUserInfo = true
	updates, err := bot.GetUpdatesChan(u)

	log.Printf("Fortune-运势 启动完毕，正在运行")

	for update := range updates {
		//message类信息触发
		if update.PostType == "message" {

			groupID := update.GroupID
			userID := update.UserID
			text := update.Message.Text

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
			}

			types = Types

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

				log.Printf("[%s] %s", update.Message.From.String(), text)
				if reply != "" {
					bot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, reply)
				}

				client := "go"
				version := "4"
				botQQ := fmt.Sprintf("%v", conf.Bot)
				types := types
				fromGroup := fmt.Sprintf("%v", groupID)
				fromQQ := fmt.Sprintf("%v", userID)
				ask := text
				limit := limit

				fromDataParm := &FromDataParm{
					client,
					version,
					botQQ,
					types,
					fromGroup,
					fromQQ,
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
				path := PathExecute() + "output.jpg"

				if code != 200 {
					//服务器无响应
					Message = "[fortune-运势] 服务器失联中......"
					Message += fmt.Sprintf(" code: %v", code)
					bot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, Message)
				} else if fortuneJson.Code != 200 {
					//服务器状态异常
					Message = fortuneJson.Msg
					bot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, Message)
				} else if fortuneJson.Msg != "success" {
					Message = fortuneJson.Msg
					bot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, Message)
				} else if fortuneJson.Warn != "" {
					Message = fortuneJson.Warn
					pic(apiPic, fromDataParm, headerParm)
					Message += "[CQ:image,file=file:///" + path + "]"
					bot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, Message)
				} else if fortuneJson.Info != "" {
					if notSend() {
						Message = fortuneJson.Info
					}
					pic(apiPic, fromDataParm, headerParm)
					Message += "[CQ:image,file=file:///" + path + "]"
					bot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, Message)
				} else {
					pic(apiPic, fromDataParm, headerParm)
					Message += "[CQ:image,file=file:///" + path + "]"
					bot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, Message)
				}
			}
		}
	}
}

func returnSetting(conf *JsonConfig, groupID int64) *SettingConfig {
	group := fmt.Sprintf("%v", groupID)
	setting := conf.Setting[0]
	for index, _ := range conf.Setting {
		settingCache := conf.Setting[index]
		if group == settingCache.Group {
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
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authkey", headerParm.Authkey)
	req.Header.Set("autime", headerParm.Au_time)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	code := resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fortuneJson := FortuneJson{}
	err = json.Unmarshal(body, &fortuneJson)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authkey", headerParm.Authkey)
	req.Header.Set("autime", headerParm.Au_time)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("output.jpg", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
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
		Host: "127.0.0.1",
		Port: 8000,
		Bot:  12345678,
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
