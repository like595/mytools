package vModbus

import (
	"fmt"
	"github.com/like595/mytools/vtcp"
	"github.com/like595/mytools/vtools"
)

type MyModbusRTU struct {
	ip                          string
	port                        int
	address                     int
	modbusReceiveDataBackClient ModbusReceiveDataBackClient
	modbusConnectBackClient     ModbusConnectBackClient
	modbusDisConnectBackClient  ModbusDisConnectBackClient
	tcpClient                   vtcp.MyTcpClient
}

/*
接收数据回调函数
*/
type ModbusReceiveDataBackClient func(data []byte, len int)

/*
连接成功回调函数
*/
type ModbusConnectBackClient func()

/*
失去连接回调函数
*/
type ModbusDisConnectBackClient func()

/*
*
启动
参数：ip，端口号，设备地址
*/
func (this *MyModbusRTU) Start(ip string, port, address int, modbusReceiveDataBackClient ModbusReceiveDataBackClient,
	modbusConnectBackClient ModbusConnectBackClient, modbusDisConnectBackClient ModbusDisConnectBackClient) {
	this.modbusReceiveDataBackClient = modbusReceiveDataBackClient
	this.modbusDisConnectBackClient = modbusDisConnectBackClient
	this.modbusConnectBackClient = modbusConnectBackClient
	this.address = address

	this.tcpClient = vtcp.MyTcpClient{}
	this.tcpClient.ConnectTcpServer(fmt.Sprintf("%s:%d", ip, port), this.receiveDataBackClient, this.connectBackClient, this.disConnectBackClient)

}

/*
读取数据
功能码，起始地址，读数据长度。
功能码：
1：读线圈寄存器
2：读离散输入寄存器；
3：读保持寄存器；
4：读输入寄存区；
*/
func (this *MyModbusRTU) Read(funCode int, begin int, len int) {
	data := make([]byte, 0)
	data = append(data, byte(this.address))
	data = append(data, byte(funCode))
	data = append(data, byte(begin>>8))
	data = append(data, byte(begin))
	data = append(data, byte(len>>8))
	data = append(data, byte(len))
	crc := vtools.CalculateCRC16(data)
	data = append(data, crc[0])
	data = append(data, crc[1])

	vtools.SugarLogger.Info("发送数据：", vtools.BytesToString(data))
	this.tcpClient.WriteData(data)
}

/*
读取数据
功能码，起始地址，写数据长度，数据。
功能码：
5：写单个线圈寄存器；
6：写单个保持寄存器；
15：写多个线圈寄存器；未实现
16：写多个保持寄存器；
*/
func (this *MyModbusRTU) Write(funCode int, begin int, len int, sdata *[]byte) {
	len = len / 2
	data := make([]byte, 0)
	data = append(data, byte(this.address))
	data = append(data, byte(funCode))
	data = append(data, byte(begin>>8))
	data = append(data, byte(begin))
	if funCode == 5 || funCode == 6 {
		data = append(data, (*sdata)[0])
		data = append(data, (*sdata)[1])
	} else if funCode == 15 || funCode == 16 {
		data = append(data, byte(len>>8))
		data = append(data, byte(len))
		data = append(data, byte(len*2))
		for _, b := range *sdata {
			data = append(data, b)
		}
	}

	crc := vtools.CalculateCRC16(data)
	data = append(data, crc[0])
	data = append(data, crc[1])

	vtools.SugarLogger.Info("发送数据：", vtools.BytesToString(data))
	this.tcpClient.WriteData(data)
}

/*
接收数据回调函数
*/
func (this *MyModbusRTU) receiveDataBackClient(data []byte, len int) {
	if this.modbusReceiveDataBackClient != nil {
		this.modbusReceiveDataBackClient(data, len)
	}

}

/*
连接成功回调函数
*/
func (this *MyModbusRTU) connectBackClient() {
	if this.modbusConnectBackClient != nil {
		this.modbusConnectBackClient()
	}

}

/*
失去连接回调函数
*/
func (this *MyModbusRTU) disConnectBackClient() {
	if this.modbusDisConnectBackClient != nil {
		this.modbusDisConnectBackClient()
	}
}
