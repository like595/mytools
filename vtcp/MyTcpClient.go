package vtcp


import (
	"fmt"
	"net"
	"time"
)
/*
接收数据回调函数
*/
type ReceiveDataBackClient func(data []byte,len int)
/*
连接成功回调函数
*/
type ConnectBackClient func()
/*
失去连接回调函数
*/
type DisConnectBackClient func()


type MyTcpClient struct {
	url string
	conn net.Conn
	receiveDataBackFun ReceiveDataBackClient
	connectBack ConnectBackClient
	disConnectBack DisConnectBackClient
	zt bool
}

/*
连接tcp服务端，启动函数
*/
func (this *MyTcpClient)ConnectTcpServer(url string,receiveDataBackFun ReceiveDataBackClient,connectBack ConnectBackClient,disConnectBack DisConnectBackClient){
	this.zt = true
	//go this.trunQuere()

	this.url = url
	this.receiveDataBackFun = receiveDataBackFun
	this.connectBack = connectBack
	this.disConnectBack = disConnectBack
	fmt.Println("启动tcp客户端",this.url)
	this.connectToServer()
	go this.readData()
	go this.reconnection()
}

/*
停止tcp
*/
func (this *MyTcpClient)StopTcpClient()  {
	this.conn.Close()
}

/*
连接服务端
*/
func (this *MyTcpClient)connectToServer()  {
	conn, err := net.Dial("tcp", this.url)
	if err == nil {
		this.conn = conn
		this.connectBack()
		fmt.Println("连接成功***************")
	}
}

/*
向服务端发送数据
*/
func (this *MyTcpClient)WriteData(data []byte)  bool{

	for(true){
		if this.zt{
			break
		}
		time.Sleep(time.Millisecond)
	}
	this.zt = false

	//tools.SugarLogger.Info("发送tcp数据")
	if this.conn == nil{
		this.zt = true
		return false
	}
	_,err := this.conn.Write(data)
	//time.Sleep(50*time.Millisecond)
	time.Sleep(1e3)
	if err == nil{
		this.zt = true
		return true
	}else {
		this.conn.Close()
		this.conn = nil

		this.disConnectBack()
		//this.connectToServer()
		this.zt = true
		return false
	}
}


/*
接收服务端数据
*/
func (this *MyTcpClient)readData()  {
	for true{
		buf := make([]byte, 10*1024)
		if this.conn == nil{
			time.Sleep(1e3)
			continue
		}
		len, err := this.conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			if this.conn != nil{
				this.conn.Close()
			}
			this.disConnectBack()
			time.Sleep(time.Second)
			continue
			//return //终止程序
		}
		go this.receiveDataBackFun(buf,len)
		time.Sleep(1e3)
	}
}

/*
重连线程
*/
func (this *MyTcpClient)reconnection()  {
	for true{
		time.Sleep(1e10)
		if this.conn == nil{
			this.connectToServer()
		}
	}
}

