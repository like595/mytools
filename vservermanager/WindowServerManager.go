package ServerManager

import (
	"fmt"
	"github.com/like595/mytools/vtools"
	"golang.org/x/sys/windows"
	_ "golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"os/exec"
	"time"
)

/*
读取windows服务状态；1：启动；2：停止；0：错误；
 */
func getWindowsServerStatus(serverName string) int{
	m, err := mgr.Connect()
	if err != nil {
		vtools.SugarLogger.Error(fmt.Sprintf("连接【%s】失败，错误: %v\n",serverName, err))
		return 0
	}
	defer m.Disconnect()

	s, err := m.OpenService(serverName) // 替换为你的服务名
	if err != nil {
		vtools.SugarLogger.Error(fmt.Sprintf("启动【%s】失败，错误: %v\n",serverName, err))
		return 0
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		vtools.SugarLogger.Error(fmt.Sprintf("打开服务【%s】失败: %v\n",serverName, err))
		return 0
	}
	zt := 0
	//停止
	if status.State == 1{
		zt = 2
	}else if status.State == 4{
		//正在运行
		zt = 1
	}
	return zt
}

/*
控制Windows服务 1：启动；2：停止；3：重启；
 */
func setWindowsServerStatus(serverName string,status int) {
	m, err := mgr.Connect()
	if err != nil {
		vtools.SugarLogger.Error("Failed to connect to service control manager: %v\n", err)
		return
	}
	defer m.Disconnect()



	s, err := m.OpenService(serverName)
	if err != nil {
		vtools.SugarLogger.Error("Failed to open service: %v\n", err)
		return
	}
	defer s.Close()

	if status == 1{
		//启动
		// 启动服务
		if err := s.Start(); err != nil {
			vtools.SugarLogger.Error("Failed to start service: %v\n", err)
		} else {
			vtools.SugarLogger.Error("Service %s started successfully.\n", serverName)
		}
	}else if status == 2{
		//停止
		_, err = s.Control(windows.SERVICE_CONTROL_STOP)
		if err != nil {
			vtools.SugarLogger.Error("Failed to stop service: %v\n", err)
		} else {
			vtools.SugarLogger.Error("Service %s stop successfully.\n", serverName)
		}
	}else if status == 3{
		//重启
		//停止
		_, err = s.Control(windows.SERVICE_CONTROL_STOP)
		if err != nil {
			vtools.SugarLogger.Error("Failed to stop service: %v\n", err)
		} else {
			vtools.SugarLogger.Error("Service %s stop successfully.\n", serverName)
		}
		time.Sleep(time.Second)
		// 启动服务
		if err := s.Start(); err != nil {
			vtools.SugarLogger.Error("Failed to start service: %v\n", err)
		} else {
			vtools.SugarLogger.Error("Service %s started successfully.\n", serverName)
		}
	}
}

/*
重启计算机
*/
func restartComputerWindows()  error {
	// 构建shutdown命令，使用-r参数表示重启，-t参数表示在多少秒之后执行，这里设置为0表示立即执行
	cmd := exec.Command("shutdown", "-r", "-t", "0")

	// 运行命令
	err := cmd.Run()
	if err != nil {
		// 如果命令执行失败，打印错误信息
		fmt.Println("Error restarting the system:", err)
		return err
	}

	// 命令执行成功，但通常这里不会执行，因为系统会立即重启
	fmt.Println("System is restarting...")
	return nil
}

