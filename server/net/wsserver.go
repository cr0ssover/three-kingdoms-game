package net

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/cr0ssover/three-kingdoms-game/server/utils"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
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

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		WsConn:   wsConn,
		outChan:  make(chan *WsMsgRsp, 10000),
		property: make(map[string]interface{}),
		Seq:      0,
	}
}

func (w *wsServer) Router(router *Router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}

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
		err := w.writer(msg)
		if err != nil {
			log.Println("数据写入失败:", err)
		}
	}
}

func (w *wsServer) writer(msg *WsMsgRsp) error {
	data, err := json.Marshal(msg.Body)
	if err != nil {
		log.Println(err)
	}
	secretKey, err := w.GetProperty("secretKey")
	if err == nil {
		// 数据加密
		key := secretKey.(string)
		data, err := utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		if err != nil {
			log.Println(err)
		}
		// 压缩数据
		if data, err := utils.Zip(data); err == nil {
			// 写数据
			if err = w.WsConn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				// 写入失败关闭连接
				w.Close()
				return err
			}
		} else {
			// 压缩失败返回错误
			return err
		}
	}
	return err
}

// 读数据
func (w *wsServer) readMsgLoop() {
	defer func() {
		if err := recover(); err != nil {
			w.Close()
		}
	}()

	for {
		// 读取数据
		_, data, err := w.WsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误:", err)
			break
		}
		// 数据解压
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错:", err)
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
				log.Println("数据格式有误，解密失败:", err)
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
			log.Println("解析数据失败:", err)
		} else {
			// 获取到数据后进行处理
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{
				Name: body.Name,
				Seq:  req.Body.Seq,
			}}
			w.router.Run(req, rsp)
			w.outChan <- rsp
		}
	}
	w.Close()
}

func (w *wsServer) Close() {
	_ = w.WsConn.Close()
}

const HandshakeMsg = "handshake"

func (w *wsServer) Handshake() {
	var secretKey string
	key, err := w.GetProperty(secretKey)
	if err == nil {
		secretKey = key.(string)
	} else {
		secretKey = utils.RandSeq(16)
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
