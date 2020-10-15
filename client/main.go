package main

import (
	"github.com/Yiwen-Chan/qq-bot-api"
	log "github.com/sirupsen/logrus"

	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type FromDataParm struct {
	Client  string
	Version string
	Types   string
	FromQQ  string
	Ask     string
	Limit   string
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
	Host    string           `json:"host"`
	Port    uint16           `json:"port"`
	Master  int64            `json:"master"`
	Trigger string           `json:"trigger"`
	Setting []*SettingConfig `json:"setting"`
}

type SettingConfig struct {
	Group int64  `json:"group"`
	Types string `json:"types"`
	Limit string `json:"limit"`
}

func DefaultConfig() *JsonConfig {
	return &JsonConfig{
		Host:    "127.0.0.1",
		Port:    8000,
		Master:  12345678,
		Trigger: "运势",
		Setting: []*SettingConfig{
			{
				Group: 0,
				Types: "车万",
				Limit: "on",
			},
			{
				Group: 87654321,
				Types: "碧蓝幻想",
				Limit: "on",
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

func PathExecute() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dir)

	return dir + "/"
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
			if update.Message.Text == conf.Trigger {
				groupID := update.GroupID
				userID := update.UserID

				log.Printf("[%s] %s", update.Message.From.String(), update.Message.Text)

				setting := returnSetting(conf, groupID)
				types := setting.Types
				fromQQ := fmt.Sprintf("%v", userID)
				ask := conf.Trigger
				limit := setting.Limit

				fromDataParm := &FromDataParm{
					"go",
					"3",
					types,
					fromQQ,
					ask,
					limit,
				}
				authkey, au_time := key("test")
				headerParm := &HeaderParm{
					authkey,
					au_time,
				}

				apiFortune := "http://www.kanri.ml:10086/fortune"
				apiPic := "http://www.kanri.ml:10086/fortune.jpg"
				fortuneJson := fortune(apiFortune, fromDataParm, headerParm)

				Message := ""
				path := PathExecute() + "test.jpg"
				if fortuneJson.Code == 200 {
					if fortuneJson.Msg == "success" {
						if fortuneJson.Warn != "" {
							Message = fortuneJson.Warn
							pic(apiPic, fromDataParm, headerParm)
							Message += "[CQ:image,file=file:///" + path + "]"
							bot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, Message)
						} else {
							if fortuneJson.Info != "" {
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
					} else {
						Message = fortuneJson.Msg
						bot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, Message)
					}
				} else {
					//服务器无响应
					Message = "服务器失去连接，请到GitHub提交issue"
					bot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, Message)
				}

			}
		}
	}
}

func returnSetting(conf *JsonConfig, groupID int64) *SettingConfig {
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

func fortune(api string, fromDataParm *FromDataParm, headerParm *HeaderParm) FortuneJson {
	data := url.Values{}
	data.Set("client", fromDataParm.Client)
	data.Set("version", fromDataParm.Version)
	data.Set("types", fromDataParm.Types)
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
	fortuneJson := FortuneJson{}
	err = json.Unmarshal(body, &fortuneJson)
	if err != nil {
		panic(err)
	}
	return fortuneJson
}

func pic(api string, fromDataParm *FromDataParm, headerParm *HeaderParm) {
	data := url.Values{}
	data.Set("client", fromDataParm.Client)
	data.Set("version", fromDataParm.Version)
	data.Set("types", fromDataParm.Types)
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

	f, err := os.OpenFile("test.jpg", os.O_WRONLY|os.O_TRUNC, 0600)
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
