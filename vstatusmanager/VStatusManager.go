package vstatusmanager

import (
	"github.com/like595/mytools/vdbhelper"
	"github.com/like595/mytools/vpo"
	"time"
)

// 发送设备状态变更信息
type SendDeviceStatusData func()

// 设备状态管理
type VStatusManager struct {
	//设备ID
	deviceId string
	//设备状态；0-未知；1-正常；2-故障；3-告警；4-离线；31：告警恢复；
	deviceStatus *vpo.DeviceStatus
	//发送设备状态变更信息函数
	sendDeviceStatusData SendDeviceStatusData
}

// 启动状态管理；
// deviceId：设备id
// deviceStatus：设备状态信息
func (this *VStatusManager) Start(deviceId string, deviceStatus *vpo.DeviceStatus, sendDeviceStatusData SendDeviceStatusData) {
	this.deviceId = deviceId
	this.deviceStatus = deviceStatus
	//故障设备，更新成离线，重新检测状态
	//if this.deviceStatus.Status == 2{
	//	this.deviceStatus.Status = 4
	//}

	this.sendDeviceStatusData = sendDeviceStatusData
}

// 设置设备状态
// 0-未知；1-正常；2-故障；3-告警；4-离线；31：告警恢复；
func (this *VStatusManager) SetDeviceStatus(status int) {
	//ys := this.deviceStatus.Status
	isChange := true
	//状态发生改变
	if this.deviceStatus.Status != status {

		if status == 31 {
			//告警恢复，更改设备状态为正常
			status = 1
		} else if status == 2 && this.deviceStatus.Status == 4 {
			//已经离线的设备，不能更新成故障
			isChange = false
		} else if this.deviceStatus.Status == 3 {
			if status == 1 {
				isChange = false
			}
		}
		if isChange {
			this.deviceStatus.Status = status
			this.deviceStatus.ChangeTime = time.Now().Format("2006-01-02 15:04:05")
			if this.sendDeviceStatusData != nil {
				this.sendDeviceStatusData()
			}
			//保存数据库
			dbHelper := vdbhelper.MySqlDBHelper{}
			dbHelper.Open()
			sql := "update a_device set device_status=?,device_statsu_time=? where id=?"
			dbHelper.Exec(sql, this.deviceStatus.Status, this.deviceStatus.ChangeTime, this.deviceId)
			dbHelper.Close()
		}

	} /*else {
		isChange = false
	}*/
	//sprintf := fmt.Sprintf("要更新成%d。原始是%d，更新后是%d。", status, ys, this.deviceStatus.Status)
	//fmt.Println(sprintf, "***", isChange)
}
