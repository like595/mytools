package main

import (
	"fmt"
	"github.com/like595/mytools/vpo"
	"github.com/like595/mytools/vstatusmanager"
)

func main() {
	//gengXin(0, 0)
	//gengXin(0, 1)
	//gengXin(0, 2)
	//gengXin(0, 3)
	//gengXin(0, 4)
	//gengXin(0, 31)

	//gengXin(1, 0)
	//gengXin(1, 1)
	//gengXin(1, 2)
	//gengXin(1, 3)
	//gengXin(1, 4)
	//gengXin(1, 31)

	//gengXin(2, 0)
	//gengXin(2, 1)
	//gengXin(2, 2)
	//gengXin(2, 3)
	//gengXin(2, 4)
	//gengXin(2, 31)

	//gengXin(3, 0)
	//gengXin(3, 1)
	//gengXin(3, 2)
	//gengXin(3, 3)
	//gengXin(3, 4)
	//gengXin(3, 31)

	gengXin(4, 0)
	gengXin(4, 1)
	gengXin(4, 2)
	gengXin(4, 3)
	gengXin(4, 4)
	gengXin(4, 31)
}

func sendData() {
	fmt.Println("发送消息")
}

func gengXin(yzt, xtz int) {
	deviceStatus := vpo.DeviceStatus{}
	deviceStatus.Status = yzt
	deviceId := "1"
	statusManager := vstatusmanager.VStatusManager{}
	statusManager.Start(deviceId, &deviceStatus, sendData)
	//0-未知；1-正常；2-故障；3-告警；4-离线；31：告警恢复；
	statusManager.SetDeviceStatus(xtz)
}
