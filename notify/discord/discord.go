package discord

import (
	"deploy/notify/common"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	once sync.Once
)

type Handler struct {
	Name    string
	HookURL string
	Token   string
	Bot     *discordgo.Session
}

func Init(name, hookUrl, proxy, token string) (common.Notify, error) {
	if hookUrl == "" || token == "" {
		log.Println("NewInit discord failed with empty configure")
		return nil, errors.New("NewInit discord failed with empty configure")
	}
	fmt.Println(name, hookUrl, token, proxy)
	var err error
	defer func() {
		if err != nil {
			log.Println("NewInit discord failed", err)
		}
	}()
	c, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	if proxy != "" {
		u, _ := url.Parse(proxy)
		c.Client = &http.Client{
			Timeout: time.Second * 3,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(u),
			},
		}
		c.Dialer = &websocket.Dialer{
			Proxy: http.ProxyURL(u),
		}
	}
	if err = c.Open(); err != nil {
		return nil, err
	}
	return &Handler{Name: name, HookURL: hookUrl, Token: token, Bot: c}, nil

}

func (h *Handler) Send(messages common.Messages) error {
	msg := fmt.Sprintf("*%s*: %s \n", "当前时间", time.Now().Format(time.DateTime))
	for _, vo := range messages {
		if !vo.IsCode {
			msg += fmt.Sprintf("*%s*: %s \n", vo.Name, vo.Value)
		} else {
			msg += fmt.Sprintf("*%s*: ```%s``` \n", vo.Name, vo.Value)
		}
	}
	if _, err := h.Bot.ChannelMessageSendComplex(h.HookURL, &discordgo.MessageSend{
		Content: msg,
		TTS:     false,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Hello, Embed!",
				Description: "This is an embedded message.",
			},
		},
	}); err != nil {
		return err
	}
	return nil
}
