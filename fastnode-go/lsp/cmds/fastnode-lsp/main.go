package main

import (
	"log"
	"os"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lsp"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/lsp/jsonrpc2"
	"github.com/khulnasoft-lab/fastnode/fastnode-go/lsp/process"
)

func main() {
	log.Println("Fastnode-LSP Starting...")

	startFastnode()

	stdio := jsonrpc2.NewReaderWriterConnection(os.Stdin, os.Stdout)

	var lspServer *lsp.Server
	lspServer = lsp.New()

	rpcConn := jsonrpc2.NewRPCConnection(stdio, lspServer)

	err := rpcConn.Run()
	if err != nil {
		log.Println(err)
		log.Println("RPC connection closed.")
	}
}

func startFastnode() {
	// Attempt to start Fastnode if it's not running.
	isRunning, err := process.IsRunning(process.Name)
	if err != nil {
		log.Println("Could not check if Fastnode is running. Continuing initialization...")
		return
	}
	if !isRunning {
		log.Println("Fastnode not running! Attempting to start...")
		err = process.Start()
		if err != nil {
			log.Println("Could not autostart Fastnode:", err)
		}
	}
}
