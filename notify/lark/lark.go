package lark

import (
	"deploy/notify/common"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	client   *http.Client
	template = `{"card":{"elements":[{"tag":"div","text":{"content":"%s","tag":"lark_md"},"fields":[]}],"config":{"wide_screen_mode":true},"header":{"template":"blue","title":{"content":"%s","tag":"plain_text"}}},"msg_type":"interactive"}`
)

// Handler https://open.larksuite.com/open-apis/bot/v2/hook/81ccbebe-0bee-435c-b50c-11654637cce9
type Handler struct {
	Url string
}

func Init(urlStr, proxy string) (common.Notify, error) {
	if len(urlStr) == 0 {
		return nil, errors.New("url is not empty")
	}
	client = &http.Client{Timeout: time.Second * 10}
	if proxy != "" {
		u, _ := url.Parse(proxy)
		client = &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(u),
			},
		}
	}
	return &Handler{Url: urlStr}, nil
}

func (h *Handler) Send(messages common.Messages) error {
	var (
		title   string
		content string
	)

	content = fmt.Sprintf("%s：%s\\n", "**当前时间**", time.Now().Format(time.DateTime))
	for _, item := range messages {
		if len(title) == 0 && item.Name == "服务名称" {
			title = item.Value
			continue
		}
		if item.Name == "程序版本" {
			item.Value = "<font color='green'>**" + item.Value + "**</font>"
		}

		if strings.HasPrefix(item.Name, "错误") || strings.Contains(item.Name, "SQL") || item.Name == "GIT信息" {
			item.Value = "<font color='red'>" + item.Value + "</font> "
		}
		item.Value = strings.ReplaceAll(item.Value, "\n\t", "\\n")
		item.Value = strings.ReplaceAll(item.Value, "\n", "\\n")
		item.Value = strings.ReplaceAll(item.Value, "`", "")
		content += fmt.Sprintf("**%s**：%s\\n", item.Name, item.Value)
	}
	str := fmt.Sprintf(template, content, title)
	resp, err := client.Post(h.Url, "application/json", strings.NewReader(str))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	b, _ := io.ReadAll(resp.Body)
	return errors.New(string(b))
}
