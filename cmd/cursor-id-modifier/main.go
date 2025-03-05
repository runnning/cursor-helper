package main

import (
	"github.com/yuaotian/go-cursor-help/internal/config"
	"github.com/yuaotian/go-cursor-help/internal/process"
	"github.com/yuaotian/go-cursor-help/internal/ui"
	"github.com/yuaotian/go-cursor-help/pkg/idgen"
	"github.com/yuaotian/go-cursor-help/pkg/system"
)

func main() {
	// 检查是否已经是管理员权限
	isAdmin, _ := system.CheckAdminPrivileges()

	// 初始化组件
	username := system.GetCurrentUser()
	configManager, _ := config.NewManager(username)
	generator := idgen.NewGenerator()
	processManager := process.NewManager(nil, nil)

	// 创建并运行UI
	display := ui.NewFyneDisplay(configManager, generator, processManager, isAdmin)
	display.Initialize()
	display.Run()
}
