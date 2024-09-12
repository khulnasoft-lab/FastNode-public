package text

type token int

const (
	sof token = iota
	eof
	text
	unknownToken // typically comes from words that are too large
)

// String ...
func (t token) String() string {
	switch t {
	case sof:
		return "fastnode-textlex184-SOF"
	case eof:
		return "fastnode-textlex184-EOF"
	case text:
		return "fastnode-textlex184-TEXT"
	default:
		return "fastnode-textlex184-UNK"
	}
}

var allTokens = []token{
	sof,
	eof,
	text,
	unknownToken,
}
