package p2p

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	lr := NewTCPTransport(listenAddr)

	assert.Equal(t, lr.listenAddress, listenAddr)

	assert.NilError(t, lr.ListenAndAccept())
}
