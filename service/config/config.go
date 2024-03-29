package config

import (
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var VERSION = "v1.0.0"

func GetKey(key, defaultValue string) string {
	file, err := os.ReadFile(path.Join("data", "config", "config.sh"))
	if err != nil {
		return defaultValue
	}
	compile := regexp.MustCompile(key + `="(.*?)"`)
	if !compile.Match(file) {
		return defaultValue
	}
	datas := compile.FindAllStringSubmatch(string(file), 1)
	return datas[0][1]
}

func GetVersion() string {
	return strings.TrimPrefix(VERSION, "v")
}

var address string

func SetAddress(add string) {
	address = add
}

func ListenPort() int {
	split := strings.Split(address, ":")
	port, _ := strconv.Atoi(split[1])
	return port
}
