package config

import (
	"os"
	"path"
	"regexp"
)

func GetKey(key string) string {
	file, err := os.ReadFile(path.Join("data", "config", "config.sh"))
	if err != nil {
		return ""
	}
	compile := regexp.MustCompile(key + `="(.*?)"`)
	if !compile.Match(file) {
		return ""
	}
	datas := compile.FindAllStringSubmatch(string(file), 1)
	return datas[0][1]
}
