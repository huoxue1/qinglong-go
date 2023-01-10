package notification

import (
	"os"
	"path"
)

func HandlePush(config string) error {
	push, err := parsePush(config)
	if err != nil {
		return err
	}
	err = push.Send("青龙测试", "青龙消息推送测试")
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join("data", "config", "push.json"), []byte(config), 0666)
	if err != nil {
		return err
	}
	Push = push
	return nil
}
