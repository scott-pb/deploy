package api

import (
	"deploy/config"
	"deploy/constant"
	"deploy/internal"
	"deploy/log"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocket struct {
}

var mapMu = make(map[string]*sync.Mutex)
var mapWsMu = make(map[*websocket.Conn]map[string]*sync.Mutex)

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

	if !isLogin {
		s.Flush("no login", ws)
		return
	}

	defer func() {
		log.Error("ws Close", err)
		_ = ws.Close()
		if v, ok := mapWsMu[ws]; ok {
			for _, mv := range v {
				mv.Unlock()
			}
		}

	}()
	// 处理WebSocket消息
	for {
		_ = ws.SetReadDeadline(time.Now().Add(time.Second * 30))
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Error("ReadMessage err", err)
			break
		}
		if !isLogin {
			s.Flush("no login", ws)
			continue
		}

		if string(p) == "ping" {
			s.Flush("pong", ws)
			continue
		}

		var message internal.Message
		_ = json.Unmarshal(p, &message)
		message.UserName = username
		key := message.Project + "_" + message.Env
		if message.Env == constant.Production {
			key = constant.Production
		}
		v, ok := mapMu[key]
		if !ok {
			v = &sync.Mutex{}
			mapMu[key] = v
		}
		for !v.TryLock() {
			s.Flush(message.Project+" "+message.Env+" 运行中...", ws)
			time.Sleep(time.Second)
		}

		_, wok := mapWsMu[ws]
		if !wok {
			mapWsMu[ws] = make(map[string]*sync.Mutex)
		}
		mapWsMu[ws][key] = v

		log.Info(clientIP, username, fmt.Sprintf("%+v", message))
		go func() {
			defer func() {
				s.Flush("finished", ws)
				v.Unlock()
				delete(mapWsMu[ws], key)
				fmt.Println("close==========")
			}()
			if message.Env == constant.Production {
				s.ServerProduction(ws, message)
				return
			}

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
			case constant.Front:
				switch message.Env {
				case constant.Test:
					s.ServerTest(ws, message)
				case constant.Release:
					s.ServerRelease(ws, message)
				}
			}
		}()
	}

}
