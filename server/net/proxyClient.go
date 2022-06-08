package net

import (
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ProxyClient struct {
	proxy string
	conn  *ClientConn
}

func NewProxyClient(proxy string) *ProxyClient {
	return &ProxyClient{
		proxy: proxy,
	}
}

// 连接websocket服务端
func (c *ProxyClient) Connect() error {
	var dialer = websocket.Dialer{
		Subprotocols:     []string{"p1", "p2"},
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
	}

	ws, _, err := dialer.Dial(c.proxy, nil)
	if err != nil {
		log.Panicln("连接异常，err", err)
		return err
	}
	c.conn = NewClientConn(ws)
	if !c.conn.Start() {
		return errors.New("握手失败")
	}
	return nil
}

func (c *ProxyClient) SetProperty(key string, data interface{}) {
	if c.conn != nil {
		c.conn.SetProperty(key, data)
	}
}

func (c *ProxyClient) SetOnPush(hook func(conn *ClientConn, body *RspBody)) {
	if c.conn != nil {
		c.conn.SetOnPush(hook)
	}
}

func (c *ProxyClient) Send(name string, msg interface{}) (*RspBody, error) {
	if c.conn != nil {
		return c.conn.Send(name, msg)
	}
	return nil, errors.New("未找到连接")
}
