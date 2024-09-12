# Fastnode LSP

`fastnode-lsp` is an intermediary between editor-clients that speak [Language Server Protocol](https://microsoft.github.io/language-server-protocol/), and the Fastnode Engine. It maintains an LSP session, and translates between LSP requests and Fastnode API requests.

## Setup with JupyterLab
First, make sure you've installed [JupyterLab](https://github.com/jupyterlab/jupyterlab), [`jupyterlab-fastnode`](https://github.com/khulnasoft-lab/jupyterlab-fastnode#installation), and the Fastnode Engine.

Then, to build `fastnode-lsp` run the following:
```bash
go install github.com/khulnasoft-lab/fastnode/fastnode-go/lsp/cmds/fastnode-lsp
```
This will build the `fastnode-lsp` binary and then place it in `$GOPATH/bin`.

To make `jupyterlab-fastnode` aware of your build of `fastnode-lsp`, move to your Jupyter config folder (usually `$HOME/.jupyter`) and create a file called `jupyter_notebook_config.json`, with the following contents:
```json
{
  "LanguageServerManager": {
      "language_servers": {
          "fastnodels": {
              "argv": [
                  "YOUR_FASTNODE_LSP_LOCATION"
              ],
              "languages": [
                  "python"
              ],
              "version": 2
          }
      }
  }
}
```

If you built `fastnode-lsp` using the `go install` instruction above, then `YOUR_FASTNODE_LSP_LOCATION` will be your `GOPATH` plus `bin/fastnode-lsp`, and this needs to be provided as an absolute path (i.e. don't just put `"$GOPATH/bin/fastnode-lsp"` in the config).

## Usage
After installation, make sure that you have the Fastnode Engine running on your machine, and then start JupyterLab by running:
```bash
jupyter lab
```

