package litemark_test

import (
	"strings"
	"testing"

	"github.com/desktopgame/ckro/internal/litemark"
	"github.com/stretchr/testify/assert"
)

func Test01(t *testing.T) {
	text := `
# H1
## H2
### H3
`
	text = strings.Trim(text, " \t\n")

	r := litemark.StringReader{
		Source: strings.Split(text, "\n"),
	}
	blocks := litemark.Parse(&r)

	h1 := blocks[0].(*litemark.Heading)
	assert.Equal(t, h1.Level, 1)

	h2 := blocks[1].(*litemark.Heading)
	assert.Equal(t, h2.Level, 2)

	h3 := blocks[2].(*litemark.Heading)
	assert.Equal(t, h3.Level, 3)
}

func Test02(t *testing.T) {
	lines := []string{
		"```lang",
		"L1",
		"L2",
		"L3",
		"```",
	}

	r := litemark.StringReader{
		Source: lines,
	}
	blocks := litemark.Parse(&r)

	cb := blocks[0].(*litemark.CodeBlock)
	assert.Equal(t, cb.LineIndex, 0)
	assert.Equal(t, cb.LineCount, 5)
}
