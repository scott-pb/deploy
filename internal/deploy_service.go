package internal

import (
	"archive/zip"
	"deploy/config"
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
	"strings"
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
	d, _    = os.Getwd()
	dir     = strings.ReplaceAll(d, "\\", "/") + "/"
	gitAuth = &http.BasicAuth{
		Username: "bin.peng",
		Password: "pb123654",
	}
	testConf = &ClientConfig{
		Host:     "192.168.0.30",
		Port:     "22",
		User:     "root",
		Password: "admin@123",
	}
	releaseConf = &ClientConfig{
		Host:     "8.212.40.18",
		Port:     "22",
		User:     "root",
		Password: "X@135791667",
	}
)

var (
	testAdminClientPath    = "/root/soga/soga_im_admin_server"
	releaseAdminClientPath = "/root/soga_admin_server"
)

var (
	enterpriseProjectPath    = dir + "/project/soga_im_enterprise"
	enterpriseBinPath        = enterpriseProjectPath + "/bin"
	enterpriseZipName        = "im_enterprise.zip"
	enterpriseZip            = enterpriseBinPath + "/im_enterprise.zip"
	enterpriseGitUrl         = "http://192.168.0.13/server/soga_im_enterprise"
	testEnterpriseClientPath = "/root/soga/soga_im_enterprise/bin"
	nameMap                  = map[string]string{
		"api-chat":     enterpriseBinPath + "/soga_api_chat",
		"api-chatroom": enterpriseBinPath + "/soga_api_chatroom",
		"rpc-chat":     enterpriseBinPath + "/soga_rpc_chat",
		"rpc-game":     enterpriseBinPath + "/soga_rpc_game",
		"cron":         enterpriseBinPath + "/soga_cron",
	}
	enterpriseBuildPathMap = map[string]string{
		"api-chat":     enterpriseProjectPath + "/cmd/api/chat",
		"api-chatroom": enterpriseProjectPath + "/cmd/api/chatroom",
		"rpc-chat":     enterpriseProjectPath + "/cmd/rpc/chat",
		"rpc-game":     enterpriseProjectPath + "/cmd/rpc/game",
		"cron":         enterpriseProjectPath + "/cmd/cron",
	}

	enterpriseBinNameMap = map[string]string{
		"api-chat":     "soga_api_chat",
		"api-chatroom": "soga_api_chatroom",
		"rpc-chat":     "soga_rpc_chat",
		"rpc-game":     "soga_rpc_game",
		"cron":         "soga_cron",
	}

	enterpriseRestartMap = map[string]string{
		"api-chat":     "soga_api_chat",
		"api-chatroom": "soga_api_chatroom",
		"rpc-chat":     "soga_api_rpc_chat",
		"rpc-game":     "soga_api_rpc_game",
		"cron":         "soga_cron",
	}
)

type Message struct {
	Env     string   `json:"env"`
	Project string   `json:"project"`
	Branch  string   `json:"branch"`
	Items   []string `json:"items"`
}

func (d *DeployService) AdminTest(conn *websocket.Conn, branch string) {
	adminConf := config.Config.AdminTest
	auth := &http.BasicAuth{
		Username: adminConf.GitConfig.UserName,
		Password: adminConf.GitConfig.PassWord,
	}
	gitLog, err := d.Git(adminConf.ProjectPath, adminConf.GitUrl, branch, auth, conn)
	if err != nil {
		return
	}

	m := map[string]string{
		adminConf.BinPath: adminConf.ProjectPath,
	}
	if err = d.Build(m, adminConf.BinPath, gitLog, conn); err != nil {
		return
	}
	//if err = d.ZipFiles(adminZip, []string{adminBuildBinName}, conn); err != nil {
	//	return
	//}
	//_ = d.ScpUpload(testConf, adminZip, testAdminClientPath, adminZipName, adminBinName, "pm2 restart soga_admin", conn)
	return
}

