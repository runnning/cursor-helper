package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yuaotian/go-cursor-help/internal/config"
	"github.com/yuaotian/go-cursor-help/internal/process"
	"github.com/yuaotian/go-cursor-help/pkg/idgen"
)

// FyneDisplay 实现了基于Fyne的GUI界面
type FyneDisplay struct {
	app            fyne.App
	window         fyne.Window
	progress       *widget.ProgressBarInfinite
	status         *widget.Label
	configManager  *config.Manager
	generator      *idgen.Generator
	processManager *process.Manager
	isAdmin        bool
	logFile        *os.File
}

func (d *FyneDisplay) setupLogging() error {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0666); err != nil {
		return err
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logDir, fmt.Sprintf("cursor_modifier_%s.log", timestamp))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	d.logFile = logFile
	d.log("日志系统初始化完成")
	return nil
}

func (d *FyneDisplay) log(format string, args ...interface{}) {
	if d.logFile != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(d.logFile, "[%s] %s\n", timestamp, fmt.Sprintf(format, args...))
		d.logFile.Sync()
	}
}

// NewFyneDisplay 创建一个新的Fyne显示实例
func NewFyneDisplay(configManager *config.Manager, generator *idgen.Generator, processManager *process.Manager, isAdmin bool) *FyneDisplay {
	fyneApp := app.New()
	fyneApp.Settings().SetTheme(theme.DarkTheme())

	window := fyneApp.NewWindow("Cursor ID Modifier")
	window.Resize(fyne.NewSize(400, 300))

	display := &FyneDisplay{
		app:            fyneApp,
		window:         window,
		status:         widget.NewLabel(""),
		progress:       widget.NewProgressBarInfinite(),
		configManager:  configManager,
		generator:      generator,
		processManager: processManager,
		isAdmin:        isAdmin,
	}

	if err := display.setupLogging(); err != nil {
		dialog.ShowError(fmt.Errorf("无法创建日志文件: %v", err), window)
	}

	return display
}

// Initialize 初始化UI组件
func (d *FyneDisplay) Initialize() {
	d.log("初始化UI组件")
	title := widget.NewLabel("Cursor ID Modifier")
	title.TextStyle = fyne.TextStyle{Bold: true}

	var content *fyne.Container
	if !d.isAdmin && runtime.GOOS == "windows" {
		d.log("当前无管理员权限，显示权限请求界面")
		content = container.NewVBox(
			title,
			widget.NewLabel(""),
			widget.NewLabel("需要管理员权限才能运行此程序"),
			widget.NewLabel(""),
			widget.NewButton("获取管理员权限", d.handleElevation),
			widget.NewButton("退出", func() {
				d.log("用户选择退出程序")
				d.window.Close()
			}),
		)
	} else {
		d.log("显示主操作界面")
		content = container.NewVBox(
			title,
			widget.NewLabel(""),
			d.status,
			d.progress,
			widget.NewLabel(""),
			widget.NewButton("开始修改", d.handleModification),
		)
	}

	d.window.SetContent(content)
	d.progress.Hide()
}

// Run 运行UI
func (d *FyneDisplay) Run() {
	d.log("启动UI")
	d.window.ShowAndRun()
}

// ShowError 显示错误对话框
func (d *FyneDisplay) ShowError(title, message string) {
	d.log("显示错误: %s - %s", title, message)
	dialog.ShowError(fmt.Errorf(message), d.window)
}

// ShowInfo 显示信息对话框
func (d *FyneDisplay) ShowInfo(title, message string) {
	d.log("显示信息: %s - %s", title, message)
	dialog.ShowInformation(title, message, d.window)
}

// RequestPrivileges 请求管理员权限
func (d *FyneDisplay) RequestPrivileges() {
	dialog.ShowInformation("权限请求", "程序需要管理员权限才能运行\n点击确定后将请求权限", d.window)
}

func (d *FyneDisplay) handleElevation() {
	d.log("开始权限提升流程")
	exe, err := os.Executable()
	if err != nil {
		d.log("获取程序路径失败: %v", err)
		d.ShowError("错误", "无法获取程序路径: "+err.Error())
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		d.log("获取当前目录失败: %v", err)
		d.ShowError("错误", "无法获取当前目录: "+err.Error())
		return
	}

	var cmdArgs string
	if len(os.Args) > 1 {
		cmdArgs = strings.Join(os.Args[1:], " ")
	}

	var argStr string
	if cmdArgs != "" {
		argStr = fmt.Sprintf("-ArgumentList '%s'", cmdArgs)
	}

	psCmd := fmt.Sprintf(`$env:ELEVATED=1; Start-Process -FilePath '%s' %s -Verb runas`,
		exe, argStr)

	d.log("执行权限提升命令: %s", psCmd)
	cmd := exec.Command("powershell.exe", "-Command", psCmd)
	cmd.Dir = cwd

	if err := cmd.Start(); err != nil {
		d.log("权限提升失败: %v", err)
		d.ShowError("错误", "获取管理员权限失败: "+err.Error())
		return
	}

	d.log("权限提升成功，关闭当前窗口")
	d.window.Close()
}

func (d *FyneDisplay) handleModification() {
	go func() {
		d.progress.Show()
		defer d.progress.Hide()

		d.log("开始修改操作")
		d.status.SetText("正在关闭 Cursor...")
		if err := d.processManager.KillCursorProcesses(); err != nil {
			d.log("关闭Cursor进程失败: %v", err)
			d.ShowError("错误", "无法关闭 Cursor 进程，请手动关闭后重试")
			return
		}

		d.status.SetText("正在读取配置...")
		oldConfig, err := d.configManager.ReadConfig()
		if err != nil {
			d.log("读取配置失败: %v", err)
			d.ShowError("错误", "读取配置失败: "+err.Error())
			return
		}

		d.status.SetText("正在生成新ID...")
		d.log("生成新的配置ID")
		newConfig := d.generateNewConfig(oldConfig)

		d.status.SetText("正在保存配置...")
		if err := d.configManager.SaveConfig(newConfig, false); err != nil {
			d.log("保存配置失败: %v", err)
			d.ShowError("错误", "保存配置失败: "+err.Error())
			return
		}

		d.log("操作完成")
		d.status.SetText("操作完成！请重启 Cursor")
		d.ShowInfo("成功", "ID修改完成！\n请重启 Cursor 以应用更改。")
	}()
}

func (d *FyneDisplay) generateNewConfig(oldConfig *config.StorageConfig) *config.StorageConfig {
	newConfig := &config.StorageConfig{}

	if machineID, err := d.generator.GenerateMachineID(); err == nil {
		d.log("生成新的机器ID: %s", machineID)
		newConfig.TelemetryMachineId = machineID
	}

	if macMachineID, err := d.generator.GenerateMacMachineID(); err == nil {
		d.log("生成新的MAC机器ID: %s", macMachineID)
		newConfig.TelemetryMacMachineId = macMachineID
	}

	if deviceID, err := d.generator.GenerateDeviceID(); err == nil {
		d.log("生成新的设备ID: %s", deviceID)
		newConfig.TelemetryDevDeviceId = deviceID
	}

	if oldConfig != nil && oldConfig.TelemetrySqmId != "" {
		d.log("保留原有的SQM ID: %s", oldConfig.TelemetrySqmId)
		newConfig.TelemetrySqmId = oldConfig.TelemetrySqmId
	} else if sqmID, err := d.generator.GenerateSQMID(); err == nil {
		d.log("生成新的SQM ID: %s", sqmID)
		newConfig.TelemetrySqmId = sqmID
	}

	return newConfig
}
