package router

import (
	"deploy/api"
	"deploy/config"
	"deploy/log"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar,XMLHttpRequest,language,token") // 跨域关键设置 让浏览器可以解析
		c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                                                         // 缓存请求信息 单位为秒
		c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                                                 //  跨域请求是否需要带cookie信息 默认设置为true
		c.Header("content-type", "application/json")                                                                                                                                                                                         // 设置返回格式是json
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.POST("/login", login)
	r.GET("/webSocket", (&api.WebSocket{}).WebSocketHandler)
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, "%s", html(ip+":"+config.Config.Port))
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

func html(url string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>soga 发布系统 </title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            display: flex;
            flex-direction: column;
            min-height: 100vh;
            background-color: #f0f0f0;
        }

        .toolbar {
            position: sticky;
            top: 0;
            left: 0;
            right: 0;
            z-index: 1000;
            display: flex;
            justify-content: center;
            background-color: #007BFF;
            color: white;
            padding: 10px;
        }

        .toolbar button {
            padding: 8px 16px;
            font-size: 14px;
            color: white;
            background-color: #5daf34;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            margin: 5px;
        }

        .toolbar button:hover {
            background-color: #0056b3;
        }

        .form-container {
            padding: 20px;
            background-color: white;
            border: 1px solid #ccc;
            border-radius: 4px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
            flex-grow: 1;
            overflow-y: auto;
        }

        .form-container fieldset {
            border: none;
            padding: 0;
            margin: 0;
        }

        .form-container legend {
            font-weight: bold;
            color: #007BFF;
            margin-bottom: 10px;
        }

        .form-container p {
            margin-bottom: 10px;
        }

        .form-container label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #555;
        }

        .form-container input[type="radio"],
        .form-container select {
            padding: 8px;
            border: 1px solid #ccc;
            border-radius: 4px;
            box-sizing: border-box;
        }

        .form-container input[type="button"] {
            width: 100%;
            padding: 8px;
            border: none;
            border-radius: 4px;
            background-color: #5daf34;
            color: white;
            cursor: pointer;
        }

        .form-container input[type="submit"]:hover {
            background-color: #45a049;
        }


        .radio-group label {
            display: flex;
            align-items: center;
        }

        .group {
            display: flex;
            margin: 10px;
            flex-direction: row;
            align-items: center; /* 垂直居中对齐 */
            gap: 10px; /* 元素之间的间距 */
        }

        /* 如果希望label文本与select框对齐，可以添加如下样式 */
        label {
            white-space: nowrap; /* 防止文本换行 */
        }

        .custom-input {
            width: 150px;
            height: 35px;
            padding: 0 10px;
            border: 1px solid #ced4da;
            border-radius: 4px;
            font-size: 16px;
            color: #495057;
        }

        /* 模态对话框样式 */
        .modal {
            display: none; /* 默认隐藏 */
            position: fixed;
            z-index: 1001;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            overflow: auto;
            background-color: rgba(0, 0, 0, 0.4);
        }

        .modal-content {
            background-color: white;
            margin: 15% auto;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            max-width: 400px;

            text-align: center;
        }

        .modal-content h2 {
            margin-top: 0;
            color: #333;
        }

        .modal-content .input-group {
            margin-bottom: 15px;
            text-align: left;
        }

        .modal-content .input-group label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #555;
        }

        .modal-content .input-group input {
            width: 100%;
            padding: 10px;
            border: 1px solid #ccc;
            border-radius: 4px;
            box-sizing: border-box;
        }

        .modal-content button {
            width: 100%;
            padding: 10px;
            border: none;
            border-radius: 4px;
            background-color: #5daf34;
            color: white;
            cursor: pointer;
            font-size: 16px;
        }

        .modal-content button:hover {
            background-color: #45a049;
        }

        .close {
            position: absolute;
            top: 10px;
            right: 10px;
            font-size: 28px;
            font-weight: bold;
            color: #aaa;
            cursor: pointer;
        }

        .close:hover,
        .close:focus {
            color: black;
            text-decoration: none;
        }

    </style>
