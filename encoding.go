package gelf

import (
	"os"
	"time"
	"bytes"
	"encoding/json"
)

func Encode(header *Common, additional interface{}) ([]byte, error) {
	h := *header
	if h.Version == "" {
		h.Version = "1.1"
	}
	if h.Host == "" {
		hostname, err := os.Hostname()
		if err == nil {
			h.Host = hostname
		} else {
			h.Host = "localhost"
		}
	}
	if h.ShortMessage == "" {
		h.ShortMessage = "A log message"
	}
	if h.Timestamp.IsZero() {
		h.Timestamp = Time(time.Now())
	}

	buf, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	// if there are no additional fields, then we're done already
	if additional == nil {
		return buf, nil
	}
	extra, err := json.Marshal(additional)
	if err != nil {
		return nil, err
	}
	if len(extra) == 2 || extra[0] != '{' {
		// another way we might get no additional fields
		return buf, nil
	}

	// add on the additional stuff; create a buffer starting with
	// the header data but stripping off the trailing '}'
	accum := bytes.NewBuffer(buf[:len(buf)-1])
	accum.WriteByte(',');
	accum.Write(extra[1:])
	return accum.Bytes(), nil
}
