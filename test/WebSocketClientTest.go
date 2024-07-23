package main

import (
	"fmt"
	"github.com/like595/mytools/vwebsocket"
	"time"
)

func main() {
	vWebSocketClient := vwebsocket.VWebSocketClient{}
	vWebSocketClient.Start("ws://192.168.204.201:8888", webSocketReceiveData)
	for true {
		time.Sleep(time.Hour)
	}
}

func webSocketReceiveData(messageType int, data *[]byte) {
	fmt.Println("接收带哦消息：", string(*data))
}