</head>
<body>
<div class="toolbar">
    <button class="connect" id="connect">连接</button>
    <button class="disconnect" id="disconnect">断开</button>
    <button class="disconnect" id="clear">清空</button>
    <button class="disconnect" id="openLoginBtn">登录</button>

</div>
<div class="form-container">
    <div id="status" class="disconnect" style="text-align: right">
        <span style="color: red;font-weight: bold;">状态：未连接</span>
    </div>
    <div>
        <form id="form">
            <div class="group">
                <label>环境:</label>
                <label><input type="radio" name="env" value="test" checked> test </label>
                <label><input type="radio" name="env" value="release"> release </label>
            </div>

            <!-- 下拉选择按钮 -->
            <div class="group">
                <label for="project">项目:</label>
                <select id="project" name="project" class="custom-select" style="width: 54%">
                    <option value="admin">admin</option>
                    <option value="enterprise">enterprise</option>
                    <option value="server">server</option>
                </select>
            </div>
            <div class="group">
                <label>分支</label>
                <input name="branch" type="text" value="dev" class="custom-input"/>
            </div>

            <div class="group custom-select" id="item-select" style="display: flex;flex-wrap: wrap; ">

            </div>
            <div class="group">
                <label>重启:</label>
                <label><input type="radio" name="restart" value="true" checked> true </label>
                <label><input type="radio" name="restart" value="false"> false </label>
            </div>

            <!-- 提交按钮 -->
            <p>
                <input type="button" value="提交" id="onsubmit">
            </p>
        </form>
    </div>
    <div id="messages"></div>

    <!-- 模态对话框 -->
    <div id="loginModal" class="modal">
        <div class="modal-content">
            <span class="close">&times;</span>
            <h2>登录</h2>
            <form id="loginForm">
                <div class="input-group">
                    <label for="username">Username:</label>
                    <input type="text" id="username" name="username" required>
                </div>
                <div class="input-group">
                    <label for="password">Password:</label>
                    <input type="password" id="password" name="password" required>
                </div>
                <button type="submit">登录</button>
            </form>
        </div>
    </div>
</div>


