Go GELF Library
===============

This package provides support for talking GELF over UDP or TCP,
with support for annotations from the `go-logging` package.


```
package main

import (
	"github.com/dkolbly/gelf"
	logging "github.com/dkolbly/go-logging"
	"time"
)

var log = logging.MustGetLogger("test")

func main() {
	g, err := gelf.Dial("gelf+tcp://graylog.internal/")
	if err != nil {
		panic(err)
	}

	logging.SetBackend(g)

	log.Info("Test message")
	// if we exit too soon, it'll be before we can send the message!
	<-time.After(time.Second)
}
```
