package gelf

import (
	"net"
	logging "github.com/dkolbly/go-logging"
)

type UDPConn struct {
	raw *net.UDPConn
}


func DialUDP(hostport string) (Conn, error) {
	addr, err := net.ResolveUDPAddr("udp", hostport)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	c := &UDPConn{
		raw: conn,
	}
	return c, nil
}

func (c *UDPConn) Send(header *Common, additional interface{}) error {
	buf, err := Encode(header, additional)
	if err != nil {
		return err
	}
	// fire and forget
	c.raw.Write(buf)
	return nil
}

func (c *UDPConn) Close() error {
	sock := c.raw
	c.raw = nil
	return sock.Close()
}

// Log implements the logging.Backend interface
func (c *UDPConn) Log(level logging.Level, calldepth int, rec *logging.Record) error {
	return c.Send(backend(level, calldepth+1, rec))
}

