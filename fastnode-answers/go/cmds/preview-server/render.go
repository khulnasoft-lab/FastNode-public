package main

import (
	"encoding/json"
	"net/http"

	"github.com/khulnasoft-lab/fastnode/fastnode-answers/go/execution"
	"github.com/khulnasoft-lab/fastnode/fastnode-answers/go/render"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lang/python/pythonresource"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/fastnodectx"
)

func doRender(w http.ResponseWriter, sandbox execution.Manager,
	resourceMgr pythonresource.Manager, src []byte) {
	raw, err := render.ParseRaw(src)
	var out render.Rendered
	if err != nil {
		var i render.CodeBlockItem
		i.Output = &execution.Output{
			Type:  "text",
			Title: "post parse error",
			Data:  err.Error(),
		}
		out.Content = append(out.Content, render.Block{CodeBlock: []render.CodeBlockItem{i}})
	} else {
		out, _ = render.Render(fastnodectx.TODO(), sandbox, resourceMgr, raw)
	}
	res, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)
	return
}

func rootAssetsAndApp() (http.Handler, http.Handler) {
	index := MustAsset("assets/index.html")
	app := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(index)
	}
	return http.FileServer(assetFS()), http.HandlerFunc(app)
}
