# Cursor ID Modifier

English | [简体中文](README_zh.md)

An improved version of the Cursor editor ID modification tool based on [go-cursor-help](https://github.com/yuaotian/go-cursor-help), featuring a redesigned interface using the Fyne GUI framework for better user experience.

## Features

- Modern graphical interface
- Automatic administrator privileges request
- Automatic Cursor process termination
- New device identifier generation
- Detailed operation logging
- Dark theme support

## Usage

1. Run the program directly
2. Click "Get Administrator Privileges" button if needed
3. Click "Start Modification" button in the main interface
4. Wait for the operation to complete
5. Restart Cursor editor

## Requirements

- Windows operating system
- Administrator privileges required
- Go 1.21 or higher (if compilation needed)

## Build Instructions

```bash
# Install dependencies
go mod tidy

# Build the program
go build -ldflags "-H windowsgui" -o cursor-id-modifier.exe ./cmd/cursor-id-modifier
```

## Important Notes

- Save your work in Cursor before modification
- The program will automatically close all Cursor processes
- Operation logs are saved in the `logs` directory

## Tech Stack

- [Fyne](https://fyne.io/) - Cross-platform GUI framework
- [logrus](https://github.com/sirupsen/logrus) - Structured logging
- Go standard library

## Acknowledgments

- Thanks to the original project [go-cursor-help](https://github.com/yuaotian/go-cursor-help) for providing the base functionality
- Thanks to [Fyne](https://fyne.io/) for providing the excellent GUI framework

