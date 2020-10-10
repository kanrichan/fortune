package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	// 基于验证的key
	au_key := "test"

	// 基于验证的时间
	au_time := fmt.Sprintf("%v", time.Now().Unix())

	// 将验证的key与时间合并成一个字符
	au_key_time := fmt.Sprintf("%v|%v", au_key, au_time)

	m := md5.New()
	m.Write([]byte(au_key_time))
	authkey := hex.EncodeToString(m.Sum(nil))

	// 将生成加密的 KEY 与 时间传递至服务端
	api := "http://www.kanri.ml:10086/fortune.jpg"
	data := url.Values{}
	data.Set("fromQQ", "test")
	data.Set("types", "诺亚幻想")

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
	req.Header.Set("authkey", authkey)
	req.Header.Set("autime", au_time)

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
