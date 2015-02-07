package gelf

import (
	"time"
	"testing"
)

var t0 Time

func init() {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		panic(err)
	}
	// October 9th of this 2015 is even cooler
	t0 = Time(time.Date(2015, 1, 25, 15, 43, 42, 222000000, loc))
}
	
func TestEncodeTriv(t *testing.T) {
	a, err := Encode(&Common{Host:"test", Timestamp:t0}, nil)
	if err != nil {
		t.Fatal(err)
	}
	x := `{"version":"1.1","host":"test","short_message":"A log message","timestamp":1422222222.222}`
	if string(a) != x {
		t.Fatalf("Got %q, expected %q", a, x)
	}
}


func TestEncodeExtra(t *testing.T) {
	extra := map[string]interface{}{"_foo": 123}
	a, err := Encode(&Common{Host:"test", Timestamp:t0}, extra)
	if err != nil {
		t.Fatal(err)
	}
	x := `{"version":"1.1","host":"test","short_message":"A log message","timestamp":1422222222.222,"_foo":123}`
	if string(a) != x {
		t.Fatalf("Got %q, expected %q", a, x)
	}
	
}