//func (d *DeployService) AdminRelease(conn *websocket.Conn, branch string) {
//	gitLog, err := d.Git(adminProjectPath, adminGitUrl, branch, conn)
//	if err != nil {
//		return
//	}
//	m := map[string]string{
//		adminBuildBinName: adminProjectPath,
//	}
//	if err = d.Build(m, adminBinPath, gitLog, conn); err != nil {
//		return
//	}
//	if err = d.ZipFiles(adminZip, []string{adminBuildBinName}, conn); err != nil {
//		return
//	}
//	_ = d.ScpUpload(releaseConf, adminZip, releaseAdminClientPath, adminZipName, adminBinName, "supervisorctl restart soga_admin", conn)
//	return
//}

//func (d *DeployService) EnterpriseTest(conn *websocket.Conn, msg Message) {
//	gitLog, err := d.Git(enterpriseProjectPath, enterpriseGitUrl, msg.Branch, conn)
//	if err != nil {
//		return
//	}
//
//	m := make(map[string]string)
//	files := make([]string, 0)
//	binName := make([]string, 0)
//	restartName := make([]string, 0)
//	if slices.Contains(msg.Items, "all") {
//		for k, v := range nameMap {
//			files = append(files, v)
//			binName = append(binName, enterpriseBinNameMap[k])
//			restartName = append(restartName, enterpriseRestartMap[k])
//			m[v] = enterpriseBuildPathMap[k]
//		}
//	} else {
//		for _, k := range msg.Items {
//			binName = append(binName, enterpriseBinNameMap[k])
//			restartName = append(restartName, enterpriseRestartMap[k])
//			files = append(files, nameMap[k])
//			m[nameMap[k]] = enterpriseBuildPathMap[k]
//		}
//	}
//
//	if err = d.Build(m, adminBinPath, gitLog, conn); err != nil {
//		return
//	}
//
//	if err = d.ZipFiles(enterpriseZip, files, conn); err != nil {
//		return
//	}
//	_ = d.ScpUpload(testConf, enterpriseZip, testEnterpriseClientPath, enterpriseZipName, strings.Join(binName, " "), "pm2 restart "+strings.Join(restartName, " "), conn)
//	return
//}

