package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/huoxue1/qinglong-go/utils"
	"github.com/imroc/req/v3"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"os"
	"path"
)

type Pusher interface {
	Send(title, message string) error
}

var (
	Push   Pusher
	client *req.Client
)

func init() {
	client = req.C()
}

func init() {
	if utils.FileExist(path.Join("data", "config", "push.json")) {
		data, _ := os.ReadFile(path.Join("data", "config", "push.json"))
		push, err := parsePush(string(data))
		if err != nil {
			return
		}
		Push = push
	} else {
		Push = &defaultPush{}
	}
	_ = Push.Send("上线通知", "你的青龙已上线！")
}

func parsePush(config string) (Pusher, error) {
	t := gjson.Get(config, "type").String()
	log.Infoln("采用全局推送方式： ", t)
	switch t {
	case "telegramBot":
		return getBot[TgBot](config)

	case "goCqHttpBot":
		return getBot[goCqHttpBot](config)

	case "serverChan":
		return getBot[ServerChan](config)

	case "pushDeer":
		return getBot[pushDeer](config)

	case "gotify":
		return getBot[gotify](config)

	default:
		return nil, errors.New("not found type")
	}
}

func getBot[T any](config string) (*T, error) {
	bot := new(T)
	err := json.Unmarshal([]byte(config), bot)
	if err != nil {
		return nil, err
	}
	return bot, nil
}

type defaultPush struct {
}

func (d *defaultPush) Send(title, message string) error {
	log.Infoln(fmt.Sprintf("[push] %s,%s", title, message))
	return nil
}

type TgBot struct {
	Type                 string `json:"type"`
	TelegramBotToken     string `json:"telegramBotToken"`
	TelegramBotUserId    string `json:"telegramBotUserId"`
	TelegramBotProxyHost string `json:"telegramBotProxyHost"`
	TelegramBotProxyPort string `json:"telegramBotProxyPort"`
	TelegramBotProxyAuth string `json:"telegramBotProxyAuth"`
	TelegramBotApiHost   string `json:"telegramBotApiHost"`
}

type goCqHttpBot struct {
	Type             string `json:"type"`
	GoCqHttpBotUrl   string `json:"goCqHttpBotUrl"`
	GoCqHttpBotToken string `json:"goCqHttpBotToken"`
	GoCqHttpBotQq    string `json:"goCqHttpBotQq"`
}

type ServerChan struct {
	Type          string `json:"type"`
	ServerChanKey string `json:"serverChanKey"`
}

type pushDeer struct {
	PushDeerKey string `json:"pushDeerKey"`
}

type gotify struct {
	GotifyUrl      string `json:"gotifyUrl"`
	GotifyToken    string `json:"gotifyToken"`
	GotifyPriority string `json:"gotifyPriority"`
}
