package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"os"
)

func main() {
	// 设置克隆的远程仓库地址
	remoteURL := "http://192.168.0.13/platform/soga_admin.git"

	// 设置克隆到本地的路径
	clonePath := "/"
	f, _ := os.Open(clonePath)

	// 设置认证信息
	auth := &http.BasicAuth{
		Username: "user",
		Password: "password",
	}

	// 初始化内存存储（用于存储Git对象）
	storer := memory.NewStorage()

	// 克隆仓库
	_, err := git.Clone(storer, nil, &git.CloneOptions{
		URL:               remoteURL,
		Progress:          f,
		Auth:              auth,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		fmt.Println("Error cloning repo:", err)
		return
	}

	fmt.Println("Repository cloned successfully")
}
