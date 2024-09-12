package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFuncCandidates(t *testing.T) {
	test := "start_new_thread(hello, (Fastnode,)"

	funcs := findFuncCandidates(test)

	assert.Equal(t, 1, len(funcs))
}
