package internal

import (
	"archive/zip"
	"bytes"
	"deploy/config"
	dlog "deploy/log"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gorilla/websocket"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DeployService struct {
}
type ClientConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

var (
	d, _ = os.Getwd()
	dir  = strings.ReplaceAll(d, "\\", "/") + "/"
	mu   sync.Mutex
)

type Message struct {
	Env     string   `json:"env"`
	Project string   `json:"project"`
	Branch  string   `json:"branch"`
	Restart bool     `json:"restart"`
	Items   []string `json:"items"`
}

type newWriter struct {
	Wr io.Writer
}

func (w *newWriter) Write(p []byte) (n int, err error) {
	if len(p) > 0 {
		p = bytes.ReplaceAll(p, []byte("\r"), []byte(""))
		p = bytes.ReplaceAll(p, []byte(" "), []byte(""))
		if len(p) > 0 {
			p = append(p, []byte("<br>")...)
			return w.Wr.Write(p)
		}

	}
	return
}

func (d *DeployService) AdminTest(conn *websocket.Conn, msg Message) {
	adminConf := config.Config.AdminTest
	err := os.Chdir(dir)
	if err != nil {
		d.Flush("Chdir err"+err.Error(), conn)
	}
	gitLog, err := d.Git(adminConf, msg.Branch, conn)
	if err != nil {
		return
	}

	if err = d.Build(adminConf, gitLog, conn); err != nil {
		return
	}

	if err = d.ZipFiles(adminConf.ProjectPath, adminConf.ZipFilePath, []string{adminConf.BuildConfigs[0].BinName}, conn); err != nil {
		return
	}

	_ = d.ScpUpload(adminConf, adminConf.BuildConfigs[0].Name, "pm2 restart soga_admin", msg.Restart, conn)
	return
}

func (d *DeployService) AdminRelease(conn *websocket.Conn, msg Message) {
	adminConf := config.Config.AdminRelease
	if err := os.Chdir(dir); err != nil {
		d.Flush("Chdir err"+err.Error(), conn)
	}

	gitLog, err := d.Git(adminConf, msg.Branch, conn)
	if err != nil {
		return
	}

	if err = d.Build(adminConf, gitLog, conn); err != nil {
		return
	}

	if err = d.ZipFiles(adminConf.ProjectPath, adminConf.ZipFilePath, []string{adminConf.BuildConfigs[0].BinName}, conn); err != nil {
		return
	}

	_ = d.ScpUpload(adminConf, adminConf.BuildConfigs[0].Name, "supervisorctl restart soga_admin", msg.Restart, conn)
	return
}

func (d *DeployService) EnterpriseTest(conn *websocket.Conn, msg Message) {
	cfg := config.Config.EnterpriseTest

	gitLog, err := d.Git(cfg, msg.Branch, conn)
	if err != nil {
		return
	}

	files := make([]string, 0)
	fileNames := make([]string, 0)
	binNames := make([]string, 0)

	newBuildConfig := make([]config.BuildConfig, 0)
	for _, item := range msg.Items {
		for _, bcfg := range cfg.BuildConfigs {
			if item == bcfg.Env {
				newBuildConfig = append(newBuildConfig, bcfg)
			}
		}
	}
	cfg.BuildConfigs = newBuildConfig
	if len(cfg.BuildConfigs) == 0 {
		d.Flush("没有可打包的💔💔💔", conn)
		return
	}

	for _, bcfg := range cfg.BuildConfigs {
		files = append(files, bcfg.BinName)
		fileNames = append(fileNames, bcfg.Name)
		switch bcfg.Name {
		case "soga_tool":
			continue
		case "soga_rpc_chat":
			binNames = append(binNames, "soga_api_rpc_chat")
		case "soga_rpc_game":
			binNames = append(binNames, "soga_api_rpc_game")
		default:
			binNames = append(binNames, bcfg.Name)
		}
	}

	if err = d.Build(cfg, gitLog, conn); err != nil {
		return
	}

	if err = d.ZipFiles(cfg.ProjectPath, cfg.ZipFilePath, files, conn); err != nil {
		return
	}

	_ = d.ScpUpload(cfg, strings.Join(fileNames, " "), "pm2 restart "+strings.Join(binNames, " "), msg.Restart, conn)

	return
}

