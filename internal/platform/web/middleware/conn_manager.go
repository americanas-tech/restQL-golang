package middleware

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"io"
	"net"
	"sync"
	"syscall"
	"time"
)

// ConnManager keeps track of active TCP connections
// and their corresponding contexts.
//
// When a connection is closed by the client it cancel
// its context.
type ConnManager struct {
	log          restql.Logger
	contextIndex sync.Map
}

// NewConnManager creates a connection manager
func NewConnManager(log restql.Logger) *ConnManager {
	return &ConnManager{log: log}
}

// ContextForConnection get the connections context, if it exists.
// Otherwise it create a context and start the connection watcher.
func (cm *ConnManager) ContextForConnection(conn net.Conn) context.Context {
	connCtx, found := cm.contextIndex.Load(conn)
	if !found {
		return cm.initializeConnContext(conn)
	}

	ctx, ok := connCtx.(context.Context)
	if !ok {
		return cm.initializeConnContext(conn)
	}

	return ctx
}

func (cm *ConnManager) initializeConnContext(conn net.Conn) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go cm.watchConn(conn, func() {
		cm.contextIndex.Delete(conn)
		cancel()
	})

	cm.contextIndex.Store(conn, ctx)

	return ctx
}

func (cm *ConnManager) watchConn(conn net.Conn, callback func()) {
	rc, err := conn.(syscall.Conn).SyscallConn()
	if err != nil {
		cm.log.Error("failed to cast connection to syscall interface", err)
		callback()
		return
	}

	ticker := time.NewTicker(10 * time.Millisecond)

	var sysErr error = nil
	for {
		select {
		case <-ticker.C:
			err = rc.Read(func(fd uintptr) bool {
				var buf = []byte{0}
				n, _, err := syscall.Recvfrom(int(fd), buf, syscall.MSG_PEEK|syscall.MSG_DONTWAIT)
				switch {
				case n == 0 && err == nil:
					sysErr = io.EOF
				case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
					sysErr = nil
				default:
					sysErr = err
				}
				return true
			})
			if err != nil {
				callback()
				return
			}

			if sysErr != nil {
				callback()
				return
			}
		}
	}
}
