package vModbus

import (
	"github.com/like595/mytools/vtcp"
)

// ModbusTCP工具
type VModbusTCPSync struct {
	url                        string
	address                    int
	modbusConnectBackClient    ModbusConnectBackClient
	modbusDisConnectBackClient ModbusDisConnectBackClient
	tcpClient                  vtcp.MyTcpSyncClient
	index                      int
}

// 启动
// url:格式，ip:port，设备地址，接收数据回调函数，连接成功回调函数，连接失败回调函数
func (this *VModbusTCPSync) Start(url string, address int, modbusReceiveDataBackClient ModbusReceiveDataBackClient,
	modbusConnectBackClient ModbusConnectBackClient, modbusDisConnectBackClient ModbusDisConnectBackClient) {
	this.modbusDisConnectBackClient = modbusDisConnectBackClient
	this.modbusConnectBackClient = modbusConnectBackClient
	this.address = address

	this.tcpClient = vtcp.MyTcpSyncClient{}
	go this.tcpClient.ConnectTcpServer(url, this.connectBackClient, this.disConnectBackClient, 5)

}

// 读取数据
// 功能码，起始地址，读数据长度。
// 功能码：
// 1：读线圈寄存器
// 2：读离散输入寄存器；
// 3：读保持寄存器；
// 4：读输入寄存区；
// 发送数据后，同步返回接收到的内容
func (this *VModbusTCPSync) Read(funCode int, begin int, len int) ([]byte, error) {
	data := make([]byte, 0)
	//索引
	//data = append(data, byte(ix>>8))
	//data = append(data, byte(ix))
	data = append(data, byte(begin>>8))
	data = append(data, byte(begin))
	//TCP
	data = append(data, 0x00)
	data = append(data, 0x00)
	//长度
	data = append(data, 0x00)
	data = append(data, 0x06)
	//地址
	data = append(data, byte(this.address))
	//功能码
	data = append(data, byte(funCode))
	data = append(data, byte(begin>>8))
	data = append(data, byte(begin))
	data = append(data, byte(len>>8))
	data = append(data, byte(len))

	return this.tcpClient.SendAndReceive(data)
}

// 读取数据
// 功能码，起始地址，写数据长度，数据。
// 功能码：
// 5：写单个线圈寄存器；
// 6：写单个保持寄存器；
// 15：写多个线圈寄存器；未实现
// 16：写多个保持寄存器；
// 发送数据后，同步返回接收到的内容
func (this *VModbusTCPSync) Write(funCode int, begin int, len int, sdata *[]byte) ([]byte, error) {
	len = len
	ix := this.getIndex()
	data := make([]byte, 0)
	//索引
	data = append(data, byte(ix>>8))
	data = append(data, byte(ix))
	//TCP
	data = append(data, 0x00)
	data = append(data, 0x00)

	if funCode == 5 || funCode == 6 {
		len += 5
		//功能码0x5: 写单个线圈寄存器0x05：
		//MBAP header(7字节) + 功能码(1字节) + 线圈寄存器起始地址的高位（1字节） + 线圈寄存器起始地址的低位（1字节） + 要写的值的高位（1字节） + 要写的值的低位（1字节）
		//长度
		data = append(data, byte(len>>8))
		data = append(data, byte(len))
		//地址
		data = append(data, byte(this.address))
		//功能码
		data = append(data, byte(funCode))
		data = append(data, byte(begin>>8))
		data = append(data, byte(begin))
		data = append(data, (*sdata)[0])
		data = append(data, (*sdata)[1])
	} else if funCode == 0x0F || funCode == 0x10 {
		//MBAP 功能码 + 起始地址H 起始地址L + 输出数量H 输出数量L + 字节长度 + 输出值H 输出值L

		//长度
		data = append(data, byte((len*2+7)>>8))
		data = append(data, byte(len*2+7))
		//地址
		data = append(data, byte(this.address))
		//功能码
		data = append(data, byte(funCode))
		//起始地址H 起始地址L
		data = append(data, byte(begin>>8))
		data = append(data, byte(begin))
		//输出数量H 输出数量L
		data = append(data, byte((len)>>8))
		data = append(data, byte(len))
		//字节长度
		data = append(data, byte(len*2))
		//输出值H 输出值L
		for _, b := range *sdata {
			data = append(data, b)
		}
	}

	return this.tcpClient.SendAndReceive(data)
}

/*
连接成功回调函数
*/
func (this *VModbusTCPSync) connectBackClient() {
	if this.modbusConnectBackClient != nil {
		this.modbusConnectBackClient()
	}
}

/*
失去连接回调函数
*/
func (this *VModbusTCPSync) disConnectBackClient() {
	if this.modbusDisConnectBackClient != nil {
		this.modbusDisConnectBackClient()
	}
}

/*
生成索引
*/
func (this *VModbusTCPSync) getIndex() int {
	this.index += 1
	if this.index > 65535 {
		this.index = 0
	}
	return this.index
}
