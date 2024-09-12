package pigeon

import "github.com/khulnasoft-lab/fastnode/fastnode-go/lang/javascript/ast"

// Node is the raw representation of a
// javascript node from the generated parser.
type Node struct {
	Begin, Len int
	Type       ast.Type
	Children   interface{}
}
