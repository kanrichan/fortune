package fortune

import (
	"fmt"
	"io/ioutil"
	"os"

	"yaya/core"
)

func PathExecute() string {
	dir, _ := os.Getwd()
	return dir + "/"
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

func CreatePath(path string) {
	if !PathExists(path) {
		err := os.MkdirAll(path, 0644)
		if err != nil {
			ERROR("生成应用目录失败")
		}
	}
}

func ProtectRun(entry func(), label string) {
	defer func() {
		err := recover()
		if err != nil {
			ERROR("[协程] %v协程发生了不可预知的错误，请在GitHub提交issue：%v", label, err)
		}
	}()
	entry()
}

func INFO(s string, v ...interface{}) {
	core.OutPutLog("[INFO] " + fmt.Sprintf(s, v...))
}

func ERROR(s string, v ...interface{}) {
	core.OutPutLog("[ERROR] " + fmt.Sprintf(s, v...))
}
