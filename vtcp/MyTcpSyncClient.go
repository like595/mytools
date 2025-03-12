package vtcp

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// 接收数据回调函数
type ReceiveSyncDataBackClient func(data []byte, len int)

// 连接成功回调函数
type ConnectSyncBackClient func()

// 失去连接回调函数
type DisConnectSyncBackClient func()

// Tcp客户端工具
type MyTcpSyncClient struct {
	url     string
	conn    net.Conn
	mu      sync.Mutex
	zt      bool
	timeout int
}

// 连接tcp服务端，启动函数
func (this *MyTcpSyncClient) ConnectTcpServer(url string,
	connectBack ConnectSyncBackClient, disConnectBack DisConnectSyncBackClient, timeout int) {
	this.timeout = timeout
	this.mu.Lock()
	defer this.mu.Unlock()
	this.url = url
	fmt.Println("启动tcp客户端", this.url)
	for !this.connectToServer() {
		time.Sleep(time.Minute)
	}
}

// 停止tcp
func (this *MyTcpSyncClient) StopTcpServer() {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.conn != nil {
		this.conn.Close()
	}
}

// 连接服务端
func (this *MyTcpSyncClient) connectToServer() bool {
	conn, err := net.DialTimeout("tcp", this.url, time.Second)
	if err == nil {
		this.conn = conn
		return true
	}
	return false
}

// 向服务端发送数据并同步接收响应
func (this *MyTcpSyncClient) SendAndReceive(data []byte) ([]byte, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.conn == nil {
		return nil, errors.New("connection is nil")
	}

	// 发送数据
	_, err := this.conn.Write(data)
	if err != nil {
		return nil, err
	}

	// 设置读取超时
	this.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))

	buf := make([]byte, 10*1024)
	len, err := this.conn.Read(buf)
	if err != nil {
		// 检查是否是超时错误
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, errors.New("timeout waiting for response")
		}
		return nil, err
	}
	return buf[:len], nil
}
