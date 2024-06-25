package main

import (
	vModbus "github.com/like595/mytools/MyModbus"
	"github.com/like595/mytools/vtools"
	"time"
)

func main() {

	mymodbusRTU := vModbus.MyModbusRTU{}
	mymodbusRTU.Start("192.168.204.201", 6010, 5, rreceiveDataBackClient, cconnectBackClient, ddisConnectBackClient)
	mymodbusRTU.Read(1, 0x1234, 10)
	mymodbusRTU.Read(3, 0x1235, 10)
	//sdata := &[]byte{0, 1, 0, 2, 0, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 10}
	sdata1 := &[]byte{0, 1}
	//mymodbusRTU.Write(16, 0x1235, len((*sdata)), sdata)
	//mymodbusRTU.Write(15, 0x1235, len((*sdata)), sdata)
	mymodbusRTU.Write(6, 0x1235, len((*sdata1)), sdata1)
	for true {
		time.Sleep(time.Hour)
	}
}

/*
接收数据回调函数
*/
func rreceiveDataBackClient(data []byte, len int) {
	vtools.SugarLogger.Info("接收到回应数据：", vtools.BytesToString(data))
}

/*
连接成功回调函数
*/
func cconnectBackClient() {
	vtools.SugarLogger.Info("连接成功")
}

/*
失去连接回调函数
*/
func ddisConnectBackClient() {
	vtools.SugarLogger.Info("失去连接")
}
