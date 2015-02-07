// Package gelf implements TCP and UDP connectivity in graylog format (GELF)
package gelf

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	logging "github.com/dkolbly/go-logging"
)

type Time time.Time

type Common struct {
	Version      string `json:"version"`
	Host         string `json:"host"`
	ShortMessage string `json:"short_message"`
	FullMessage  string `json:"full_message,omitempty"`
	Timestamp    Time   `json:"timestamp"`
	Level        int    `json:"level,omitempty"`
	Facility     string `json:"facility,omitempty"`
}

func (t Time) IsZero() bool {
	raw := time.Time(t)
	return raw.IsZero()
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("0"), nil
	}
	raw := time.Time(t)
	ns := raw.UnixNano()
	sec := ns / 1000000000
	msec := (ns / 1000000) % 1000
	return []byte(fmt.Sprintf("%d.%03d", sec, msec)), nil
}

type Conn interface {
	Send(header *Common, additional interface{}) error
	Close() error
	Log(level logging.Level, calldepth int, rec *logging.Record) error
}

var ErrInvalidScheme = errors.New("invalid GELF server scheme, expected gelf+udp or gelf+tcp")

func Dial(server string) (Conn, error) {
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	h := u.Host
	if !strings.Contains(h, ":") {
		h = h + ":12201"
	}
	if u.Scheme == "gelf+udp" {
		return DialUDP(h)
	} else if u.Scheme == "gelf+tcp" {
		return DialTCP(h)
	} else {
		return nil, ErrInvalidScheme
	}
}
