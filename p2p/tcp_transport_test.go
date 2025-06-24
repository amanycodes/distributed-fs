package p2p

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestNewTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	opts := TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       GOBDecoder{},
	}
	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.ListenAddr, listenAddr)

	err := tr.ListenAndAccept()
	assert.NilError(t, err)

	// Give the listener a moment to start, then close it to end the test
	time.Sleep(100 * time.Millisecond)
	if tr.listener != nil {
		tr.listener.Close()
	}
}
