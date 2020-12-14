package fortune

// 全局变量初始化
var (
	AppInfoJson string

	FirstStart bool = true

	ApiHost    = "api.kanri.top"
	ApiPort    = "10086"
	ApiFortune = "http://" + ApiHost + ":" + ApiPort + "/fortune"
	ApiPic     = "http://" + ApiHost + ":" + ApiPort + "/fortune.jpg"

	ClientKey  = "233666"
	ClientName = "xq"
	ClientVer  = "5"

	AppPath    = PathExecute() + "data/app/fortune/"
	ConfPath   = AppPath + "config.yml"
	PicPath    = AppPath + "output.jpg"
	ResultPath = AppPath + "output.txt"

	Conf = &YamlConfig{}
)

func init() {}

func Main() {}

func onStart() {
	CreatePath(AppPath)
	if FirstStart {
		INFO("[OneBot-YaYa] 夜夜は世界一かわいい")
		INFO("[fortune-运势] 项目地址：https://github.com/Yiwen-Chan/fortune")
		INFO("[fortune-运势] 配置文件：%s%s", AppPath, "config.yml")
		Conf = Load(AppPath + "config.yml")
		if Conf == nil {
			ERROR("晚安~")
			return
		}
	}
	FirstStart = false
}
