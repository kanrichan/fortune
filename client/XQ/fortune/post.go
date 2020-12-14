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

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authkey", headerParm.Authkey)
	req.Header.Set("autime", headerParm.Au_time)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		ERROR("[GET INFO] 网络错误 ❌")
	}

	code := resp.StatusCode

	body, _ := ioutil.ReadAll(resp.Body)
	fortuneJson := FortuneJson{}
	_ = json.Unmarshal(body, &fortuneJson)

	f, err := os.OpenFile(ResultPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		ERROR("[GET INFO] 写入信息错误 ❌")
	} else {
		f.Write(body)
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

	req, _ := http.NewRequest("POST", api, fromData)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authkey", headerParm.Authkey)
	req.Header.Set("autime", headerParm.Au_time)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		ERROR("[GET PIC] 网络错误 ❌")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	f, err := os.OpenFile(PicPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		ERROR("[GET PIC] 写入图片错误 ❌")
	} else {
		f.Write(body)
	}
}
