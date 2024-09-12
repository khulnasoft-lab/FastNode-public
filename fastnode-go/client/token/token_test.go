package token

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Token(t *testing.T) {
	token := NewToken()
	headers := make(http.Header)

	//empty data at first
	token.AddToHeader(headers)
	assert.Equal(t, 0, len(headers))

	//extract from headers and add to new headers to make sure it's working
	headers.Add(tokenHeaderKey, "fastnode-key")
	headers.Add(tokenDataHeaderKey, "fastnode-data")
	token.UpdateFromHeader(headers)

	newHeaders := make(http.Header)
	token.AddToHeader(newHeaders)
	assert.Equal(t, 2, len(newHeaders))
	assert.Equal(t, "fastnode-key", newHeaders.Get(tokenHeaderKey))
	assert.Equal(t, "fastnode-data", newHeaders.Get(tokenDataHeaderKey))

	//clear
	newHeaders = make(http.Header)
	token.Clear()
	token.AddToHeader(newHeaders)
	assert.Equal(t, 0, len(newHeaders))
}
