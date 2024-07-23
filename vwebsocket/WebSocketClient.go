package vwebsocket

import (
	"github.com/gorilla/websocket"
	"github.com/like595/mytools/vtools"
	"time"
)

type VWebSocketClient struct {
	url                  string
	webSocketReceiveData WebSocketReceiveData
	conn                 *websocket.Conn
}

// 接收数据回调函数
type WebSocketReceiveData func(messageType int, data *[]byte)

// 启动
// url：WebSocket服务端地址
// webSocketReceiveData：接收到数据回调函数
func (this *VWebSocketClient) Start(url string, webSocketReceiveData WebSocketReceiveData) {
	this.url = url
	this.webSocketReceiveData = webSocketReceiveData
	go this.connect()
}

func (this *VWebSocketClient) connect() {
	for true {
		// 连接到WebSocket服务器
		c, _, err := websocket.DefaultDialer.Dial(this.url, nil)
		if err != nil {
			vtools.SugarLogger.Error("连接WebSocket失败。url=", this.url)
			time.Sleep(time.Second)
			continue
		} else {
			vtools.SugarLogger.Error("连接WebSocket成功。url=", this.url)
			this.conn = c
			go this.receive()
			break
		}
	}
}

// 停止
func (this *VWebSocketClient) Stop() {
	this.conn.Close()
}

// 发送WebSocket消息
func (this *VWebSocketClient) SendData(data *[]byte) bool {
	// 发送消息
	err := this.conn.WriteMessage(websocket.TextMessage, *data)
	if err != nil {
		return false
	} else {
		return true
	}
}

// 接收消息
func (this *VWebSocketClient) receive() {
	// 接收消息
	for true {
		messageType, message, err := this.conn.ReadMessage()
		if err != nil {
			vtools.SugarLogger.Info("WebSocketClient接收消息错误：", err)
			go this.connect()
			break
		}
		if this.webSocketReceiveData != nil {
			this.webSocketReceiveData(messageType, &message)
		}
		time.Sleep(time.Millisecond)
	}
}
