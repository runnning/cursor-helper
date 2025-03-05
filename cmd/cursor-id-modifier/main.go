package main

import (
	"os/exec"
	"os/user"
	"runtime"

	"github.com/yuaotian/go-cursor-help/internal/config"
	"github.com/yuaotian/go-cursor-help/internal/process"
	"github.com/yuaotian/go-cursor-help/internal/ui"
	"github.com/yuaotian/go-cursor-help/pkg/idgen"
)

func main() {
	// 检查是否已经是管理员权限
	isAdmin, _ := checkAdminPrivileges()

	// 初始化组件
	username := getCurrentUser()
	configManager, _ := config.NewManager(username)
	generator := idgen.NewGenerator()
	processManager := process.NewManager(nil, nil)

	// 创建并运行UI
	display := ui.NewFyneDisplay(configManager, generator, processManager, isAdmin)
	display.Initialize()
	display.Run()
}

func getCurrentUser() string {
	if user, err := user.Current(); err == nil {
		return user.Username
	}
	return ""
}

func checkAdminPrivileges() (bool, error) {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("net", "session").Run() == nil, nil
	case "darwin", "linux":
		if user, err := user.Current(); err != nil {
			return false, err
		} else {
			return user.Uid == "0", nil
		}
	default:
		return false, nil
	}
}
