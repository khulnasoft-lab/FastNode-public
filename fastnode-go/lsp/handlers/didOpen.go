package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/component"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lsp/types"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/complete/data"
)

// DidOpen handles an LSP textDocument/didOpen call
func (h *Handlers) DidOpen(params types.DidOpenTextDocumentParams) error {
	filepath, err := filepathFromURI(params.TextDocument.URI)
	if err != nil {
		return err
	}

	text := params.TextDocument.Text
	h.files[filepath] = text

	event := component.EditorEvent{
		Source:   string(data.JupyterEditor),
		Action:   "focus",
		Filename: filepath,
		Text:     text,
		// TODO: Selections
	}
	buf, _ := json.Marshal(event)
	_, err = http.Post(eventURL, contentType, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}
