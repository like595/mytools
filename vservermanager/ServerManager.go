package vServerManager

import (
	"fmt"
	"net"
	"runtime"
	"time"
)

/*
读取windows服务状态；1：启动；2：停止；0：错误；
*/
func GeServerStatus(serverName string) int {
	if getSystemType() == 1 {
		return getWindowsServerStatus(serverName)
	} else if getSystemType() == 3 {
		return getLinuxServerStatus(serverName)
	}
	return 0
}

/*
控制Windows服务 1：启动；2：停止；3：重启；
*/
func SetServerStatus(serverName string, status int) {
	if getSystemType() == 1 {
		setWindowsServerStatus(serverName, status)
	} else if getSystemType() == 3 {
		setLinuxServerStatus(serverName, status)
	}
}

/*
重启计算机
*/
func RestartComputer() error {
	if getSystemType() == 1 {
		return restartComputerWindows()
	} else if getSystemType() == 3 {
		return restartComputerLinux()
	}
	return nil
}

/*
获取操作系统类型；1-Windows；2-macos；3-linux；
*/
func getSystemType() int {
	os := runtime.GOOS
	switch os {
	case "windows":
		return 1
	case "darwin":
		return 2
	case "linux":
		return 3
	default:
		return 0
	}
}

/*
监测端口号状态；true：占用；false：未占用；
*/
func GetPortStatus(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 1*time.Second)
	if err != nil {
		// 如果连接失败，端口可能未被占用
		return false
	}
	// 如果连接成功，关闭连接并返回false
	conn.Close()
	return true
}
