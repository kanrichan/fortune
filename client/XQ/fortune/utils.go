package fortune

import (
	"io/ioutil"
	"os"
)

func PathExecute() string {
	dir, err := os.Getwd()
	if err != nil {
		OutPutLog("[fortune-运势] 判断当前运行路径失败")
	}
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

func CreatePath(path string) error {
	err := os.MkdirAll(path, 0644)
	if err != nil {
		return err
	}
	return nil
}
