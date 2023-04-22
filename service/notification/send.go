package notification

import (
	"errors"
	"fmt"
	log "github.com/huoxue1/go-utils/base/log"
	"github.com/tidwall/gjson"
	"net/http"
	"net/url"
	"strings"
)

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

func (s *ServerChan) Send(title, message string) error {
	_, err := client.R().SetFormData(map[string]string{
		"text": title,
		"desp": strings.ReplaceAll(message, "\n", "\n\n"),
	}).Post(fmt.Sprintf("https://sc.ftqq.com/%s.send", s.ServerChanKey))
	if err != nil {
		return err
	}
	return nil
}

func (p *pushDeer) Send(title, message string) error {
	_, err := client.R().SetFormData(map[string]string{
		"pushKey": p.PushDeerKey,
		"text":    title,
		"desp":    message,
		"type":    "markdown",
	}).Post("https://api2.pushdeer.com/message/push")
	if err != nil {
		return err
	}
	return nil
}

func (g *gotify) Send(title, message string) error {
	_, err := client.R().SetQueryParam("token", g.GotifyToken).SetFormData(map[string]string{
		"title":    title,
		"message":  message,
		"priority": g.GotifyPriority,
	}).Post(fmt.Sprintf("%s/message", g.GotifyUrl))
	if err != nil {
		return err
	}
	return nil
}