func (d *DeployService) EnterpriseRelease(conn *websocket.Conn, msg Message) {
	cfg := config.Config.EnterpriseRelease

	gitLog, err := d.Git(cfg, msg.Branch, conn)
	if err != nil {
		return
	}

	files := make([]string, 0)
	fileNames := make([]string, 0)
	binNames := make([]string, 0)

	newBuildConfig := make([]config.BuildConfig, 0)
	for _, item := range msg.Items {
		for _, bcfg := range cfg.BuildConfigs {
			if item == bcfg.Env {
				newBuildConfig = append(newBuildConfig, bcfg)
			}
		}
	}
	cfg.BuildConfigs = newBuildConfig
	if len(cfg.BuildConfigs) == 0 {
		d.Flush("没有可打包的💔💔💔", conn)
		return
	}

	restartCmd := ""
	for _, bcfg := range cfg.BuildConfigs {
		files = append(files, bcfg.BinName)
		fileNames = append(fileNames, bcfg.Name)
		switch bcfg.Name {
		case "soga_tool":
			continue
		case "soga_rpc_chat":
			binNames = append(binNames, "soga_api_rpc_chat")
		case "soga_rpc_game":
			binNames = append(binNames, "soga_api_rpc_game")
		case "soga_cron":
			restartCmd = "mv /root/soga_im_enterprise/bin/soga_cron /root/soga_im_cron/ && mv /root/soga_im_cron/soga_cron /root/soga_im_cron/soga_im_cron && "
		default:
			binNames = append(binNames, bcfg.Name)
		}
	}

	if err = d.Build(cfg, gitLog, conn); err != nil {
		return
	}

	if err = d.ZipFiles(cfg.ProjectPath, cfg.ZipFilePath, files, conn); err != nil {
		return
	}

	_ = d.ScpUpload(cfg, strings.Join(fileNames, " "), restartCmd+"pm2 restart "+strings.Join(binNames, " "), msg.Restart, conn)

	return
}

func (d *DeployService) ServerTest(conn *websocket.Conn, msg Message) {
	cfg := config.Config.ServerTest

	files := make([]string, 0)
	fileNames := make([]string, 0)

	newBuildConfig := make([]config.BuildConfig, 0)
	for _, item := range msg.Items {
		for _, bcfg := range cfg.BuildConfigs {
			if item == bcfg.Env {
				newBuildConfig = append(newBuildConfig, bcfg)
			}
		}
	}
	cfg.BuildConfigs = newBuildConfig

	for _, bcfg := range cfg.BuildConfigs {
		files = append(files, bcfg.BinName)
		fileNames = append(fileNames, bcfg.Name)
	}

	if len(cfg.BuildConfigs) == 0 {
		d.Flush("没有可打包的💔💔💔", conn)
		return
	}

	gitLog, err := d.Git(cfg, msg.Branch, conn)
	if err != nil {
		return
	}

	if err = d.Build(cfg, gitLog, conn); err != nil {
		return
	}

	if err = d.ZipFiles(cfg.ProjectPath, cfg.ZipFilePath, files, conn); err != nil {
		return
	}

	_ = d.ScpUpload(cfg, strings.Join(fileNames, " "), "pm2 restart "+strings.Join(fileNames, " "), msg.Restart, conn)

	return
}

func (d *DeployService) ServerRelease(conn *websocket.Conn, msg Message) {
	cfg := config.Config.ServerRelease

	files := make([]string, 0)
	fileNames := make([]string, 0)

	newBuildConfig := make([]config.BuildConfig, 0)
	for _, item := range msg.Items {
		for _, bcfg := range cfg.BuildConfigs {
			if item == bcfg.Env {
				newBuildConfig = append(newBuildConfig, bcfg)
			}
		}
	}
	cfg.BuildConfigs = newBuildConfig

	for _, bcfg := range cfg.BuildConfigs {
		files = append(files, bcfg.BinName)
		fileNames = append(fileNames, bcfg.Name)
	}

	if len(cfg.BuildConfigs) == 0 {
		d.Flush("没有可打包的💔💔💔", conn)
		return
	}

	gitLog, err := d.Git(cfg, msg.Branch, conn)
	if err != nil {
		return
	}

	if err = d.Build(cfg, gitLog, conn); err != nil {
		return
	}

	if err = d.ZipFiles(cfg.ProjectPath, cfg.ZipFilePath, files, conn); err != nil {
		return
	}

	_ = d.ScpUpload(cfg, strings.Join(fileNames, " "), "supervisorctl restart "+strings.Join(fileNames, " "), msg.Restart, conn)

	return
}

