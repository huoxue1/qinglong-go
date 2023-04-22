package fileInit

import (
	_ "embed"
	"os"
	"path"
)

//go:embed config_sample.sh
var sample []byte

func InitConfig() {
	_ = os.MkdirAll(path.Join("data", "config"), 0666)
	_ = os.WriteFile(path.Join("data", "config", "config.sample.sh"), sample, 0666)
	_ = os.WriteFile(path.Join("data", "config", "config.sh"), sample, 0666)
}
