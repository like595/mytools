package vModbus

import (
	"github.com/like595/mytools/vtcp"
	"github.com/like595/mytools/vtools"
)

type VModbusRTU struct {
	url                         string
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

// 启动
// url:格式，ip:port，设备地址，接收数据回调函数，连接成功回调函数，连接失败回调函数
func (this *VModbusRTU) Start(url string, address int, modbusReceiveDataBackClient ModbusReceiveDataBackClient,
	modbusConnectBackClient ModbusConnectBackClient, modbusDisConnectBackClient ModbusDisConnectBackClient) {
	this.modbusReceiveDataBackClient = modbusReceiveDataBackClient
	this.modbusDisConnectBackClient = modbusDisConnectBackClient
	this.modbusConnectBackClient = modbusConnectBackClient
	this.address = address

	this.tcpClient = vtcp.MyTcpClient{}
	this.tcpClient.ConnectTcpServer(url, this.receiveDataBackClient, this.connectBackClient, this.disConnectBackClient)

}

// 读取数据
// 功能码，起始地址，读数据长度。
// 功能码：
// 1：读线圈寄存器
// 2：读离散输入寄存器；
// 3：读保持寄存器；
// 4：读输入寄存区；
func (this *VModbusRTU) Read(funCode int, begin int, len int) {
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

	this.tcpClient.WriteData(data)
}

// 读取数据
// 功能码，起始地址，写数据长度，数据。
// 功能码：
// 5：写单个线圈寄存器；
// 6：写单个保持寄存器；
// 15：写多个线圈寄存器；未实现
// 16：写多个保持寄存器；
func (this *VModbusRTU) Write(funCode int, begin int, len int, sdata *[]byte) {
	len = len / 2
	data := make([]byte, 0)
	data = append(data, byte(this.address))
	data = append(data, byte(funCode))
	data = append(data, byte(begin>>8))
	data = append(data, byte(begin))
	if funCode == 5 || funCode == 6 {
		data = append(data, (*sdata)[0])
		data = append(data, (*sdata)[1])
	} else if funCode == 16 {
		data = append(data, byte(len>>8))
		data = append(data, byte(len))
		data = append(data, byte(len*2))
		for _, b := range *sdata {
			data = append(data, b)
		}
	} else if funCode == 15 {
		data = append(data, byte(len>>8))
		data = append(data, byte(len))
		if len % 8 == 0	{
			data = append(data, byte(len/8))
		}else {
			data = append(data, byte(len/8 + 1))
		}
		for _, b := range *sdata {
			data = append(data, b)
		}
	}

	crc := vtools.CalculateCRC16(data)
	data = append(data, crc[0])
	data = append(data, crc[1])

	this.tcpClient.WriteData(data)
}

/*
接收数据回调函数
*/
func (this *VModbusRTU) receiveDataBackClient(data []byte, len int) {
	if this.modbusReceiveDataBackClient != nil {
		this.modbusReceiveDataBackClient(data, len)
	}

}

/*
连接成功回调函数
*/
func (this *VModbusRTU) connectBackClient() {
	if this.modbusConnectBackClient != nil {
		this.modbusConnectBackClient()
	}

}

/*
失去连接回调函数
*/
func (this *VModbusRTU) disConnectBackClient() {
	if this.modbusDisConnectBackClient != nil {
		this.modbusDisConnectBackClient()
	}
}