func (d *DeployService) Flush(msg string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	_ = conn.SetReadDeadline(time.Now().Add(time.Minute))
	_ = conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (d *DeployService) gitPull(worktree *git.Worktree, auth *http.BasicAuth, try int) (err error) {
	err = worktree.Pull(&git.PullOptions{
		Auth:  auth,
		Force: true,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		if try > 3 {
			return err
		}
		return d.gitPull(worktree, auth, try+1)
	}
	return nil
}

func (d *DeployService) gitCheckout(wo *git.Worktree, branch string, try int) (err error) {
	err = wo.Checkout(&git.CheckoutOptions{
		Force:  true,
		Branch: plumbing.NewRemoteReferenceName(git.DefaultRemoteName, branch),
	})

	if errors.Is(err, git.ErrUnstagedChanges) {
		if try > 3 {
			return err
		}
		return d.gitCheckout(wo, branch, try+1)
	} else {
		return
	}
}

func (d *DeployService) GitLog(depth int) (str string, err error) {
	gLog, err := exec.Command("git", "log", "-n", strconv.Itoa(depth), "--format='%h- %an, %ar : %s'").Output()
	if err != nil {
		return
	}
	if strings.Index(string(gLog), "Merge remote-tracking branch") == -1 {
		return string(gLog), nil
	} else {
		if depth == 1 {
			return d.GitLog(depth + 1)
		}
		return string(gLog), nil
	}

}

// Git 拉取代码
func (d *DeployService) Git(cfg config.Configure, branch string, conn *websocket.Conn) (log string, err error) {
	d.Flush("git 开始拉取... 🚀🚀🚀", conn)
	defer func() {
		_ = os.Chdir(dir)
		if err != nil {
			d.Flush("git 错误 💔💔💔"+err.Error(), conn)
			dlog.Error(err, string(debug.Stack()))
		} else {
			dlog.Info("git Success 👌👌👌")
			d.Flush("git Success 👌👌👌", conn)
		}
	}()
	if _, err = os.Stat(cfg.ProjectPath); err != nil {
		if err = os.MkdirAll(cfg.ProjectPath, fs.ModePerm); err != nil {
			return
		}
	}

	_ = os.Chdir(cfg.ProjectPath)
	mu.Lock()
	defer mu.Unlock()
	mw, _ := conn.NextWriter(websocket.TextMessage)

	gitCmdFun := func(w io.Writer, arg ...string) (err error) {
		cmd := exec.Command("git", arg...)
		cmd.Stdout = w
		cmd.Stderr = w
		err = cmd.Run()
		return
	}

	if _, err = os.Stat(cfg.ProjectName); err != nil {
		if err = gitCmdFun(mw, "clone", cfg.GitUrl); err != nil {
			return
		}
	}

	if err = os.Chdir(cfg.ProjectName); err != nil {
		return
	}

	if err = gitCmdFun(mw, "fetch", "origin"); err != nil {
		return
	}

	if err = gitCmdFun(mw, "checkout", branch); err != nil {
		return
	}

	if err = gitCmdFun(mw, "pull", "origin", branch, "--force"); err != nil {
		return
	}

	log, err = d.GitLog(1)
	if err != nil {
		return
	}

	if err = gitCmdFun(mw, "submodule", "update", "--init", "--recursive"); err != nil {
		return
	}

	_ = os.Chdir("depends")

	if err = gitCmdFun(mw, "fetch", "origin"); err != nil {
		return
	}

	if err = gitCmdFun(mw, "checkout", branch); err != nil {
		return
	}

	if err = gitCmdFun(mw, "pull", "origin", branch, "--force"); err != nil {
		return
	}

	dLog, err := d.GitLog(1)
	if err != nil {
		return
	}

	return log + "【depends】:" + dLog, nil
}

// Build 更新
func (d *DeployService) Build(cfg config.Configure, gitLog string, conn *websocket.Conn) (err error) {
	d.Flush("开始打包...🚀🚀🚀 ", conn)
	var version string
	defer func() {
		_ = os.Chdir(dir)
		if err != nil {
			dlog.Error(err, string(debug.Stack()))
			d.Flush("打包错误 💔💔💔"+err.Error(), conn)
		} else {
			dlog.Info("打包版本【" + version + "】 Success 💯💯💯")
			d.Flush("打包版本【"+version+"】 Success 💯💯💯", conn)
		}
	}()
	// 存放bin的目录
	_, err = os.Stat(cfg.BinPath)
	if !os.IsExist(err) {
		// 目录不存在，则创建它
		err = os.MkdirAll(cfg.BinPath, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	// 设置编译环境
	_ = os.Setenv("CGO_ENABLED", "0")
	_ = os.Setenv("GOOS", "linux")

	// 版本信息
	version = "v" + time.Now().Format("20060102150405")
	gitLog = strings.ReplaceAll(gitLog, "\n", ";")
	gitLog = strings.ReplaceAll(gitLog, "\t", "-")
	gitLog = strings.ReplaceAll(gitLog, " ", ",")
	ldflags := fmt.Sprintf(`-ldflags=-X main.version=%s -X "main.gitInfo=%s"`, version, gitLog)

	for _, build := range cfg.BuildConfigs {
		if err = os.Chdir(build.ModPath); err != nil {
			return err
		}
		d.Flush("【"+build.Name+"】 go mod tidy start...", conn)
		w, _ := conn.NextWriter(websocket.TextMessage)
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Stdout = io.MultiWriter(os.Stdout, w)
		cmd.Stderr = io.MultiWriter(os.Stdout, w)
		err := cmd.Run()
		if err != nil {
			return err
		}

		d.Flush("【"+build.Name+"】go mod tidy finished...", conn)

		d.Flush("go build 【"+build.Name+"】 start... 🚀🚀🚀", conn)
		buildOut, err := exec.Command("go", "build", "-o", dir+cfg.ProjectPath+"/"+build.BinName, "-gcflags=all=-N -l", ldflags, "-trimpath").CombinedOutput()
		if len(buildOut) > 0 {
			d.Flush(string(buildOut), conn)
		}
		if err != nil {
			return err
		}

		d.Flush("go build 【"+build.Name+"】 success... 👌👌👌", conn)
		_ = os.Chdir(dir)
	}

	d.Flush("go build all finished...👍👍👍", conn)
	return
}

func (d *DeployService) ZipFiles(projectPath, zipFilePath string, files []string, conn *websocket.Conn) (err error) {
	_ = os.Chdir(projectPath)
	d.Flush("开始删除压缩文件"+zipFilePath+"...🚀🚀🚀", conn)
	// 删除压缩文件
	if _, err = os.Stat(zipFilePath); err == nil {
		if err = os.Remove(zipFilePath); err != nil {
			return err
		}
	}
	d.Flush("删除压缩文件成功"+zipFilePath+"...✔️✔️✔️", conn)

	defer func() {
		_ = os.Chdir(dir)
		if err != nil {
			dlog.Error(err, string(debug.Stack()))
			d.Flush("压缩 错误 💔💔💔"+err.Error(), conn)
		} else {
			dlog.Info("压缩 Success 👌👌👌")
			d.Flush("压缩 Success 👌👌👌", conn)
		}
	}()
	d.Flush("开始压缩...🚀🚀🚀", conn)
	// 创建 ZIP 文件
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("无法创建 ZIP 文件: %w", err)
	}
	defer zipFile.Close()

	// 创建 ZIP 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		d.Flush("开始压缩文件"+file+"...🚀🚀🚀", conn)
		// 打开要压缩的文件
		fileToZip, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("无法打开文件 %s: %w", file, err)
		}
		defer fileToZip.Close()

		// 获取文件信息
		info, err := fileToZip.Stat()
		if err != nil {
			return fmt.Errorf("无法获取文件信息 %s: %w", file, err)
		}

		// 创建 ZIP 文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("无法创建文件头 %s: %w", file, err)
		}

		// 将文件名相对化
		header.Name = filepath.Base(file)
		header.Method = zip.Deflate // 设置压缩方法

		// 创建写入器
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("无法创建写入器 %s: %w", file, err)
		}

		// 将文件内容写入 ZIP
		if _, err := io.Copy(writer, fileToZip); err != nil {
			return fmt.Errorf("写入文件 %s 到 ZIP 失败: %w", file, err)
		}
		d.Flush("压缩文件"+file+"...👌👌👌", conn)
	}

	return nil
}

