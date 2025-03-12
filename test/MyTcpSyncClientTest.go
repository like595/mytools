package main

import (
	"fmt"
	"github.com/like595/mytools/vtcp"
	"time"
)

// 接收数据回调函数
func receiveSyncDataBackClients(data []byte, len int) {
	// 接收数据回调函数
	fmt.Println("ReceiveSyncDataBackClient", string(data[:len]))
}

// 连接成功回调函数
func ConnectSyncBackClients() {
	fmt.Println("ConnectSyncBackClient")
}

// 失去连接回调函数
func DisConnectSyncBackClients() {
	fmt.Println("DisConnectSyncBackClient")
}

var syncClient vtcp.MyTcpSyncClient

func main() {
	fmt.Println("Hello World")
	syncClient = vtcp.MyTcpSyncClient{}
	syncClient.ConnectTcpServer("127.0.0.1:9601", ConnectSyncBackClients, DisConnectSyncBackClients, 30)
	time.Sleep(time.Second)
	go helloWord()
	go fuck()
	for true {
		time.Sleep(time.Hour)
	}
}

func helloWord() {
	index := 0
	for {
		index++
		data, err := syncClient.SendAndReceive([]byte(fmt.Sprintf("你好，世界 %d", index)))
		fmt.Println("Receive:", string(data), "****", err)
		time.Sleep(time.Second)
	}
}

func fuck() {
	index := 0
	for {
		index++
		data, err := syncClient.SendAndReceive([]byte(fmt.Sprintf("fuck %d", index)))
		fmt.Println("Receive:", string(data), "****", err)
		time.Sleep(time.Second)
	}
}
