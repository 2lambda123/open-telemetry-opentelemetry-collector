// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package confignet

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddrConfigTimeout(t *testing.T) {
	nac := &AddrConfig{
		Endpoint:  "localhost:0",
		TransportType: TransportTypeTCP,
		DialerConfig: DialerConfig{
			Timeout: -1 * time.Second,
		},
	}
	_, err := nac.Dial(context.Background())
	assert.Error(t, err)
	var netErr net.Error
	if errors.As(err, &netErr) {
		assert.True(t, netErr.Timeout())
	} else {
		assert.Fail(t, "error should be a net.Error")
	}
}

func TestTCPAddrConfigTimeout(t *testing.T) {
	nac := &TCPAddrConfig{
		Endpoint: "localhost:0",
		DialerConfig: DialerConfig{
			Timeout: -1 * time.Second,
		},
	}
	_, err := nac.Dial(context.Background())
	assert.Error(t, err)
	var netErr net.Error
	if errors.As(err, &netErr) {
		assert.True(t, netErr.Timeout())
	} else {
		assert.Fail(t, "error should be a net.Error")
	}
}

func TestAddrConfig(t *testing.T) {
	nas := &AddrConfig{
		Endpoint:  "localhost:0",
		TransportType: TransportTypeTCP,
	}
	ln, err := nas.Listen(context.Background())
	assert.NoError(t, err)
	done := make(chan bool, 1)

	go func() {
		conn, errGo := ln.Accept()
		assert.NoError(t, errGo)
		buf := make([]byte, 10)
		var numChr int
		numChr, errGo = conn.Read(buf)
		assert.NoError(t, errGo)
		assert.Equal(t, "test", string(buf[:numChr]))
		assert.NoError(t, conn.Close())
		done <- true
	}()

	nac := &AddrConfig{
		Endpoint:  ln.Addr().String(),
		Transport: "tcp",
	}
	var conn net.Conn
	conn, err = nac.Dial(context.Background())
	assert.NoError(t, err)
	_, err = conn.Write([]byte("test"))
	assert.NoError(t, err)
	assert.NoError(t, conn.Close())
	<-done
	assert.NoError(t, ln.Close())
}

func Test_NetAddr_Validate(t *testing.T) {
	na := &AddrConfig{
		TransportType: TransportTypeTCP,
	}
	assert.NoError(t, na.Validate())

	na = &AddrConfig{
		TransportType: transportTypeEmpty,
	}
	assert.Error(t, na.Validate())

	na = &AddrConfig{
		TransportType: "random string",
	}
	assert.Error(t, na.Validate())
}

func TestTCPAddrConfig(t *testing.T) {
	nas := &TCPAddrConfig{
		Endpoint: "localhost:0",
	}
	ln, err := nas.Listen(context.Background())
	assert.NoError(t, err)
	done := make(chan bool, 1)

	go func() {
		conn, errGo := ln.Accept()
		assert.NoError(t, errGo)
		buf := make([]byte, 10)
		var numChr int
		numChr, errGo = conn.Read(buf)
		assert.NoError(t, errGo)
		assert.Equal(t, "test", string(buf[:numChr]))
		assert.NoError(t, conn.Close())
		done <- true
	}()

	nac := &TCPAddrConfig{
		Endpoint: ln.Addr().String(),
	}
	var conn net.Conn
	conn, err = nac.Dial(context.Background())
	assert.NoError(t, err)
	_, err = conn.Write([]byte("test"))
	assert.NoError(t, err)
	assert.NoError(t, conn.Close())
	<-done
	assert.NoError(t, ln.Close())
}
