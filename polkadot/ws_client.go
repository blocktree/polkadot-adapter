/*
 * Copyright 2019 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package polkadot

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"net/http"
	"sync"
	"time"
)

//局部常量
const (
	WriteWait              = 60 * time.Second
	PongWait               = 30 * time.Second
	PingPeriod             = (PongWait * 9) / 10
	MaxMessageSize         = 1 * 1024
	DefaultTimoutSEC       = 60
	ErrNetworkDisconnected = 99999
	ErrRequestTimeout      = 99998
)

type requestEntry struct {
	method   string
	respChan chan gjson.Result
	time     int64
}

type WSClient struct {
	wm                       *WalletManager //钱包管理者
	url                      string
	debug                    bool
	addr                     *string
	c                        *websocket.Conn
	isConnect                bool
	closeOnce                sync.Once
	reconnect                chan bool
	disconnected             chan struct{}
	stopWritePump            chan struct{}
	_send                    chan []byte
	nonceGen                 *snowflake.Node        //nonce生成器
	mu                       sync.RWMutex           //读写锁
	timeout                  time.Duration          //超时时间
	startRequestTimeoutCheck bool                   //是否启动了请求超时检查
	requestQueue             map[int64]requestEntry //节点的请求队列
}

func NewWSClient(wm *WalletManager, url string, timeoutSEC int, debug bool) *WSClient {

	nonceGen, _ := snowflake.NewNode(1)

	ws := &WSClient{
		wm:            wm,
		url:           url,
		debug:         debug,
		nonceGen:      nonceGen,
		timeout:       time.Duration(timeoutSEC) * time.Second,
		requestQueue:  make(map[int64]requestEntry),
		reconnect:     make(chan bool),
		disconnected:  make(chan struct{}),
		stopWritePump: make(chan struct{}),
		_send:         make(chan []byte, MaxMessageSize),
	}

	//开启自动重连
	go ws.autoReconnectNode()
	//延迟1秒，确保连接上后，可以马上调用Call
	time.Sleep(1 * time.Second)
	return ws
}

func (ws *WSClient) Call(method string, request map[string]interface{}) (*gjson.Result, error) {

	var (
		err      error
		respChan = make(chan gjson.Result)
	)

	//检查是否已经连接服务
	if !ws.isConnect {
		return nil, fmt.Errorf("websocket client is disconnected")
		//err = ws.connectNode()
		//if err != nil {
		//	return nil, err
		//}
	}

	if request == nil {
		request = make(map[string]interface{})
	}

	//添加请求队列到Map，处理完成回调方法
	nonce := ws.nonceGen.Generate().Int64()
	//nonce := int64(99)
	time := time.Now().Unix()
	//封装数据包
	request["id"] = fmt.Sprintf("%d", nonce)
	request["command"] = method

	//添加请求到队列，异步或同步等待结果，应该在发送前就添加请求，如果发送失败，删除请求
	err = ws.addRequest(nonce, time, method, respChan)
	if err != nil {
		return nil, err
	}

	//向节点发送请求
	err = ws.send(request)
	if err != nil {
		//发送失败移除请求
		ws.removeRequest(nonce)
		return nil, err
	}

	//等待返回
	result := <-respChan

	return &result, nil
}

//AddRequest 添加请求到队列
func (ws *WSClient) addRequest(nonce int64, time int64, method string, respChan chan gjson.Result) error {

	ws.mu.Lock()
	defer ws.mu.Unlock()

	if !ws.isConnect {
		return fmt.Errorf("peer has been closed")
	}

	if !ws.startRequestTimeoutCheck {
		go ws.timeoutRequestHandle()
	}

	if ws.requestQueue == nil {
		ws.requestQueue = make(map[int64]requestEntry)
	}

	if _, exist := ws.requestQueue[nonce]; exist {
		return fmt.Errorf("OWTP: nonce exist. ")
	}

	ws.requestQueue[nonce] = requestEntry{method: method, respChan: respChan, time: time}

	return nil
}

//RemoveRequest 移除请求
func (ws *WSClient) removeRequest(nonce int64) error {

	ws.mu.Lock()
	defer ws.mu.Unlock()

	delete(ws.requestQueue, nonce)

	return nil
}

//resetRequestQueue 重置请求队列
func (ws *WSClient) resetRequestQueue() {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	//处理所有未完成的请求，返回连接断开的异常
	for n, r := range ws.requestQueue {
		resp, err := responseError(r.method, "network disconnected", ErrNetworkDisconnected)
		if err != nil {
			continue
		}
		r.respChan <- resp
		delete(ws.requestQueue, n)
	}
}

//connectNode 连接节点
func (ws *WSClient) connectNode() error {

	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 60 * time.Second,
	}

	c, _, err := dialer.Dial(ws.url, nil)
	if err != nil {
		return err
	}

	ws.c = c
	ws.isConnect = true
	ws.openPipe()

	return nil
}

//autoReconnectTransmitNode 自动重连
func (ws *WSClient) autoReconnectNode() error {

	var (
		err error
		//重连时的等待时间
		reconnectWait = 5
	)

	//连接状态通道
	ws.reconnect = make(chan bool, 1)
	//断开状态通道
	ws.disconnected = make(chan struct{}, 1)

	defer func() {
		close(ws.reconnect)
		close(ws.disconnected)
	}()

	//启动连接
	ws.reconnect <- true

	//节点运行时
	for {
		select {
		case <-ws.reconnect:
			//重新连接
			ws.wm.Log.Info("Connecting to XRP node")
			err = ws.connectNode()
			if err != nil {
				ws.wm.Log.Errorf("Connect XRP node failed unexpected error: %v", err)
				ws.disconnected <- struct{}{}
			} else {
				ws.wm.Log.Infof("Connect XRP node successfully.")
			}
			//ws.Call("world", nil)
		case <-ws.disconnected:
			//重新连接，前等待
			ws.wm.Log.Info("Auto reconnect after", reconnectWait, "seconds...")
			time.Sleep(time.Duration(reconnectWait) * time.Second)
			ws.reconnect <- true
		}
	}

	return nil
}

//Close 关闭连接
func (ws *WSClient) close() error {
	var err error

	//保证节点只关闭一次
	//ws.closeOnce.Do(func() {

	if !ws.isConnect {
		//ws.wm.Log.Debug("end close")
		return nil
	}

	ws.wm.Log.Infof("websocket connection closed")

	err = ws.c.Close()
	ws.isConnect = false
	//断开需要重置所有请求
	ws.resetRequestQueue()
	//断开连接通知
	ws.stopWritePump <- struct{}{}
	ws.disconnected <- struct{}{}
	//})
	return err
}

//Send 发送消息
func (ws *WSClient) send(data map[string]interface{}) error {

	//ws.wm.Log.Emergency("Send DataPacket:", data)
	respBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ws._send <- respBytes
	return nil
}

//OpenPipe 打开通道
func (ws *WSClient) openPipe() error {

	//if ws.debug {
		ws.wm.Log.Debug("openPipe")
	//}

	if !ws.isConnect {
		return fmt.Errorf("client is not connect")
	}

	//发送通道
	go ws.writePump()

	//监听消息
	go ws.readPump()

	return nil
}

// WritePump 发送消息通道
func (ws *WSClient) writePump() {

	//if ws.debug {
		ws.wm.Log.Debug("writePump start")
	//}

	ticker := time.NewTicker(PingPeriod) //发送心跳间隔事件要<等待时间
	defer func() {
		ticker.Stop()
		//ws.close()
		//if ws.debug {
			ws.wm.Log.Debug("writePump end")
		//}
	}()
	for {
		select {
		case <-ws.stopWritePump: //断开则停止写入管道
			return
		case message, ok := <-ws._send:
			//发送消息
			if !ok {
				ws.write(websocket.CloseMessage, []byte{})
				continue
			}
			if ws.debug {
				ws.wm.Log.Debug("Send: ", string(message))
			}
			if err := ws.write(websocket.TextMessage, message); err != nil {
				continue
			}
		case <-ticker.C:
			//定时器的回调,发送心跳检查,
			err := ws.write(websocket.PingMessage, []byte{})

			if err != nil {
				continue //客户端不响应心跳就停止
			}
		}
	}
}

// write 输出数据
func (ws *WSClient) write(mt int, message []byte) error {
	ws.c.SetWriteDeadline(time.Now().Add(WriteWait)) //设置发送的超时时间点
	return ws.c.WriteMessage(mt, message)
}

// ReadPump 监听消息
func (ws *WSClient) readPump() {

	//if ws.debug {
		ws.wm.Log.Debug("readPump start")
	//}

	ws.c.SetReadDeadline(time.Now().Add(PongWait)) //设置客户端心跳响应的最后限期
	ws.c.SetPongHandler(func(string) error {
		ws.c.SetReadDeadline(time.Now().Add(PongWait)) //设置下一次心跳响应的最后限期
		return nil
	})
	defer func() {
		ws.close()
		//if ws.debug {
			ws.wm.Log.Debug("readPump end")
		//}
	}()

	for {
		_, message, err := ws.c.ReadMessage()
		if err != nil {
			ws.wm.Log.Error("Read unexpected error: ", err)
			return
		}

		if ws.debug {
			ws.wm.Log.Debug("Read: ", string(message))
		}

		//开一个goroutine处理消息
		go ws.onNewDataReceived(gjson.ParseBytes(message))

	}
}

func (ws *WSClient) onNewDataReceived(response gjson.Result) {
	nonce := response.Get("id").Int()
	//wsType := response.Get("type").String()
	ws.mu.Lock()
	//nonce := int64(99)
	f, ok := ws.requestQueue[nonce]
	if ok {
		f.respChan <- response
		delete(ws.requestQueue, nonce)
	}
	ws.mu.Unlock()
}

// timeoutRequestHandle 超时请求检查
func (ws *WSClient) timeoutRequestHandle() {
	if ws.timeout == 0 {
		ws.timeout = DefaultTimoutSEC * time.Second
	}

	period := (ws.timeout * 6) / 10
	//ws.wm.Log.Debug("mux.timeout:", mux.timeout)
	//ws.wm.Log.Debug("period:", period)
	ticker := time.NewTicker(period) //检查超时过程要<超时时间
	defer func() {
		ticker.Stop()
	}()

	ws.startRequestTimeoutCheck = true

	for {
		select {
		case <-ticker.C:
			//ws.wm.Log.Printf("check request timeout \n")
			ws.mu.Lock()
			//定时器的回调，处理超时请求
			//如果节点已经关闭马上取消请求

			for n, r := range ws.requestQueue {

				currentServerTime := time.Now()

				//计算客户端过期时间
				requestTimestamp := time.Unix(r.time, 0)
				expiredTime := requestTimestamp.Add(ws.timeout)

				//ws.wm.Log.Printf("requestTimestamp = %s \n", requestTimestamp.String())
				//ws.wm.Log.Printf("currentServerTime = %s \n", currentServerTime.String())
				//ws.wm.Log.Printf("expiredTime = %s \n", expiredTime.String())

				if currentServerTime.Unix() > expiredTime.Unix() {
					//ws.wm.Log.Printf("request expired time")
					//返回超时响应
					errInfo := fmt.Sprintf("request timeout over %s", ws.timeout.String())
					resp, err := responseError(r.method, errInfo, ErrRequestTimeout)
					if err != nil {
						continue
					}
					r.respChan <- resp
					delete(ws.requestQueue, n)
				}

			}

			ws.mu.Unlock()
		}
	}
}

//responseError 返回一个错误数据包
func responseError(method, message string, status uint64) (gjson.Result, error) {

	/*
		{
			"error": "unknownCmd",
			"error_code": 32,
			"error_message": "Unknown method.",
			"id": 1,
			"status": "error",
			"type": "response"
		}
	*/

	resp := map[string]interface{}{
		"error":         method,
		"error_code":    status,
		"error_message": message,
		"status":        "error",
		"type":          "response",
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return gjson.Result{}, err
	}

	r := gjson.ParseBytes(respJson)

	return r, nil
}
