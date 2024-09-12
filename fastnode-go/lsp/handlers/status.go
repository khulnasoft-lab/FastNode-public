package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lsp/types"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/enginestatus"
)

// Status gets the state of the Fastnode Engine for the given file.
func (h *Handlers) Status(params types.FastnodeStatusParams) (enginestatus.Response, error) {
	filepath, err := filepathFromURI(params.URI)
	if err != nil {
		return enginestatus.Response{}, err
	}
	statusPath, err := buildURL(statusURL, map[string]string{"filename": filepath})
	if err != nil {
		return enginestatus.Response{}, err
	}
	res, err := http.Get(statusPath)
	if err != nil {
		return enginestatus.Response{}, err
	}

	fastnodeResponse := enginestatus.Response{}
	err = json.NewDecoder(res.Body).Decode(&fastnodeResponse)
	if err != nil {
		return enginestatus.Response{}, err
	}
	return fastnodeResponse, nil
}