func (d *DeployService) ScpUpload(conf config.Configure, binName, restartCmd string, restart bool, conn *websocket.Conn) (err error) {
	_ = os.Chdir(conf.ProjectPath)
	d.Flush("开始远程服务器 "+conf.Host+" 执行...🚀🚀🚀 ", conn)
	defer func() {
		_ = os.Chdir(dir)
		if err != nil {
			dlog.Error(err, string(debug.Stack()))
			d.Flush("服务器执行失败 💔💔💔"+err.Error(), conn)
		} else {
			dlog.Info("服务器执行 Success 💯💯💯")
			d.Flush("服务器执行 Success 💯💯💯", conn)
		}
	}()
	// SSH 配置
	c := &ssh.ClientConfig{
		User: conf.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 建立 SSH 连接
	addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	client, err := ssh.Dial("tcp", addr, c)
	if err != nil {
		return fmt.Errorf("无法连接到服务器: %w", err)
	}
	defer client.Close()

	// 打开本地zip文件
	localFile, err := os.Open(conf.ZipFilePath)
	if err != nil {
		return fmt.Errorf("无法打开本地文件: %w", err)
	}
	defer localFile.Close()

	// 创建 SSH 会话
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("无法创建 SSH 会话: %w", err)
	}
	defer session.Close()

	// 使用 SCP 传输文件命令
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("无法获取标准输入管道: %w", err)
	}

	// 启动 SCP 命令来接收文件
	if err := session.Start(fmt.Sprintf("scp -qt %s", conf.ServerPath)); err != nil {
		return fmt.Errorf("无法启动会话: %w", err)
	}

	d.Flush("开始上传 "+conf.Host+" 🚀🚀🚀 ", conn)

	// 文件传输前，必须要向远程服务器发送文件头信息，包括文件大小和权限
	fileInfo, err := localFile.Stat()
	if err != nil {
		return fmt.Errorf("无法获取本地文件信息: %w", err)
	}
	fileSize := fileInfo.Size()
	fileName := fileInfo.Name()
	_, _ = fmt.Fprintf(stdin, "C0644 %d %s\n", fileSize, fileName)

	mu.Lock()
	writer, _ := conn.NextWriter(websocket.TextMessage)

	newWriter := &newWriter{Wr: writer}

	bar := progressbar.NewOptions64(
		fileSize,
		progressbar.OptionSetDescription("uploading:"),
		progressbar.OptionSetWidth(10),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWriter(newWriter),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: "-",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	bar.StartWithoutRender()

	// 包装文件读取器，监控进度
	_, _ = fmt.Fprintf(writer, "%s", bar.String())
	progressReader := io.TeeReader(localFile, bar)

	// 复制文件内容到远程服务器
	if _, err := io.Copy(stdin, progressReader); err != nil {
		return fmt.Errorf("文件传输失败: %w", err)
	}

	// 结束文件传输
	_, _ = fmt.Fprint(stdin, "\x00")
	_ = stdin.Close()
	mu.Unlock()

	d.Flush("上传完成 "+conf.Host+" ✌️✌️✌️ ", conn)

	// 等待会话结束
	if err := session.Wait(); err != nil {
		return fmt.Errorf("文件传输会话执行失败: %w", err)
	}
	d.Flush("<br>文件上传成功...✔️✔️✔️", conn)

	d.Flush("服务器开始解压...🚀🚀🚀", conn)

	// 解压
	unsession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("无法创建 SSH 会话: %w", err)
	}
	defer unsession.Close()

	uncmd := fmt.Sprintf("cd %s && unzip -o %s && chmod +x %s", conf.ServerPath, conf.ZipName, binName)
	un, err := unsession.Output(uncmd)
	if err != nil {
		return fmt.Errorf("解压会话执行失败 ssh: command %v failed", err)
	}
	d.Flush(string(un), conn)
	d.Flush("服务器解压成功...✔️✔️✔️", conn)

	// 需要重启
	if restart {
		// 重启
		d.Flush("服务器开始重启...🚀🚀🚀", conn)
		resession, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("无法创建 SSH 会话: %w", err)
		}
		defer resession.Close()

		re, err := resession.Output(restartCmd)
		if err != nil {
			return fmt.Errorf("重启会话执行失败 ssh: command %v failed", err)
		}
		d.Flush(string(re), conn)
		d.Flush("服务器重启成功...✔️✔️✔️", conn)
	}

	return nil
}
