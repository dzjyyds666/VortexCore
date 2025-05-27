package vortex

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

// 根据请求的前几个字节做协议识别，并并发连接
// 分发器
type Dispatcher struct {
	ctx  context.Context
	conn net.Conn
	peek bytes.Buffer
	r    io.Reader
}

func NewDispatcher(ctx context.Context, conn net.Conn) *Dispatcher {
	return &Dispatcher{
		ctx:  ctx,
		conn: conn,
	}
}

func (d *Dispatcher) Parse() (string, error) {
	peekLen := 24 // Enough for HTTP/2 preface
	buf := make([]byte, peekLen)
	n, err := d.conn.Read(buf)
	if err != nil {
		return "", err
	}

	d.peek.Write(buf[:n]) // 保存起来

	if isHttp1(d.peek.Bytes()) {
		// 重新构造一个 reader 来解析完整的 HTTP 请求头
		fullReader := io.MultiReader(bytes.NewReader(d.peek.Bytes()), d.conn)
		req, parseErr := http.ReadRequest(bufio.NewReader(fullReader))
		if parseErr == nil && req.Header.Get("Upgrade") == "websocket" {
			return webSocket, nil
		}
		return http1, nil
	} else if isHttp2Preface(d.peek.Bytes()) {
		return http2, nil
	} else {
		return "unknown", nil
	}
}

func (d *Dispatcher) Response(resp []byte) error {
	_, err := d.conn.Write(resp) // 发送响应
	if nil != err {
		return err
	}
	return nil
}

func (d *Dispatcher) Read(p []byte) (int, error) {
	if d.r == nil {
		// 如果未调用 Parse()，直接读原始连接
		return d.conn.Read(p)
	}
	return d.r.Read(p)
}

// 代理其余 Conn 方法
func (d *Dispatcher) Write(p []byte) (int, error)        { return d.conn.Write(p) }
func (d *Dispatcher) Close() error                       { return d.conn.Close() }
func (d *Dispatcher) LocalAddr() net.Addr                { return d.conn.LocalAddr() }
func (d *Dispatcher) RemoteAddr() net.Addr               { return d.conn.RemoteAddr() }
func (d *Dispatcher) SetDeadline(t time.Time) error      { return d.conn.SetDeadline(t) }
func (d *Dispatcher) SetReadDeadline(t time.Time) error  { return d.conn.SetReadDeadline(t) }
func (d *Dispatcher) SetWriteDeadline(t time.Time) error { return d.conn.SetWriteDeadline(t) }

func (d *Dispatcher) GetReadBuffer() io.Reader {
	if d.r == nil {
		// 如果未调用 Parse()，直接返回原始连接
		return d.conn
	}
	return d.r
}

func isHttp2Preface(buf []byte) bool {
	return bytes.Equal(buf, []byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"))
}

func isHttp1(buf []byte) bool {
	return bytes.HasPrefix(buf, []byte("GET ")) ||
		bytes.HasPrefix(buf, []byte("POST")) ||
		bytes.HasPrefix(buf, []byte("HEAD")) ||
		bytes.HasPrefix(buf, []byte("PUT ")) ||
		bytes.HasPrefix(buf, []byte("DELE"))
}

func isWebSocket(buf []byte) bool {
	return bytes.HasPrefix(buf, []byte("GET ")) &&
		bytes.Contains(buf, []byte("Upgrade: websocket"))
}