<script>
    let ws = null;
    let heartbeatInterval = null;

    function sendHeartbeat() {
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send('ping');
        }
    }

    function startHeartbeat(interval) {
        heartbeatInterval = setInterval(sendHeartbeat, interval || 10000); // 默认间隔10秒
    }

    function stopHeartbeat() {
        clearInterval(heartbeatInterval);
    }


    document.getElementById('connect').addEventListener('click', function () {
        if (!window.WebSocket) {
            alert('您的浏览器不支持 WebSocket');
            return;
        }

        if (ws && ws.readyState === WebSocket.OPEN) {
            alert('已经连接到WebSocket服务器');
            return;
        }

        let token = localStorage.getItem("token");
        if (token === "" || token == null) {
            alert("请先登录")
            return;
        }

        ws = new WebSocket('ws://` + url + `/webSocket?' + token);

        ws.onopen = function () {
            document.getElementById('status').innerHTML = '<span style="color: #5daf34;font-weight: bold;">状态：连接成功</span>';
 			document.getElementById('messages').innerHTML = '';
            startHeartbeat(); // 开始发送心跳
            sendHeartbeat(); // 连接成功后立即发送一次心跳
        };

        ws.onmessage = function (event) {
            if (event.data === "pong") {
                return
            }
            if (event.data === "no login") {
                document.getElementById('messages').innerHTML = '<span style="color: red;font-weight: bold; margin: 10% auto">请先登录!</span>';
                ws.close(3008, "noLogin")
                return
            }
            document.getElementById('messages').innerHTML += event.data + '<br>';
        };

        ws.onerror = function (error) {
            console.log(error)
            alert('WebSocket 连接错误');
        };

        ws.onclose = function (event) {
            console.log(event.code);
            let msg = '';
            if (event.code === 3008) {
                msg = '请先登录';
            } else {
                msg = '连接关闭';
            }
            document.getElementById('status').innerHTML = '<span style="color: red;font-weight: bold;">状态：' + msg + '</span>';
            ws = null
            stopHeartbeat();
            alert(msg);
        };
    });

    document.getElementById('clear').addEventListener('click', function () {
        document.getElementById('messages').innerHTML = '';
    });

    document.getElementById('disconnect').addEventListener('click', function () {
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.close();
        } else {
            alert("未连接websocket");
        }
    });

    document.getElementById('onsubmit').addEventListener('click', function () {
        if (ws && ws.readyState === WebSocket.OPEN) {
			document.getElementById('messages').innerHTML = '';
            if (document.querySelector('input[name="env"]:checked').value === "release") {
                let userChoice = confirm("你确定要继续吗？");
                if (!userChoice) {
                    return
                }
            }

            const form = document.getElementById('form');
            const formData = {};
            for (const element of form.elements) {
                if (element.name.length === 0) {
                    continue
                }
                // 多选
                if (element.type === 'checkbox') {
                    if (element.checked) {
                        if (!Array.isArray(formData[element.name])) {
                            formData[element.name] = [];
                        }
                        formData[element.name].push(element.value);
                    }
                    continue
                }
                // 单选框
                if (element.type === 'radio') {
                    if (element.checked) {
                        formData[element.name] = element.value;
                    }
                    continue
                }
                formData[element.name] = element.value;
            }
            const jsonData = JSON.stringify(formData);
            ws.send(jsonData)
        } else {
            alert("未连接websocket");
        }
    });

    document.getElementById('project').addEventListener('change', function () {
        let items = [];
        switch (this.value) {
            case "enterprise":
                items = ["soga_api_chat", "soga_api_chatroom", "soga_rpc_chat", "soga_rpc_game", "soga_cron", "soga_tool"];
                document.getElementById('item-select').innerHTML = '<label>工程:</label>';
                break
            case "server":
                items = ["soga_im_api", "soga_im_msg_gateway", "soga_im_msg_transfer", "soga_im_push",
                    "soga_im_rpc_auth", "soga_im_rpc_cache", "soga_im_rpc_conversation", "soga_im_rpc_friend",
                    "soga_im_rpc_group", "soga_im_rpc_msg", "soga_im_rpc_office", "soga_im_rpc_organization",
                    "soga_im_rpc_user"];
                document.getElementById('item-select').innerHTML = '<label>工程:</label>';
                break
            default:
                document.getElementById('item-select').innerHTML = '';
        }
        for (let i = 0; i < items.length; i++) {
            document.getElementById('item-select').innerHTML += '<label><input type="checkbox" name="items" value="' + items[i] + '">' + items[i] + '</label>'
        }
    });

    // 弹出登录框
    document.getElementById('openLoginBtn').addEventListener('click', function () {
        document.getElementById('loginModal').style.display = 'block';
    });

    document.querySelector('.close').addEventListener('click', function () {
        document.getElementById('loginModal').style.display = 'none';
    });

    window.onclick = function (event) {
        let modal = document.getElementById('loginModal');
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    };

    document.getElementById('loginForm').addEventListener('submit', function (event) {
        event.preventDefault(); // 阻止表单默认提交行为
        let data = JSON.stringify({
            username: document.getElementById('username').value,
            password: document.getElementById('password').value,
        });

        const xhr = new XMLHttpRequest();
        xhr.open('POST', "http://` + url + `/" + "login", true);
        xhr.setRequestHeader("content-type", "application/json;charset=UTF-8");
        xhr.onreadystatechange = function () {
            if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
                let result = JSON.parse(xhr.response);
                if (result.status === 200) {
                    localStorage.setItem("token", result.token);
                    document.getElementById('loginModal').style.display = 'none';
                }
                alert(result.message);
            }
        };
        xhr.send(data)


    });

</script>
</body>
</html>`
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
