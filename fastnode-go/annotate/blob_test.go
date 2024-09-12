package annotate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBlob_Line(t *testing.T) {
	s := "[[FASTNODE[[LINE 15]]FASTNODE]]\n"
	expected := []blob{
		blob{Type: lineBlob, Line: 15},
	}
	actual, err := parseBlobs(s)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestParseBlob_Emit(t *testing.T) {
	s := "[[FASTNODE[[SHOW 123]]FASTNODE]]\n"
	expected := []blob{
		blob{Type: emitBlob, Content: "123"},
	}
	actual, err := parseBlobs(s)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestParseBlob_Output(t *testing.T) {
	s := "abc"
	expected := []blob{
		blob{Type: outputBlob, Content: "abc"},
	}
	actual, err := parseBlobs(s)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestParseBlob_LineAndOutput(t *testing.T) {
	s := "abc[[FASTNODE[[LINE 5]]FASTNODE]]\ndef"
	expected := []blob{
		blob{Type: outputBlob, Content: "abc"},
		blob{Type: lineBlob, Line: 5},
		blob{Type: outputBlob, Content: "def"},
	}
	actual, err := parseBlobs(s)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestParseBlob_Three(t *testing.T) {
	s := `[[FASTNODE[[LINE 1]]FASTNODE]]
abc
[[FASTNODE[[SHOW {"foo": 0}]]FASTNODE]]
`
	expected := []blob{
		blob{Type: lineBlob, Line: 1},
		blob{Type: outputBlob, Content: "abc\n"},
		blob{Type: emitBlob, Content: `{"foo": 0}`},
	}
	actual, err := parseBlobs(s)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestParseBlob_Mixed(t *testing.T) {
	s := `[[FASTNODE[[LINE 1]]FASTNODE]]
abc
[[FASTNODE[[SHOW {"foo": 0}]]FASTNODE]]
[[FASTNODE[[LINE 5]]FASTNODE]]
def
[[FASTNODE[[LINE 8]]FASTNODE]]
[[FASTNODE[[LINE 9]]FASTNODE]]
[[FASTNODE[[SHOW {"bar": 0}]]FASTNODE]]
`
	expected := []blob{
		blob{Type: lineBlob, Line: 1},
		blob{Type: outputBlob, Content: "abc\n"},
		blob{Type: emitBlob, Content: `{"foo": 0}`},
		blob{Type: lineBlob, Line: 5},
		blob{Type: outputBlob, Content: "def\n"},
		blob{Type: lineBlob, Line: 8},
		blob{Type: lineBlob, Line: 9},
		blob{Type: emitBlob, Content: `{"bar": 0}`},
	}
	actual, err := parseBlobs(s)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestParseBlob_Empty(t *testing.T) {
	actual, err := parseBlobs("")
	require.NoError(t, err)
	assert.Len(t, actual, 0)
}
