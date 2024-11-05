package router

import (
	"deploy/api"
	"deploy/log"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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
	fmt.Println("ip:", ip)
	r.GET("/webSocket", (&api.WebSocket{}).WebSocketHandler)
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, "%s", html(ip))
	})
	_ = r.Run("0.0.0.0:8088")
}

func html(ip string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WebSocket 实时通信示例</title>
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


    </style>
</head>
<body>
<div class="toolbar">
    <button class="connect" id="connect">连接</button>
    <button class="disconnect" id="disconnect">断开</button>
    <button class="disconnect" id="clear">清空</button>

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

            <div class="group custom-select" id="item-select">

            </div>

            <!-- 提交按钮 -->
            <p>
                <input type="button" value="提交" id="onsubmit">
            </p>
        </form>
    </div>
    <div id="messages"></div>
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

        ws = new WebSocket('ws://` + ip + `:8088/webSocket');

        ws.onopen = function () {
            document.getElementById('status').innerHTML = '<span style="color: #5daf34;font-weight: bold;">状态：连接成功</span>';
            startHeartbeat(); // 开始发送心跳
            sendHeartbeat(); // 连接成功后立即发送一次心跳
        };

        ws.onmessage = function (event) {
            if (event.data === "pong") {
                console.log(event.data)
                return
            }
            document.getElementById('messages').innerHTML += event.data + '<br>';
        };

        ws.onerror = function (error) {
            console.log(error)
            alert('WebSocket 连接错误');
        };

        ws.onclose = function () {
            document.getElementById('status').innerHTML = '<span style="color: red;font-weight: bold;">状态：连接关闭</span>';
            ws = null;
            stopHeartbeat();
            alert("连接关闭")
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
                items = ["api-chat", "api-chatroom", "rpc-chat", "rpc-game", "cron"];
                document.getElementById('item-select').innerHTML = '<label>工程:</label><label><input type="checkbox" name="items" checked value="all">all</label>';
                break
            case "server":
                document.getElementById('item-select').innerHTML = '<label>工程:</label><label><input type="checkbox" name="items" checked value="all">all</label>';
                break
            default:
                document.getElementById('item-select').innerHTML = '';
        }
        for (let i = 0; i < items.length; i++) {
            document.getElementById('item-select').innerHTML += '<label><input type="checkbox" name="items" value="' + items[i] + '">' + items[i] + '</label>'
        }
    })

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
