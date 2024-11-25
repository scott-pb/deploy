package router

import (
	"deploy/api"
	"deploy/config"
	"deploy/log"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"io/fs"
	"math/rand"
	"net"
	"net/http"
)

//go:embed  static
var front embed.FS

func InitRouter() {
	gin.DefaultWriter = io.MultiWriter(log.Writer)
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.Use(gin.Recovery())
	ip := getLocalIp()
	if config.Config.Ip != "" {
		ip = config.Config.Ip
	}
	log.Info(ip)
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar,XMLHttpRequest,language,token") // 跨域关键设置 让浏览器可以解析
		c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                                                         // 缓存请求信息 单位为秒
		c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                                                 //  跨域请求是否需要带cookie信息 默认设置为true
		//c.Header("content-type", "application/json")                                                                                                                                                                                         // 设置返回格式是json
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.POST("/login", login)
	r.GET("/webSocket", (&api.WebSocket{}).WebSocketHandler)
	staticFp, _ := fs.Sub(front, "static")
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(staticFp))))
	//r.StaticFS("/public", http.FS(front))
	templ := template.Must(template.New("").ParseFS(front, "static/*.html"))
	r.SetHTMLTemplate(templ)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"host": ip + ":" + config.Config.Port})
	})

	_ = r.Run("0.0.0.0:" + config.Config.Port)
}

func login(c *gin.Context) {
	var (
		isSuccess bool
		account   config.Account
	)

	if err := c.ShouldBind(&account); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	for _, acc := range config.Config.Accounts {
		if acc.Username == account.Username && acc.Password == account.Password {
			isSuccess = true
			break
		}
	}

	if !isSuccess {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "账号或密码不正确",
		})
		return
	}
	id := fmt.Sprintf("%d", 1000+rand.Intn(9000))
	config.Config.Sessions[account.Username] = id
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"token":   account.Username + "-" + id,
		"message": "登录成功",
	})
}

func getLocalIp() (ip string) {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && isPrivateIP(ipnet.IP) {
				return ipnet.IP.String()
			}
		}
	}

	return ip
}

func isPrivateIP(ip net.IP) bool {
	privateIPBlocks := []*net.IPNet{
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
