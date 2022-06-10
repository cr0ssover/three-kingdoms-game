package net

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/utils"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

// websocket 服务
type wsServer struct {
	WsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgRsp //写队列
	Seq          int64
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

var cid int64

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	s := &wsServer{
		WsConn:   wsConn,
		outChan:  make(chan *WsMsgRsp, 10000),
		property: make(map[string]interface{}),
		Seq:      0,
	}
	cid++
	s.SetProperty("cid", cid)
	return s
}

func (w *wsServer) Router(router *Router) {
	w.router = router
}

// 设置属性
func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}

// 获取属性
func (w *wsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	if value, ok := w.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	delete(w.property, key)
}
func (w *wsServer) Addr() string {
	return w.WsConn.RemoteAddr().String()
}

func (w *wsServer) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{
		Body: &RspBody{
			Name: name,
			Msg:  data,
			Seq:  0,
		},
	}
	w.outChan <- rsp
}

// 建立通道后，收发消息就要监听
func (w *wsServer) Start() {
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

// 写数据
func (w *wsServer) writeMsgLoop() {
	for {
		msg := <-w.outChan
		err := w.writer(msg.Body)
		if err != nil {
			logger.Warn("数据写入失败:", err)
		}
	}
}

func (w *wsServer) writer(msg interface{}) error {
	data, err := json.Marshal(msg.(*RspBody))
	if err != nil {
		logger.Warn(err)
	}

	// 获取密钥
	secretKey, err := w.GetProperty("secretKey")
	if err != nil {
		return err
	}
	logger.Info("服务端写入数据:", string(data))
	// 数据加密
	key := secretKey.(string)
	data, err = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	if err != nil {
		logger.Warn("ws消息加密出错,err: ", err)
		return err
	}

	// 压缩数据
	if data, err = utils.Zip(data); err != nil {
		// 压缩失败返回错误
		logger.Warn("ws消息压缩出错,err: ", err)
		return err
	}

	// 写数据
	if err = w.WsConn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		// 写入失败关闭连接
		w.Close()
		return err
	}

	return nil
}

// 读数据
func (w *wsServer) readMsgLoop() {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn("ws捕捉到异常，err: ", err)
			w.Close()
		}
	}()

	for {
		// 读取数据
		_, data, err := w.WsConn.ReadMessage()
		if err != nil {
			logger.Warn("收消息出现错误:", err)
			break
		}
		// 数据解压
		data, err = utils.UnZip(data)
		if err != nil {
			logger.Warn("解压数据出错:", err)
			continue
		}

		// 前端加密消息进行解密
		secretKey, err := w.GetProperty("secretKey")
		if err == nil {
			// 转换数据类型
			key := secretKey.(string)
			// 数据解密
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				logger.Warn("数据格式有误，解密失败:", err)
				// 出错后发起握手
				w.Handshake()
			} else {
				data = d
			}
		}

		// 反序列化数据
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			logger.Warn("解析数据失败,err", err)
		} else {
			// 获取到数据后进行处理
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{
				Name: body.Name,
				Seq:  req.Body.Seq,
			}}
			// 是否为心跳
			if req.Body.Name == "heartbeat" {
				heartbeate := &Heartbeate{}
				mapstructure.Decode(req.Body.Msg, heartbeate)
				heartbeate.Stime = time.Now().UnixNano() / 1e6
				rsp.Body.Msg = heartbeate
			} else {
				w.router.Run(req, rsp)
			}
			w.outChan <- rsp
		}
	}
	w.Close()
}

// 连接关闭
func (w *wsServer) Close() {
	_ = w.WsConn.Close()
}

const HandshakeMsg = "handshake"

func (w *wsServer) Handshake() {
	var secretKey string
	key, err := w.GetProperty(secretKey)
	if err != nil {
		secretKey = utils.RandSeq(16)
	} else {
		secretKey = key.(string)
	}

	// 发送secreKey给客户端
	handshake := &Handshake{Key: secretKey}
	body := &RspBody{Name: HandshakeMsg, Msg: handshake}
	if data, err := json.Marshal(body); err == nil {
		if secretKey != "" {
			w.SetProperty("secretKey", secretKey)
		} else {
			w.RemoveProperty("secretKey")
		}
		// 压缩数据
		if data, err = utils.Zip(data); err == nil {
			w.WsConn.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}
