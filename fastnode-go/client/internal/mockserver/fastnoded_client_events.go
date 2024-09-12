package mockserver

import (
	"bytes"
	"encoding/json"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/component"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/stringindex"
)

// PostEvent sends an event to the client. It returns the http response status code and an error. -1 is returned for the status if an error ocurred
func (t *FastnodedClient) PostEvent(event component.EditorEvent) (int, error) {
	bodyBytes, err := json.Marshal(event)
	if err != nil {
		return -1, err
	}

	resp, err := t.Post("/clientapi/editor/event", bytes.NewReader(bodyBytes))
	if err != nil {
		return -1, err
	}

	return resp.StatusCode, nil
}

// PostEventData posts an editor event to the client. It returns the http response status code and an error. -1 is returned for the status if an error ocurred
func (t *FastnodedClient) PostEventData(source, filename, text, action string, startRune, endRune int64, enc stringindex.OffsetEncoding) (int, error) {
	return t.PostEvent(component.EditorEvent{
		Source:     source,
		Filename:   filename,
		Text:       text,
		Action:     action,
		Selections: []*component.Selection{{Start: startRune, End: endRune, Encoding: enc}},
	})
}

// PostFocusEvent sends a new event of type focus to fastnoded
func (t *FastnodedClient) PostFocusEvent(editor, filename, text string, runeOffset int64) (int, error) {
	return t.PostEventData(editor, filename, text, "focus", runeOffset, runeOffset, stringindex.UTF32)
}

// PostSelectionEvent sends a new event of type selection to fastnoded
func (t *FastnodedClient) PostSelectionEvent(editor, filename, text string, runeOffset int64) (int, error) {
	return t.PostEventData(editor, filename, text, "selection", runeOffset, runeOffset, stringindex.UTF32)
}

// PostEditEvent sends a new event of type edit to fastnoded
func (t *FastnodedClient) PostEditEvent(editor, filename, text string, runeOffset int64) (int, error) {
	return t.PostEventData(editor, filename, text, "edit", runeOffset, runeOffset, stringindex.UTF32)
}
