package tui_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

func printCharacter(a []tcell.SimCell, cluster string, s tcell.Style) []tcell.SimCell {
	runes := []rune(cluster)

	if runewidth.RuneWidth(runes[0]) == 2 {
		a = append(a, tcell.SimCell{
			Bytes: []byte(cluster),
			Runes: []rune(cluster),
			Style: s,
		})
		a = append(a, tcell.SimCell{
			Bytes: nil,
			Runes: nil,
			Style: tcell.StyleDefault,
		})
	} else {
		a = append(a, tcell.SimCell{
			Bytes: []byte(cluster),
			Runes: []rune(cluster),
			Style: s,
		})
	}
	return a
}

func testScenario(t *testing.T, scenarioFile string, width, height int) {
	tb := newPlainTextBox(width, height)
	tb.ShowCursor = true

	file, err := os.Open(scenarioFile)
	if err != nil {
		assert.Error(t, err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	assert.True(t, sc.Scan())
	assert.Equal(t, sc.Text(), "%%")

	for sc.Scan() {
		line := sc.Text()
		if line == "%%" {
			break
		}

		sepPos := strings.IndexByte(line, ' ')
		var instruction, args string

		if sepPos == -1 {
			instruction = line
			args = ""
		} else {
			instruction = line[:sepPos]
			args = line[sepPos+1:]
		}

		switch instruction {
		case "TYPE":
			assert.Equal(t, args[0], byte('"'))
			assert.Equal(t, args[len(args)-1], byte('"'))
			tb.InsertString(args[1 : len(args)-1])
		case "TYPELN":
			assert.Equal(t, args[0], byte('"'))
			assert.Equal(t, args[len(args)-1], byte('"'))
			tb.InsertString(fmt.Sprintf("%s\n", args[1:len(args)-1]))
		case "REMOVE_CHAR":
			tb.RemoveChar()
		case "REMOVE_SELECTION":
			tb.RemoveSelection()
		case "SELECTION_START":
			tb.SelectionStart()
		case "SELECTION_END":
			tb.SelectionEnd()
		case "MOVE_LEFT":
			tb.MoveLeft()
		case "MOVE_RIGHT":
			tb.MoveRight()
		case "MOVE_UP":
			tb.MoveUp()
		case "MOVE_DOWN":
			tb.MoveDown()
		}
	}

	expected := []tcell.SimCell{}
	cursorRe := regexp.MustCompile(`.*(<.+>).*`)
	rows := 0
	for sc.Scan() {
		line := sc.Text()
		if line == "%%" {
			break
		}
		if cursorRe.MatchString(line) {
			leftMarker := strings.IndexByte(line, '<')
			rightMarker := strings.IndexByte(line, '>')

			leftText := line[:leftMarker]
			inText := line[leftMarker+1 : rightMarker]
			rightText := line[rightMarker+1:]

			at := 0
			rev := tcell.StyleDefault.Reverse(true)
			for _, ltCluster := range text.GraphemeClusters(leftText) {
				expected = printCharacter(expected, ltCluster, tcell.StyleDefault)
				at++
			}
			for _, inCluster := range text.GraphemeClusters(inText) {
				expected = printCharacter(expected, inCluster, rev)
				at++
			}
			for _, rtCluster := range text.GraphemeClusters(rightText) {
				expected = printCharacter(expected, rtCluster, tcell.StyleDefault)
				at++
			}
			for i := at; i < tb.Width; i++ {
				expected = append(expected, tcell.SimCell{
					Bytes: []byte{' '},
					Runes: []rune{' '},
					Style: tcell.StyleDefault,
				})
			}
		} else {
			clusters := text.GraphemeClusters(line)
			for i := 0; i < tb.Width; i++ {
				if i >= len(clusters) {
					expected = append(expected, tcell.SimCell{
						Bytes: []byte{' '},
						Runes: []rune{' '},
						Style: tcell.StyleDefault,
					})
				} else {
					cluster := clusters[i]
					expected = printCharacter(expected, cluster, tcell.StyleDefault)
				}
			}
		}
		rows++
	}

	screen := tcell.NewSimulationScreen("UTF-8")
	screen.Init()
	screen.SetSize(tb.Width, tb.Height)
	defer screen.Fini()

	g := base.Graphics{}
	g.Init(screen)
	g.Resize(tb.Width, tb.Height)

	tb.Draw(&g)

	screen.Show()

	actual, _, _ := screen.GetContents()
	for i := 0; i < rows; i++ {
		for j := 0; j < tb.Width; j++ {
			a := actual[i*tb.Width+j]
			e := expected[i*tb.Width+j]

			assert.True(t, bytes.Equal(a.Bytes, e.Bytes), "file=%s expected=%s actual=%s row=%d col=%d", scenarioFile, string(e.Bytes), string(a.Bytes), i, j)

			_, _, aMask := a.Style.Decompose()
			_, _, eMask := e.Style.Decompose()
			assert.Equal(t, aMask, eMask, "file=%s row=%d col=%d", scenarioFile, i, j)
		}
	}
}

func TestAllScenario(t *testing.T) {
	entries, err := os.ReadDir("../../testdata/")
	if err != nil {
		assert.Error(t, err)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		file := filepath.Join("../../testdata/", entry.Name())
		testScenario(t, file, 20, 10)
	}
}
