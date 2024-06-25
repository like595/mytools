package main

import (
	vModbus "github.com/like595/mytools/MyModbus"
	"github.com/like595/mytools/vtools"
	"time"
)

func main() {

	mymodbusRTU := vModbus.MyModbusTCP{}
	mymodbusRTU.Start("192.168.204.201", 6010, 5, rrreceiveDataBackClient, ccconnectBackClient, dddisConnectBackClient)
	mymodbusRTU.Read(3, 0x1234, 10)
	sdata := &[]byte{0, 1, 0, 2, 0, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 10}
	mymodbusRTU.Write(16, 0x1234, len((*sdata)), sdata)
	sdata1 := &[]byte{0, 10}
	mymodbusRTU.Write(6, 0x1234, len((*sdata1)), sdata1)
	for true {
		time.Sleep(time.Hour)
	}
}

/*
接收数据回调函数
*/
func rrreceiveDataBackClient(data []byte, len int) {
	vtools.SugarLogger.Info("接收到回应数据：", vtools.BytesToString(data))
}

/*
连接成功回调函数
*/
func ccconnectBackClient() {
	vtools.SugarLogger.Info("连接成功")
}

/*
失去连接回调函数
*/
func dddisConnectBackClient() {
	vtools.SugarLogger.Info("失去连接")
}
