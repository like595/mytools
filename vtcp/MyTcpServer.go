package service

import (
	"github.com/like595/mytools/vtools"
	"log"
	"net"
	"sync"
	"time"
)

type ReceiveDataBack func(data []byte,len int,clientURL string)
type ConnectBack func(clientURL string)
type DisConnectBack func(clientURL string)


type MyTcpServer struct {
	ip string
	port int
	clientMap	sync.Map
	receiveDataBackFun ReceiveDataBack
	connectBack ConnectBack
	disConnectBack DisConnectBack
}
//func (tn *DevicePo) ExplainData(byteArray []byte, len int)  {
func (ts *MyTcpServer)StartTcpServer(ip string,port int,receiveDataBackFun ReceiveDataBack,connectBack ConnectBack,disConnectBack DisConnectBack){
	ts.ip=ip
	ts.port=port
	ts.receiveDataBackFun = receiveDataBackFun
	ts.connectBack = connectBack
	ts.disConnectBack = disConnectBack
	vtools.SugarLogger.Info("DaoCha.DaoCha_TuoDa.","启动tcp服务器",ts.ip,ts.port)
	go ts.Accept()
}

func (ts *MyTcpServer)Accept()  {
	vtools.SugarLogger.Info("DaoCha.DaoCha_TuoDa.","启动tcp服务器线程",ts.ip,ts.port)
	address := net.TCPAddr{
		IP:   net.ParseIP(ts.ip),
		Port: ts.port,
	}
	listerer, err := net.ListenTCP("tcp4", &address)
	if err != nil {
		log.Fatal(err)
		vtools.SugarLogger.Error("DaoCha.DaoCha_TuoDa.",err)
	}
	//等等tcp连接
	for {
		conn, err := listerer.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		clientURL := conn.RemoteAddr().String()
		vtools.SugarLogger.Info("DaoCha.DaoCha_TuoDa.","远程地址：", clientURL)
		go ts.readData(conn)
		//连接回调函数
		if ts.connectBack != nil{
			ts.connectBack(clientURL)
		}
		//go echo(conn)
		ts.clientMap.Store(clientURL,conn)
		time.Sleep(1e7)
	}
}

/*
向客户端发送数据
*/
func (ts *MyTcpServer)WriteData(data []byte,clientURL string)  {
	//向单个客户端发送数据
	if clientURL != ""{
		if v,ok := ts.clientMap.Load(clientURL);ok{
			conn := (v).(net.Conn)
			if conn != nil{
				_, err := conn.Write(data)
				if err != nil {
					vtools.SugarLogger.Info("DaoCha.DaoCha_TuoDa.","发送失败1：",err)
					log.Println(err)
					conn.Close()
					//失去连接回调
					if ts.disConnectBack != nil{
						ts.disConnectBack(conn.RemoteAddr().String())
					}
					//删除该元素
					ts.clientMap.Delete(conn.RemoteAddr().String())
					return
				}
			}
		}
	}

	//向全部客户端发送数据
	ts.clientMap.Range(func(key, value interface{}) bool {
		conn := (value).(net.Conn)
		clientURL := (key).(string)
		_, err := conn.Write(data)
		if err != nil {
			vtools.SugarLogger.Info("DaoCha.DaoCha_TuoDa.","发送失败1：",err)
			log.Println(err)
			conn.Close()
			//失去连接回调
			if ts.disConnectBack != nil{
				ts.disConnectBack(clientURL)
			}
			//删除该元素
			ts.clientMap.Delete(clientURL)
		}else{
			//vtools.SugarLogger.Info("DaoCha.DaoCha_TuoDa.","发送成功")
		}
		return true
	})

}


/*
接收客户端数据
*/
func (ts *MyTcpServer)readData(conn *net.TCPConn)  {
	for true{
		buf := make([]byte, 10*1024)
		len, err := conn.Read(buf)
		if err != nil {
			vtools.SugarLogger.Info("DaoCha.DaoCha_TuoDa.","Error reading", err.Error())
			return //终止程序
		}
		go ts.receiveDataBackFun(buf,len,conn.RemoteAddr().String())

		//fmt.Printf("Received data: %v\n", string(buf[:len]))
		//fmt.Printf("Received data: %v\n", tools.BytesToString(buf[:len]))
		time.Sleep(1e3)
	}
}

/*func echo(conn *net.TCPConn) {
	tick := time.Tick(5 * time.Second) //5秒请求一次
	for now := range tick {
		n, err := conn.Write([]byte(now.String()))
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		fmt.Printf("send %d bytes to %s\n", n, conn.RemoteAddr())
	}
}*/