func flush(msg string, conn *websocket.Conn) {
	_ = conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

// Git 拉取代码
func (d *DeployService) Git(dirName, remoteURL, branch string, gitAuth *http.BasicAuth, conn *websocket.Conn) (log string, err error) {
	flush("git 开始拉取... 🚀🚀🚀", conn)
	defer func() {
		if err != nil {
			flush("git 错误 💔💔💔"+err.Error(), conn)
		} else {
			flush("git Success 👌👌👌", conn)
		}
	}()
	if _, err := os.Stat(dirName); err != nil {
		_ = os.Mkdir(dirName, 0755)
	}

	var r *git.Repository
	r, err = git.PlainOpen(dirName)
	if err != nil {
		if !errors.Is(err, git.ErrRepositoryNotExists) {
			return
		} else {
			w, _ := conn.NextWriter(websocket.TextMessage)
			r, err = git.PlainClone(dirName, false, &git.CloneOptions{
				URL:      remoteURL,
				Progress: io.MultiWriter(os.Stdout, w),
				Auth:     gitAuth,
			})
			if err != nil {
				return
			}
		}
	}

	wo, err := r.Worktree()
	if err != nil {
		return
	}

	if err = wo.Checkout(&git.CheckoutOptions{
		Force:  true,
		Branch: plumbing.NewRemoteReferenceName(git.DefaultRemoteName, branch),
	}); err != nil {
		return
	}

	err = wo.Pull(&git.PullOptions{
		Auth:  gitAuth,
		Force: true,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return
	}
	err = wo.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewRemoteReferenceName(git.DefaultRemoteName, branch),
		Force:  true,
	})
	if err != nil && !errors.Is(err, git.ErrUnstagedChanges) {
		return
	}

	sd, err := wo.Submodule("depends")
	if err != nil {
		return
	}

	if err = sd.Init(); err != nil && !errors.Is(err, git.ErrSubmoduleAlreadyInitialized) {
		return
	}

	err = sd.Update(&git.SubmoduleUpdateOptions{
		Init: true,
		Auth: gitAuth,
	})
	if err != nil {
		return
	}

	dr, _ := sd.Repository()
	dw, _ := dr.Worktree()

	if err = dw.Checkout(&git.CheckoutOptions{
		Force:  true,
		Branch: plumbing.NewRemoteReferenceName(git.DefaultRemoteName, branch),
	}); err != nil {
		return
	}

	err = dw.Pull(&git.PullOptions{
		Auth:  gitAuth,
		Force: true,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return
	}
	err = dw.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewRemoteReferenceName(git.DefaultRemoteName, branch),
	})
	if err != nil && !errors.Is(err, git.ErrUnstagedChanges) {
		return
	}
	ml, err := r.Log(&git.LogOptions{})
	if err != nil {
		return
	}
	defer ml.Close()
	commit, _ := ml.Next()
	b := fmt.Sprintf("【log】:%s %s %s %s", commit.Author.When.Format(time.DateTime), commit.Author.Name, commit.Hash.String()[:8], commit.Message)
	if strings.Index(commit.Message, "Merge remote-tracking branch") != -1 {
		commit, _ = ml.Next()
		b += fmt.Sprintf(";%s %s %s %s", commit.Author.When.Format(time.DateTime), commit.Author.Name, commit.Hash.String()[:8], commit.Message)
	}
	flush(b, conn)

	sl, err := dr.Log(&git.LogOptions{})
	if err != nil {
		return
	}
	defer sl.Close()
	scommit, _ := sl.Next()
	sb := fmt.Sprintf("【depends-log】:%s %s %s %s", commit.Author.When.Format(time.DateTime), scommit.Author.Name, scommit.Hash.String()[:8], scommit.Message)
	if strings.Index(scommit.Message, "Merge remote-tracking branch") != -1 {
		scommit, _ = sl.Next()
		sb += fmt.Sprintf("【depends-log】:%s %s %s %s", commit.Author.When.Format(time.DateTime), scommit.Author.Name, scommit.Hash.String()[:8], scommit.Message)
	}
	if err != nil {
		return
	}
	flush(sb, conn)

	return b + sb, nil
}

