package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/component"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lsp/types"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/complete/data"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/stringindex"
)

// Selection sends selection events to Fastnode
func (h *Handlers) Selection(params types.FastnodeSelectionParams) error {
	filepath, err := filepathFromURI(params.TextDocument.URI)
	if err != nil {
		return err
	}

	text := params.Text

	var selections []*component.Selection
	for _, p := range params.Positions {
		offset, err := utf8OffFromPos(text, p)
		if err != nil {
			return err
		}
		s := &component.Selection{Start: int64(offset), End: int64(offset), Encoding: stringindex.UTF8}
		selections = append(selections, s)
	}

	event := component.EditorEvent{
		Source:     string(data.JupyterEditor),
		Action:     "selection",
		Filename:   filepath,
		Text:       text,
		Selections: selections,
	}
	buf, _ := json.Marshal(event)
	_, err = http.Post(eventURL, contentType, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}
