package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/huoxue1/qinglong-go/utils"
	"github.com/imroc/req/v3"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
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

func (g *goCqHttpBot) Send(title, message string) error {
	sendUrl, err := url.Parse(g.GoCqHttpBotUrl)
	if err != nil {
		return err
	}
	if strings.Contains(sendUrl.Path, "send_private_msg") {
		sendUrl.Path = ""
		resp, err := client.R().SetHeader("Authorization", g.GoCqHttpBotToken).SetBodyJsonMarshal(map[string]any{
			"action": "send_private_msg",
			"params": map[string]any{
				"user_id": g.GoCqHttpBotQq,
				"message": map[string]any{
					"type": "text",
					"data": map[string]any{
						"text": fmt.Sprintf("%s\n\n%s", title, message),
					},
				},
			},
		}).Post(sendUrl.String())
		log.Infoln(resp.String())
		if err != nil {
			return err
		}
	} else {
		sendUrl.Path = ""
		resp, err := client.R().SetHeader("Authorization", g.GoCqHttpBotToken).SetBodyJsonMarshal(map[string]any{
			"action": "send_group_msg",
			"params": map[string]any{
				"group_id": g.GoCqHttpBotQq,
				"message": map[string]any{
					"type": "text",
					"data": map[string]any{
						"text": fmt.Sprintf("%s\n\n%s", title, message),
					},
				},
			},
		}).Post(sendUrl.String())
		log.Infoln(resp.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TgBot) Send(title, message string) error {
	if t.TelegramBotApiHost == "" {
		t.TelegramBotApiHost = "api.telegram.org"
	}
	if t.TelegramBotProxyHost != "" {
		client.SetProxyURL(fmt.Sprintf("http://%v@%v:%v", t.TelegramBotProxyAuth, t.TelegramBotProxyHost, t.TelegramBotProxyPort))
	} else {
		client.SetProxy(http.ProxyFromEnvironment)
	}
	response, err := client.R().SetFormData(map[string]string{
		"chat_id":                  t.TelegramBotUserId,
		"text":                     fmt.Sprintf("%s\n\n%s", title, message),
		"disable_web_page_preview": "true",
	}).Post(fmt.Sprintf("https://%s/bot%s/sendMessage", t.TelegramBotApiHost, t.TelegramBotToken))
	if err != nil {
		return err
	}
	if gjson.GetBytes(response.Bytes(), "ok").Bool() {
		return nil
	} else {
		return errors.New(response.String())
	}
}

func parsePush(config string) (Pusher, error) {
	t := gjson.Get(config, "type").String()
	log.Infoln("采用全局推送方式： ", t)
	switch t {
	case "telegramBot":
		{
			bot := new(TgBot)
			err := json.Unmarshal([]byte(config), bot)
			if err != nil {
				return nil, err
			}
			return bot, nil
		}
	case "goCqHttpBot":
		{
			bot := new(goCqHttpBot)
			err := json.Unmarshal([]byte(config), bot)
			if err != nil {
				return nil, err
			}
			return bot, nil
		}

	default:
		return nil, errors.New("not found type")
	}
}
