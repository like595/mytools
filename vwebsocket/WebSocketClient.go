package vwebsocket

import (
	"github.com/gorilla/websocket"
	"github.com/like595/mytools/vtools"
	"time"
)

type VWebSocketClient struct {
	url                     string
	webSocketReceiveData    WebSocketReceiveData
	webSocketConnectBack    WebSocketConnectBack
	webSocketDisConnectBack WebSocketDisConnectBack
	conn                    *websocket.Conn
	status 	int	//鏈接狀態；1：已連接；0：未連接；
}

// 接收数据回调函数
type WebSocketReceiveData func(messageType int, data *[]byte)

// 连接成功回调函数
type WebSocketConnectBack func()

// 失去连接回调函数
type WebSocketDisConnectBack func()

// 启动
// url：WebSocket服务端地址
// webSocketReceiveData：接收到数据回调函数
func (this *VWebSocketClient) Start(url string, webSocketReceiveData WebSocketReceiveData, webSocketConnectBack WebSocketConnectBack, webSocketDisConnectBack WebSocketDisConnectBack) {
	this.url = url
	this.webSocketReceiveData = webSocketReceiveData
	this.webSocketConnectBack = webSocketConnectBack
	this.webSocketDisConnectBack = webSocketDisConnectBack
	go this.connect()
}

func (this *VWebSocketClient) connect() {
	for true {
		//func NewClient(netConn net.Conn, u *url.URL, requestHeader http.Header, readBufSize, writeBufSize int) (c *Conn, response *http.Response, err error) {
		//Dial(urlStr string, requestHeader http.Header) (*Conn, *http.Response, error)
		// 连接到WebSocket服务器
		websocket.DefaultDialer.ReadBufferSize = 1024*10
		c, _, err := websocket.DefaultDialer.Dial(this.url, nil)
		if err != nil {
			vtools.SugarLogger.Error("连接WebSocket失败。url=", this.url)
			this.status = 0
			time.Sleep(time.Second)
			this.webSocketDisConnectBack()
			continue
		} else {
			vtools.SugarLogger.Error("连接WebSocket成功。url=", this.url)
			this.conn = c
			this.status = 1
			this.webSocketConnectBack()
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
	if this.status == 1{
		// 发送消息
		err := this.conn.WriteMessage(websocket.TextMessage, *data)
		if err != nil {
			return false
		} else {
			return true
		}
	}
	return false
}

// 接收消息
func (this *VWebSocketClient) receive() {
	// 接收消息
	for true {
		messageType, message, err := this.conn.ReadMessage()
		if err != nil {
			vtools.SugarLogger.Info("WebSocketClient接收消息错误：", err)
			this.status = 0
			this.webSocketDisConnectBack()
			go this.connect()
			break
		}
		if this.webSocketReceiveData != nil {
			this.webSocketReceiveData(messageType, &message)
		}
		time.Sleep(time.Millisecond)
	}
}
