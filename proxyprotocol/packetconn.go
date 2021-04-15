package proxyprotocol

import (
	"net"
	"sync"
	"time"
)

// Listener wraps a net.Listener automatically wrapping new connections with PROXY protocol support.
type PacketConn struct {
	net.PacketConn

	filter []Rule
	t      time.Duration

	mx sync.RWMutex
}


