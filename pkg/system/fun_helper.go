package system

import (
	"os/exec"
	"os/user"
	"runtime"
)

func GetCurrentUser() string {
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return ""
}

func CheckAdminPrivileges() (bool, error) {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("net", "session").Run() == nil, nil
	case "darwin", "linux":
		if u, err := user.Current(); err != nil {
			return false, err
		} else {
			return u.Uid == "0", nil
		}
	default:
		return false, nil
	}
}
