package net

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/utils"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type syncCtx struct {
	// Groutine 的上下文 包含Goroutine的运行状态
	ctx     context.Context
	cancel  context.CancelFunc
	outChan chan *RspBody
}

func NewSyncCtx() *syncCtx {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	return &syncCtx{
		ctx:     ctx,
		cancel:  cancel,
		outChan: make(chan *RspBody),
	}
}

type ClientConn struct {
	wsConn        *websocket.Conn
	isClosed      bool                   // 监听客户端是否为关闭状态
	property      map[string]interface{} // 属性
	propertyLock  sync.RWMutex
	Seq           int64     //
	handshake     bool      // 握手状态
	handshakeChan chan bool // 握手通道
	onPush        func(conn *ClientConn, body *RspBody)
	onClose       func(conn *ClientConn)
	syncCtxMap    map[int64]*syncCtx
	syncCtxLock   sync.RWMutex
}

func NewClientConn(wsConn *websocket.Conn) *ClientConn {
	return &ClientConn{
		wsConn:        wsConn,
		handshakeChan: make(chan bool),
		Seq:           0,
		isClosed:      false,
		property:      make(map[string]interface{}),
		syncCtxMap:    map[int64]*syncCtx{},
	}
}

func (c *ClientConn) Start() bool {
	c.handshake = false
	// 读消息
	go c.wsReadLoop()
	return c.waitHandShake()
}

// 等待握手连接
func (c *ClientConn) waitHandShake() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !c.handshake {
		select {
		case <-c.handshakeChan:
			log.Println("握手成功")
			return true
		case <-ctx.Done():
			log.Println("握手超时")
			return false
		}
	}

	return true
}

// 读消息
func (c *ClientConn) wsReadLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("捕捉到异常", err)
			c.Close()
		}
	}()

	for {
		_, data, err := c.wsConn.ReadMessage()
		if err != nil {
			log.Println("接收消息出现错误,err: ", err)
			break
		}
		// 解压数据
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("数据解压出错,非法格式,err: ", err)
			continue
		}

		// 解密消息
		secretKey, err := c.GetProperty("secretKey")
		if err == nil {
			key := secretKey.(string)
			data, err = utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("数据格式有误，解密失败,err: ", err)
				return
			}
		}
		// 反序列化数据
		body := &RspBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("反序列化数据失败，err: ", err)
			return
		}

		// seq为0为握手
		if body.Seq == 0 {
			if body.Name == HandshakeMsg {
				// 获取密钥
				hs := &Handshake{}
				mapstructure.Decode(body.Msg, hs)
				if hs.Key != "" {
					c.SetProperty("secretKey", hs.Key)
				} else {
					c.RemoveProperty("secretKey")
				}
				c.handshake = true
				c.handshakeChan <- true
			} else {
				if c.onPush != nil {
					c.onPush(c, body)
				}
			}
		}

		c.syncCtxLock.RLock()
		ctx, ok := c.syncCtxMap[body.Seq]
		c.syncCtxLock.RUnlock()
		if ok {
			ctx.outChan <- body
		} else {
			if body.Seq > 0 {
				log.Println("no seq syncCtx find")
			}
		}

	}
	// 关闭连接
	c.Close()
}

// 等待消息处理
func (s *syncCtx) wait() *RspBody {
	defer s.cancel()
	select {
	case msg := <-s.outChan:
		fmt.Println("代理服务器发送来的数据，msg: ", msg)
		return msg
	case <-s.ctx.Done():
		log.Println("代理服务响应超时")
		return nil
	}
}

func (c *ClientConn) Close() {
	_ = c.wsConn.Close()
}

// 获取属性
func (c *ClientConn) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 设置属性
func (c *ClientConn) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// 移除属性
func (c *ClientConn) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

// 获取连接IP地址
func (c *ClientConn) Addr() string {
	return c.wsConn.RemoteAddr().String()
}

// 发送消息
func (c *ClientConn) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{
		Body: &RspBody{
			Name: name,
			Msg:  data,
			Seq:  0,
		},
	}
	c.write(rsp.Body)
}

// 向websocket写入数据
func (c *ClientConn) write(body interface{}) error {
	// 反序列化数据
	data, err := json.Marshal(body)
	if err != nil {
		log.Println("序列化数据出错,err: ", err)
		return err
	}
	// 获取密钥
	secretKey, err := c.GetProperty("secretKey")
	if err != nil {
		log.Println("获取secretKey失败,err: ", err)
		return err
	}
	// 数据加密
	key := secretKey.(string)
	data, err = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	if err != nil {
		log.Println("数据加密出错,err: ", err)
		return err
	}
	// 数据压缩
	data, err = utils.Zip(data)
	if err != nil {
		log.Println("数据压缩出错,err: ", err)
		return err
	}

	// 发送数据
	err = c.wsConn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Println("数据发送出错,err: ", err)
		return err
	}

	return nil
}

func (c *ClientConn) SetOnPush(hook func(conn *ClientConn, body *RspBody)) {
	c.onPush = hook
}

func (c *ClientConn) Send(name string, msg interface{}) (*RspBody, error) {
	// 请求转发给登录服务
	c.Seq += 1
	seq := c.Seq
	sc := NewSyncCtx()
	c.syncCtxLock.Lock()
	c.syncCtxMap[seq] = sc
	c.syncCtxLock.Unlock()
	req := &ReqBody{
		Seq:  seq,
		Name: name,
		Msg:  msg,
	}
	err := c.write(req)
	if err != nil {
		sc.cancel()
		return nil, err
	}

	rsp := sc.wait()
	c.syncCtxLock.Lock()
	delete(c.syncCtxMap, seq)
	c.syncCtxLock.Unlock()

	return rsp, nil
}
