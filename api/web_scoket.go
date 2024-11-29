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

var mapMu = make(map[string]sync.Mutex)

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
	var (
		isLogin  bool
		username string
	)
	token := c.Request.URL.RawQuery
	if token != "" {
		s := strings.Split(token, "-")
		if len(s) == 2 && s[1] == config.Config.Sessions[s[0]] {
			isLogin = true
			username = s[0]
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
	go func() {
		// 关闭WebSocket连接
		defer ws.Close()
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

			var message internal.Message
			_ = json.Unmarshal(p, &message)
			v, ok := mapMu[message.Project+"_"+message.Env]
			if ok {
				for !v.TryLock() {
					_ = ws.WriteMessage(websocket.TextMessage, []byte("排队中..."))
					time.Sleep(time.Second)
				}
			} else {
				v = sync.Mutex{}
				mapMu[message.Project+"_"+message.Env] = v
			}
			for !v.TryLock() {
				_ = ws.WriteMessage(websocket.TextMessage, []byte(message.Project+" "+message.Env+" 排队中..."))
				time.Sleep(time.Second)
			}

			log.Info(clientIP, username, fmt.Sprintf("%+v", message))
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
			v.Unlock()
		}
	}()

}
