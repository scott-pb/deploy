package api

import (
	"deploy/config"
	"deploy/constant"
	"deploy/internal"
	"deploy/log"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WebSocket struct {
}

var mu sync.Mutex

func (w *WebSocket) WebSocketHandler(c *gin.Context) {
	upgrade := websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
	}
	var isLogin bool
	token := c.Request.URL.RawQuery
	if token != "" {
		s := strings.Split(token, "-")
		if len(s) == 2 && s[1] == config.Config.Sessions[s[0]] {
			isLogin = true
		}
	}

	// 获取WebSocket连接
	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}
	clientIP := c.ClientIP()
	fmt.Printf("Client IP: %s\n", clientIP)
	s := internal.DeployService{}
	// 处理WebSocket消息
	for {
		if err = ws.SetReadDeadline(time.Now().Add(time.Second * 20)); err != nil {
			log.Error("SetReadDeadline err", err)
			return
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if !isLogin {
			_ = ws.WriteMessage(websocket.TextMessage, []byte("no login"))
			continue
		}

		if string(p) == "ping" {
			_ = ws.WriteMessage(websocket.TextMessage, []byte("pong"))
			continue
		}
		mu.Lock()

		var message internal.Message
		_ = json.Unmarshal(p, &message)
		log.Info(clientIP, fmt.Sprintf("%+v", message))
		switch message.Project {
		case constant.Admin:
			switch message.Env {
			case constant.Test:
				s.AdminTest(ws, message)
			case constant.Release:
				s.AdminRelease(ws, message)
			}
		case constant.Enterprise:
			switch message.Env {
			case constant.Test:
				s.EnterpriseTest(ws, message)
			case constant.Release:
				s.EnterpriseRelease(ws, message)
			}
		case constant.Server:
			switch message.Env {
			case constant.Test:
				s.ServerTest(ws, message)
			case constant.Release:
				s.ServerRelease(ws, message)
			}
		}
		mu.Unlock()
	}

	// 关闭WebSocket连接
	_ = ws.Close()
}
