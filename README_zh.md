# Cursor ID Modifier

[English](README.md) | 简体中文

基于 [go-cursor-help](https://github.com/yuaotian/go-cursor-help) 项目改进的 Cursor 编辑器 ID 修改工具，使用 Fyne GUI 框架重写了界面，提供了更好的用户体验。

## 关于杀毒软件报警说明

本程序可能会被杀毒软件误报为病毒，这是因为：

1. 程序需要修改系统文件（Cursor 的配置文件）
2. 程序需要请求管理员权限
3. 程序会自动关闭其他进程（Cursor）

**安全性说明**：
- 本程序完全开源，源代码可以在 GitHub 上查看
- 您可以自行审查代码并编译程序
- 如果不放心，可以在隔离环境中运行
- 可以添加杀毒软件白名单

## 功能特点

- 现代化的图形界面
- 自动请求管理员权限
- 自动关闭 Cursor 进程
- 生成新的设备标识
- 详细的操作日志
- 暗色主题支持

## 使用方法

1. 直接运行程序
2. 如果杀毒软件报警，请添加信任或白名单
3. 如果需要管理员权限，点击"获取管理员权限"按钮
4. 在主界面点击"开始修改"按钮
5. 等待操作完成
6. 重启 Cursor 编辑器

## 运行环境要求

- Windows 操作系统
- 需要管理员权限
- Go 1.21 或更高版本（如果需要编译）

## 编译说明

```bash
# 安装依赖
go mod tidy

# 编译程序
go build -ldflags "-H windowsgui" -o cursor-id-modifier.exe ./cmd/cursor-id-modifier
```

## 注意事项

- 修改前请保存好 Cursor 中的工作内容
- 程序会自动关闭所有 Cursor 进程
- 操作日志保存在 `logs` 目录下
- 如果担心安全问题，建议自行编译使用

## 技术栈

- [Fyne](https://fyne.io/) - 跨平台 GUI 框架
- [logrus](https://github.com/sirupsen/logrus) - 结构化日志
- Go 标准库

## 致谢

- 感谢原项目 [go-cursor-help](https://github.com/yuaotian/go-cursor-help) 提供的基础功能
- 感谢 [Fyne](https://fyne.io/) 提供的优秀 GUI 框架 