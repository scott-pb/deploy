package slack

import (
	"deploy/notify/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Handler struct {
	Name    string
	HookURL string
	Proxy   string
}

func (h *Handler) Send(messages common.Messages) error {
	msg := fmt.Sprintf("*%s*: %s \n", "当前时间", time.Now().Format(time.DateTime))
	for _, vo := range messages {
		if vo.IsCode || strings.HasPrefix(vo.Name, "错误堆栈") {
			msg += fmt.Sprintf("*%s*: ```%s``` \n", vo.Name, vo.Value)
			continue
		}

		msg += fmt.Sprintf("*%s*: %s \n", vo.Name, vo.Value)
	}
	data, _ := json.Marshal(gin.H{"text": msg + "\n\n", "mrkdwn": true})

	var hc *http.Client

	if h.Proxy != "" {
		u, _ := url.Parse(h.Proxy)
		hc = &http.Client{
			Timeout: time.Second * 2,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(u),
			},
		}
	} else {
		hc = &http.Client{Timeout: time.Second * 2}
	}

	req, err := http.NewRequest(http.MethodPost, h.HookURL, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}

func Init(name, hookUrl, proxy string) (common.Notify, error) {
	if hookUrl == "" {

		return nil, errors.New("init slack failed with empty hook url")
	}
	return &Handler{Name: name, HookURL: hookUrl, Proxy: proxy}, nil
}
