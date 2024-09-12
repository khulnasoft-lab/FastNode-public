package codebase

import (
	"context"
	"sync"

	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

type terminator struct {
	terminated bool
	cancel     context.CancelFunc
	m          *sync.Mutex
}

func newTerminator() terminator {
	return terminator{
		m: new(sync.Mutex),
	}
}

func (t *terminator) terminate() {
	t.m.Lock()
	defer t.m.Unlock()

	t.terminated = true
	if t.cancel != nil {
		t.cancel()
	}
}

func (t terminator) wasTerminated() bool {
	t.m.Lock()
	defer t.m.Unlock()
	return t.terminated
}

func (t *terminator) closureWithCancel(fn func(ctx fastnodectx.Context) error) (func() error, fastnodectx.CancelFunc) {
	t.m.Lock()
	defer t.m.Unlock()

	var closure func() error
	closure, t.cancel = fastnodectx.Background().ClosureWithCancel(fn)
	return closure, t.cancel
}
