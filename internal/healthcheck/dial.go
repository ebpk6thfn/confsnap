package healthcheck

import (
	"net"
	"time"
)

// dialTCP opens a raw TCP connection to addr within the given timeout.
// It is a thin wrapper around net.DialTimeout so that tests can replace it.
var dialTCP = func(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("tcp", addr, timeout)
}
