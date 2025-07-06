package main

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTranformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some bytes of jpg"))
	if err := s.WriteStream("123", "myspecialPicture", data); err != nil {
		t.Error(err)
	}
}
