package encode

import "github.com/khulnasoft-lab/fastnode/fastnode-golib/lexicalv0/text"

// Lex ...
func Lex(buf string) []string {
	return text.SplitWithOpts(buf, true)
}
