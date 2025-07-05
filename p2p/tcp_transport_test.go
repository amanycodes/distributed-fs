package p2p

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	opts := TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.ListenAddr, listenAddr)

}
