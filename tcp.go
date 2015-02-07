package gelf

import (
	"errors"
	"net"
	"sync/atomic"
	"time"
	logging "github.com/dkolbly/go-logging"
)

type TCPConn struct {
	addr   string
	txbuf  chan []byte
	lost   uint64
	queued uint64
}

func DialTCP(hostport string) (Conn, error) {
	c := &TCPConn{
		addr:  hostport,
		txbuf: make(chan []byte, 1000),
	}
	go c.txflush()
	return c, nil
}

func (c *TCPConn) Send(header *Common, additional interface{}) error {
	buf, err := Encode(header, additional)
	if err != nil {
		return err
	}
	select {
	case c.txbuf <- buf:
		atomic.AddUint64(&c.queued, 1)
		return nil
	default:
		atomic.AddUint64(&c.lost, 1)
		return ErrLostMessage
	}
}

var ErrLostMessage = errors.New("buffer overflow, lost message")

func (c *TCPConn) reconnect() *net.TCPConn {
	for {
		addr, err := net.ResolveTCPAddr("tcp", c.addr)
		if err == nil {
			raw, err := net.DialTCP("tcp", nil, addr)
			if err == nil {
				return raw
			}
		}
		<-time.After(5 * time.Second)
	}
}

func (c *TCPConn) txflush() {
	var raw *net.TCPConn

	for item := range c.txbuf {
		if raw == nil {
			raw = c.reconnect()
		}
		// append a NUL byte
		item = append(item, 0x00)
		_, err := raw.Write(item)
		if err != nil {
			raw.Close()
			raw = nil
		}
	}
	if raw != nil {
		raw.Close()
	}
}

func (c *TCPConn) Close() error {
	close(c.txbuf)
	return nil
}

// Log implements the logging.Backend interface
func (c *TCPConn) Log(level logging.Level, calldepth int, rec *logging.Record) error {
	return c.Send(backend(level, calldepth+1, rec))
}
