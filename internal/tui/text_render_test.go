package tui_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui"
	"github.com/desktopgame/ckro/internal/tui/base"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

func printCharacter(a []tcell.SimCell, cluster string, s tcell.Style) ([]tcell.SimCell, int) {
	runes := []rune(cluster)
	added := 1

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
		added = 2
	} else {
		a = append(a, tcell.SimCell{
			Bytes: []byte(cluster),
			Runes: []rune(cluster),
			Style: s,
		})
	}
	return a, added
}

func extractQuote(s string, start int) (int, int) {
	st := strings.IndexByte(s[start:], '"')
	ed := strings.IndexByte(s[start+st+1:], '"')
	return start + st + 1, start + st + ed + 1
}

func testScenario(t *testing.T, scenarioFile string) string {
	file, err := os.Open(scenarioFile)
	if err != nil {
		assert.Error(t, err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	assert.True(t, sc.Scan())
	assert.Equal(t, sc.Text(), "%%")

	sc.Scan()
	width, _ := strconv.Atoi(sc.Text())

	sc.Scan()
	height, _ := strconv.Atoi(sc.Text())

	sc.Scan()
	textEngine := sc.Text()

	var tb *tui.TextBox
	switch textEngine {
	case "PlainText":
		tb = newPlainTextBox(width, height)
	case "StyledText":
		tb = newStyledTextBox(width, height)
	}
	tb.ShowCursor = true

	assert.True(t, sc.Scan())
	assert.Equal(t, sc.Text(), "%%")

	vars := map[string]string{}
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
		case "MOVE_TEXT_START":
			tb.MoveTextStart()
		case "MOVE_TEXT_END":
			tb.MoveTextEnd()
		case "MOVE_LINE_START":
			tb.MoveLineStart()
		case "MOVE_LINE_END":
			tb.MoveLineEnd()
		case "FIND_PREV":
			assert.Equal(t, args[0], byte('"'))
			assert.Equal(t, args[len(args)-1], byte('"'))
			assert.True(t, tb.FindPrev(args[1:len(args)-1]))
		case "FIND_NEXT":
			assert.Equal(t, args[0], byte('"'))
			assert.Equal(t, args[len(args)-1], byte('"'))
			assert.True(t, tb.FindNext(args[1:len(args)-1]))
		case "SUBMIT":
			assert.True(t, tb.Submit())
		case "BYTE_POS":
			bPosStr := strings.Split(args, " ")
			bPosRow, err := strconv.Atoi(bPosStr[0])
			if err != nil {
				assert.Error(t, err)
				return ""
			}
			bPosCol, err := strconv.Atoi(bPosStr[1])
			if err != nil {
				assert.Error(t, err)
				return ""
			}
			actualBytePos := tb.GetBytePosition()
			assert.Equal(t, bPosRow, actualBytePos.StartPosition.Row, "file=%s", scenarioFile)
			assert.Equal(t, bPosCol, actualBytePos.StartPosition.Column, "file=%s", scenarioFile)
		case "VAR":
			nameSt, nameEd := extractQuote(args, 0)
			valueSt, valueEd := extractQuote(args, nameEd+1)
			name := args[nameSt:nameEd]
			value := args[valueSt:valueEd]
			vars[name] = value
		}
	}

	expected := []tcell.SimCell{}
	cursorRe := regexp.MustCompile(`.*(<.+>).*`)
	writeRows := 0
	for sc.Scan() {
		line := sc.Text()
		if line == "%%" {
			break
		}
		for k, v := range vars {
			line = strings.ReplaceAll(line, k, v)
		}
		if cursorRe.MatchString(line) {
			leftMarker := strings.IndexByte(line, '<')
			rightMarker := strings.IndexByte(line, '>')

			leftText := line[:leftMarker]
			inText := line[leftMarker+1 : rightMarker]
			rightText := line[rightMarker+1:]

			at := 0
			rev := tcell.StyleDefault.Reverse(true)
			var added int
			for _, ltCluster := range text.GraphemeClusters(leftText) {
				expected, added = printCharacter(expected, ltCluster, tcell.StyleDefault)
				at += added
			}
			for _, inCluster := range text.GraphemeClusters(inText) {
				expected, added = printCharacter(expected, inCluster, rev)
				at += added
			}
			for _, rtCluster := range text.GraphemeClusters(rightText) {
				expected, added = printCharacter(expected, rtCluster, tcell.StyleDefault)
				at += added
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
			var added int
			for i := 0; i < tb.Width; {
				if i >= len(clusters) {
					expected = append(expected, tcell.SimCell{
						Bytes: []byte{' '},
						Runes: []rune{' '},
						Style: tcell.StyleDefault,
					})
					i++
				} else {
					cluster := clusters[i]
					expected, added = printCharacter(expected, cluster, tcell.StyleDefault)
					i += added
				}
			}
		}
		writeRows++
	}
	for writeRows < tb.Height {
		for i := 0; i < tb.Width; i++ {
			expected = append(expected, tcell.SimCell{
				Bytes: []byte{' '},
				Runes: []rune{' '},
				Style: tcell.StyleDefault,
			})
		}
		writeRows++
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
	sbuf := strings.Builder{}
	for i := 0; i < tb.Height; i++ {
		lineBuf := strings.Builder{}
		for j := 0; j < tb.Width; j++ {
			a := actual[i*tb.Width+j]
			e := expected[i*tb.Width+j]

			if a.Bytes != nil {
				lineBuf.WriteString(string(a.Bytes))
			}

			assert.True(t, bytes.Equal(a.Bytes, e.Bytes), "file=%s expected=%s actual=%s row=%d col=%d", scenarioFile, string(e.Bytes), string(a.Bytes), i, j)

			_, _, aMask := a.Style.Decompose()
			_, _, eMask := e.Style.Decompose()
			assert.Equal(t, (aMask&tcell.AttrReverse) > 0, (eMask&tcell.AttrReverse) > 0, "file=%s row=%d col=%d", scenarioFile, i, j)
		}
		sbuf.WriteString(strings.TrimRight(lineBuf.String(), " "))
		sbuf.WriteString("\n")
	}
	return sbuf.String()
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
		output := testScenario(t, file)

		outFilePath := filepath.Join("../../testdist/", entry.Name())
		outFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			assert.Error(t, err)
			continue
		}
		defer outFile.Close()
		outFile.WriteString(output)
	}
}
