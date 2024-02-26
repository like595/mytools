package ServerManager

import (
	"bytes"
	"os/exec"
	"strings"
)

/*
读取windows服务状态；1：启动；2：停止；0：错误；
*/
func getLinuxServerStatus(serverName string) int{
	err,out := execLinuxCmd(serverName,"is-active")
	if err != nil {
		// 如果有错误，返回错误信息
		return 1
	}
	zt := 2
	if strings.Index(out.String(),"active") > -1{
		zt = 1
	}
	return zt
}

/*
控制Windows服务 1：启动；2：停止；3：重启；
*/
func setLinuxServerStatus(serverName string,status int) {
	if status == 1{
		execLinuxCmd(serverName,"start")
	}else if status == 2{
		execLinuxCmd(serverName,"stop")
	}else if status == 3{
		execLinuxCmd(serverName,"restart")
	}
}


/*
重启计算机
*/
func restartComputerLinux()  error {
	execLinuxCmd("","reboot")
	return nil
}

/*
执行linux命令
 */
func execLinuxCmd(serverName,cmds string) (error,bytes.Buffer) {
	cmd := exec.Command("systemctl", cmds, serverName)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// 运行命令
	err := cmd.Run()
	return err,out
}

