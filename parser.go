package golia

import (
	"io/ioutil"
	"log"

	"strings"

	"github.com/astaxie/beego/config"
)

const defaultConfName = "config"

var Confs map[string]config.Configer

func init() {
	Confs = make(map[string]config.Configer)
	// 读取conf下所有文件
	files, err := ioutil.ReadDir("../conf")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		name, ext := sepFileNameAndExt(file.Name())
		if ext == "json" || ext == "ini" {
			Confs[name], err = config.NewConfig(ext, "../conf/"+file.Name())
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func sepFileNameAndExt(fullName string) (name string, ext string) {
	if !strings.Contains(fullName, ".") {
		name = fullName
		return
	}
	splited := strings.Split(fullName, ".")
	name, ext = splited[0], splited[1]
	return
}

func DefaultConf() config.Configer {
	return Confs[defaultConfName]
}