// Build 更新
func (d *DeployService) Build(buildMaps map[string]string, binPath, gitLog string, conn *websocket.Conn) (err error) {
	flush("开始打包...🚀🚀🚀 ", conn)
	var version string
	defer func() {
		if err != nil {
			flush("打包错误 💔💔💔"+err.Error(), conn)
		} else {
			flush("打包版本【"+version+"】 Success 💯💯💯", conn)
		}
	}()
	// 存放bin的目录
	_, err = os.Stat(binPath)
	if !os.IsExist(err) {
		// 目录不存在，则创建它
		err = os.MkdirAll(binPath, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	// 设置编译环境
	_ = os.Setenv("CGO_ENABLED", "0")
	_ = os.Setenv("GOOS", "linux")

	// 版本信息
	version = "v" + time.Now().Format("20060102150405")
	ldflags := fmt.Sprintf(`-ldflags=-X main.version=%s -X "main.gitInfo=%s"`, version, gitLog)

	for name, buildDir := range buildMaps {
		if err = os.Chdir(buildDir); err != nil {
			return err
		}
		flush("【"+name+"】 go mod tidy start...", conn)
		tidy, err := exec.Command("go", "mod", "tidy").CombinedOutput()
		if err != nil {
			return err
		}
		if len(tidy) > 0 {
			flush(string(tidy), conn)
		}

		flush("【"+name+"】go mod tidy finished...", conn)

		flush("go build 【"+name+"】 start... 🚀🚀🚀", conn)
		build, err := exec.Command("go", "build", "-o", name, "-gcflags=all=-N -l", ldflags, "-trimpath").CombinedOutput()
		if len(build) > 0 {
			flush(string(build), conn)
		}
		if err != nil {
			return err
		}

		flush("go build 【"+name+"】 success... 👌👌👌", conn)
	}

	flush("go build all finished...👍👍👍", conn)
	return
}

func (d *DeployService) ZipFiles(zipFilePath string, files []string, conn *websocket.Conn) (err error) {
	flush("开始删除压缩文件"+zipFilePath+"...🚀🚀🚀", conn)
	// 删除压缩文件
	if _, err = os.Stat(zipFilePath); err == nil {
		if err = os.Remove(zipFilePath); err != nil {
			return err
		}
	}
	flush("删除压缩文件成功"+zipFilePath+"...✔️✔️✔️", conn)

	defer func() {
		if err != nil {
			flush("压缩 错误 💔💔💔"+err.Error(), conn)
		} else {
			flush("压缩 Success 👌👌👌", conn)
		}
	}()
	flush("开始压缩...🚀🚀🚀", conn)
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
	}

	return nil
}

func (d *DeployService) ScpUpload(conf *ClientConfig, uploadZip, remoteClientPath, zipName, binName, restartCmd string, conn *websocket.Conn) (err error) {
	flush("开始远程服务器 "+conf.Host+" 执行...🚀🚀🚀 ", conn)
	defer func() {
		if err != nil {
			flush("服务器执行失败 💔💔💔"+err.Error(), conn)
		} else {
			flush("服务器执行 Success 💯💯💯", conn)
		}
	}()
	// SSH 配置
	config := &ssh.ClientConfig{
		User: conf.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 建立 SSH 连接
	addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("无法连接到服务器: %w", err)
	}
	defer client.Close()

	// 打开本地zip文件
	localFile, err := os.Open(uploadZip)
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
	if err := session.Start(fmt.Sprintf("scp -qt %s", remoteClientPath)); err != nil {
		return fmt.Errorf("无法启动会话: %w", err)
	}

	// 文件传输前，必须要向远程服务器发送文件头信息，包括文件大小和权限
	fileInfo, err := localFile.Stat()
	if err != nil {
		return fmt.Errorf("无法获取本地文件信息: %w", err)
	}
	fileSize := fileInfo.Size()
	fileName := fileInfo.Name()
	_, _ = fmt.Fprintf(stdin, "C0644 %d %s\n", fileSize, fileName)

	writer, _ := conn.NextWriter(websocket.TextMessage)

	bar := progressbar.NewOptions64(
		fileSize,
		progressbar.OptionSetDescription("文件上传中"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWriter(writer),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: "<p>",
			BarStart:      "",
			BarEnd:        "",
		}),
	)

	bar.StartWithoutRender()

	// 包装文件读取器，监控进度
	_, _ = fmt.Fprintf(writer, "\r%s", bar.String())
	progressReader := io.TeeReader(localFile, bar)

	// 复制文件内容到远程服务器
	if _, err := io.Copy(stdin, progressReader); err != nil {
		return fmt.Errorf("文件传输失败: %w", err)
	}

	// 结束文件传输
	_, _ = fmt.Fprint(stdin, "\x00")
	_ = stdin.Close()

	// 等待会话结束
	if err := session.Wait(); err != nil {
		return fmt.Errorf("文件传输会话执行失败: %w", err)
	}
	flush("<br>文件上传成功...✔️✔️✔️", conn)

	flush("服务器开始解压...🚀🚀🚀", conn)

	// 解压
	unsession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("无法创建 SSH 会话: %w", err)
	}
	defer unsession.Close()

	uncmd := fmt.Sprintf("cd %s && unzip -o %s && chmod +x %s", remoteClientPath, zipName, binName)
	un, err := unsession.Output(uncmd)
	if err != nil {
		return fmt.Errorf("解压会话执行失败 ssh: command %v failed", err)
	}
	flush(string(un), conn)
	flush("服务器解压成功...✔️✔️✔️", conn)

	// 重启
	flush("服务器开始重启...🚀🚀🚀", conn)
	resession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("无法创建 SSH 会话: %w", err)
	}
	defer resession.Close()

	re, err := resession.Output(restartCmd)
	if err != nil {
		return fmt.Errorf("重启会话执行失败 ssh: command %v failed", err)
	}
	flush(string(re), conn)
	flush("服务器重启成功...✔️✔️✔️", conn)
	return nil
}
