package fortune

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type FromDataStruct struct {
	client    string
	version   string
	bot       string
	types     string
	fromGroup string
	fromQQ    string
	ask       string
	limit     string
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

func getFromData(fromDataStruct *FromDataStruct) *strings.Reader {
	data := url.Values{}
	data.Set("client", fromDataStruct.client)
	data.Set("version", fromDataStruct.version)
	data.Set("bot", fromDataStruct.bot)
	data.Set("types", fromDataStruct.types)
	data.Set("fromGroup", fromDataStruct.fromGroup)
	data.Set("fromQQ", fromDataStruct.fromQQ)
	data.Set("ask", fromDataStruct.ask)
	data.Set("limit", fromDataStruct.limit)
	return strings.NewReader(data.Encode())
}

func getHeader(clientKey string) *HeaderParm {
	authkey, au_time := key(clientKey)
	headerParm := &HeaderParm{
		authkey,
		au_time,
	}
	return headerParm
}

func key(au_key string) (string, string) {
	au_time := fmt.Sprintf("%v", time.Now().Unix())
	au_key_time := fmt.Sprintf("%v|%v", au_key, au_time)

	m := md5.New()
	m.Write([]byte(au_key_time))
	authkey := hex.EncodeToString(m.Sum(nil))
	return authkey, au_time
}

func fortune(api string, fromDataStruct *FromDataStruct, headerParm *HeaderParm) (FortuneJson, int) {
	transport := http.Transport{
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: &transport,
	}

	fromData := getFromData(fromDataStruct)

	req, err := http.NewRequest("POST", api, fromData)
	if err != nil {
		OutPutLog("[fortune-运势] 创建POST请求失败")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authkey", headerParm.Authkey)
	req.Header.Set("autime", headerParm.Au_time)

	resp, err := client.Do(req)
	if err != nil {
		OutPutLog("[fortune-运势] POST请求失败")
	}

	defer resp.Body.Close()

	code := resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		OutPutLog("[fortune-运势] 读取POST请求数据失败")
	}
	fortuneJson := FortuneJson{}

	err = json.Unmarshal(body, &fortuneJson)
	if err != nil {
		OutPutLog("[fortune-运势] 解析JSON失败")
	}

	f, err := os.OpenFile(ResultPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		OutPutLog("[fortune-运势] 写入请求结果文件失败")
	} else {
		_, err = f.Write([]byte(string(body)))
	}

	return fortuneJson, code
}

func pic(api string, fromDataStruct *FromDataStruct, headerParm *HeaderParm) {
	transport := http.Transport{
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: &transport,
	}

	fromData := getFromData(fromDataStruct)

	req, err := http.NewRequest("POST", api, fromData)
	if err != nil {
		OutPutLog("[fortune-运势] 创建POST请求失败")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authkey", headerParm.Authkey)
	req.Header.Set("autime", headerParm.Au_time)

	resp, err := client.Do(req)
	if err != nil {
		OutPutLog("[fortune-运势] POST请求失败")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		OutPutLog("[fortune-运势] 读取POST请求数据失败")
	}

	f, err := os.OpenFile(PicPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		OutPutLog("[fortune-运势] 写入图片文件失败")
	} else {
		_, err = f.Write([]byte(string(body)))
	}
}
