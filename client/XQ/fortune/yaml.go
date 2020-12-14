package fortune

import (
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
	"time"
)

type YamlConfig struct {
	Version string         `yaml:"插件版本"`
	Setting []*GroupConfig `yaml:"运势设置"`
}

type GroupConfig struct {
	Group        string         `yaml:"设置群号"`
	GroupSetting *SettingConfig `yaml:"群聊设置"`
}

type SettingConfig struct {
	Trigger string `yaml:"触发词语"`
	Reply   string `yaml:"等待回复"`
	Types   string `yaml:"卡池类型"`
	Limit   string `yaml:"每日限制"`
	Warm    string `yaml:"超过警告"`
}

func DefaultConfig() *YamlConfig {
	return &YamlConfig{
		Version: ClientVer,
		Setting: []*GroupConfig{
			{
				Group: "默认",
				GroupSetting: &SettingConfig{
					Trigger: "运势",
					Reply:   "少女祈祷中......",
					Types:   "李清歌|碧蓝幻想|公主连结",
					Limit:   "全局文字",
					Warm:    "今天已经抽过了，请明天再来吧~",
				},
			},
			{
				Group: "请阅读项目地址配置文件，单个群设置填群号",
				GroupSetting: &SettingConfig{
					Trigger: "这里填触发关键词",
					Reply:   "这里填收到关键词的回复",
					Types:   "这里是池子类型，多个池子用|分开",
					Limit:   "每天限制一张，可填 全局图片 池子图片 全局文字 池子文字 关",
					Warm:    "若为 全局文字 或 池子文字 ，当天抽过则返回此项",
				},
			},
			{
				Group: "00000000",
				GroupSetting: &SettingConfig{
					Trigger: "抽签",
					Reply:   "少女折寿中......",
					Types:   "车万",
					Limit:   "池子图片",
					Warm:    "少女祈祷中......才怪",
				},
			},
			{
				Group: "1048452984",
				GroupSetting: &SettingConfig{
					Trigger: "运势测试",
					Reply:   "收到命令！",
					Types:   "李清歌",
					Limit:   "关",
					Warm:    "少女祈祷中......才怪",
				},
			},
		},
	}
}

func getSetting(conf *YamlConfig, groupID string) *SettingConfig {
	setting := conf.Setting[0].GroupSetting
	for index, _ := range conf.Setting {
		if groupID == conf.Setting[index].Group {
			setting := conf.Setting[index].GroupSetting
			return setting
		}
	}
	return setting
}

func Load(p string) *YamlConfig {
	if !PathExists(p) {
		c := DefaultConfig()
		c.Save(p)
	}
	c := YamlConfig{}
	err := yaml.Unmarshal([]byte(ReadAllText(p)), &c)
	if err != nil {
		ERROR("Emmm，夜夜觉得配置文件有问题")
		os.Rename(p, p+".backup"+strconv.FormatInt(time.Now().Unix(), 10))
		c := DefaultConfig()
		c.Save(p)
	}
	INFO("おはようございます。")
	c.Save(p)
	return &c
}

func (c *YamlConfig) Save(p string) {
	data, err := yaml.Marshal(c)
	if err != nil {
		ERROR("大失败！夜夜需要管理员权限")
	}
	WriteAllText(p, string(data))
}
