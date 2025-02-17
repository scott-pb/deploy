// Package dingtalk 钉钉API
package dingtalk

import (
	"deploy/notify/common"
	"errors"
	"fmt"
	"github.com/blinkbean/dingtalk"
	"log"
	"strings"
	"time"
)

type Handler struct {
	Client *dingtalk.DingTalk
}

func NewInit(tokens string, secret string) (common.Notify, error) {
	if tokens == "" || secret == "" {
		return nil, errors.New("init dingtalk failed with empty configure")
	}

	return &Handler{
		Client: dingtalk.InitDingTalkWithSecret(tokens, secret),
	}, nil
}

func Init(tokens string, secret string) (common.Notify, error) {
	return NewInit(tokens, secret)
}

func (h *Handler) Send(messages common.Messages) (err error) {

	ms := []string{fmt.Sprintf("%s：%s", "当前时间", time.Now().Format(time.DateTime)}
	for _, item := range messages {
		ms = append(ms, fmt.Sprintf("%s：%s", item.Name, item.Value))
	}
	if err = h.Client.SendTextMessage(strings.Join(ms, "\n")); err != nil {
		log.Println("notify dingtalk err:", err)
		return err
	}

	return nil
}
