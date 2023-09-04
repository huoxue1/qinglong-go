package server

import (
	"bytes"
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/huoxue1/qinglong-go/service/config"
	"github.com/huoxue1/qinglong-go/utils"
)

func Run(ctx context.Context, path string, query map[string]string, body map[string]any, header map[string]string, isLog bool) string {
	cmd := getCommand(path)
	buffer := bytes.NewBuffer([]byte(""))
	option := &utils.RunOption{
		Ctx:     ctx,
		Command: cmd,
		Env:     getEnv(query, body, header),
		OnStart: func(ctx context.Context) {

		},
		OnEnd: func(ctx context.Context) {

		},
		LogFile: buffer,
		CmdDir:  "./data/scripts/",
	}
	utils.RunWithOption(ctx, option)
	if isLog {
		return string(buffer.Bytes())
	}
	return getData(buffer.Bytes())
}

func getCommand(path string) (cmd string) {

	pythonCmd := config.GetKey("PythonCmd", "python")
	JsCmd := config.GetKey("JsCmd", "node")
	ShCmd := config.GetKey("ShCmd", "bash")

	if strings.HasSuffix(path, ".py") {
		cmd = pythonCmd + " " + path
	} else if strings.HasSuffix(path, ".js") {
		cmd = JsCmd + " " + path
	} else if strings.HasSuffix(path, ".sh") {
		cmd = ShCmd + " " + path
	}
	return
}

func getData(data []byte) (result string) {
	resultReg := `>>>(.*?)<<<`
	compile := regexp.MustCompile(resultReg)
	if compile.Match(data) {
		result = string(compile.FindAllSubmatch(data, -1)[0][1])
	}
	return result
}

func getEnv(query map[string]string, body map[string]any, header map[string]string) map[string]string {
	queryData, _ := json.Marshal(query)
	bodyData, _ := json.Marshal(body)
	headerData, _ := json.Marshal(header)
	return map[string]string{
		"QL_HEADERS": string(headerData),
		"QL_BODY":    string(bodyData),
		"QL_QUERY":   string(queryData),
	}
}
