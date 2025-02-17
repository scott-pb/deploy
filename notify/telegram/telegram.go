package telegram

import (
	"bytes"
	"depends/constant"
	"depends/models"
	"depends/notify/common"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Telegram struct {
	Url        string
	ServerName string
	ChatId     int64
	common.Notification
	builder    strings.Builder
	header     string
	httpClient *http.Client
}

type SendMessageRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func Init(sn, botToken, proxy string, cid int64) (common.Notify, error) {
	if botToken == "" || cid == 0 {
		return nil, errors.New("telegram config botToken or cid is nil")
	}
	telegram := &Telegram{
		Url:        "https://api.telegram.org/bot" + botToken + "/sendMessage",
		ServerName: sn,
		ChatId:     cid,
		builder:    strings.Builder{},
		header:     "<b>应用名称:</b>" + sn + "\n",
		httpClient: &http.Client{Timeout: time.Second},
	}
	if proxy != "" {
		u, _ := url.Parse(proxy)
		telegram.httpClient = &http.Client{
			Timeout: time.Second * 3,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(u),
			},
		}
	}
	return telegram, nil
}

func (t *Telegram) Send(message models.Messages) error {
	defer func() {
		t.builder.Reset()
		time.Sleep(time.Millisecond * 500)
	}()
	t.builder.WriteString(t.header)
	t.builder.WriteString("<b>当前时间:</b>")
	t.builder.WriteString(time.Now().Format(constant.DateTimeMilli))
	t.builder.WriteString("\n")

	for _, value := range message {
		t.builder.WriteString("<b>")
		t.builder.WriteString(value.Name)
		t.builder.WriteString(":</b>")
		t.builder.WriteString(value.Value)
		t.builder.WriteString("\n")
	}

	requestBody, err := json.Marshal(SendMessageRequest{
		ChatID:    t.ChatId,
		Text:      t.builder.String(),
		ParseMode: "HTML",
	})

	resp, err := t.httpClient.Post(t.Url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Failed to send message, status code: %d", resp.StatusCode))
	}
	return err
}
