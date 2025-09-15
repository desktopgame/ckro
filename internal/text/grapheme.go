package text

import (
	"slices"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

func GraphemeClusters(s string) []string {
	result := []string{}
	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		cluster := gr.Str()
		result = append(result, cluster)
	}
	return result
}

func GraphemeLength(s string) int {
	gr := uniseg.NewGraphemes(s)
	count := 0
	for gr.Next() {
		count++
	}
	return count
}

func GraphemeToCodepointPos(s string, graphemePos int) int {
	if graphemePos <= 0 {
		return 0
	}

	clusters := GraphemeClusters(s)
	if graphemePos >= len(clusters) {
		return len(s)
	}
	return len(strings.Join(clusters[:graphemePos], ""))
}

func CodepointToGraphemePos(s string, codepointPos int) int {
	if codepointPos <= 0 {
		return 0
	}
	if codepointPos >= len(s) {
		return GraphemeLength(s)
	}

	clusters := GraphemeClusters(s)
	currentPos := 0

	for i, cluster := range clusters {
		if currentPos >= codepointPos {
			return i
		}
		currentPos += len(cluster)
	}
	return len(clusters)
}

func GraphemeSubString(s string, startGrapheme int, endGrapheme int) string {
	clusters := GraphemeClusters(s)
	if endGrapheme == -1 {
		return strings.Join(clusters[startGrapheme:], "")
	} else {
		return strings.Join(clusters[startGrapheme:endGrapheme], "")
	}
}

func GraphemeRemove(s string, startGrapheme int, lengthGrapheme int) string {
	clusters := GraphemeClusters(s)
	if startGrapheme < 0 || startGrapheme >= len(clusters) {
		return s
	}
	endPos := min(startGrapheme+lengthGrapheme, len(clusters))
	newClusters := slices.Concat(clusters[:startGrapheme], clusters[endPos:])
	return strings.Join(newClusters, "")
}

func GraphemeInsert(s string, graphemePos int, insertString string) string {
	clusters := GraphemeClusters(s)
	if graphemePos <= 0 {
		return insertString + s
	} else if graphemePos >= len(clusters) {
		return s + insertString
	} else {
		before := strings.Join(clusters[:graphemePos], "")
		after := strings.Join(clusters[graphemePos:], "")
		return before + insertString + after
	}
}

func DisplayWidth(s string) int {
	return runewidth.StringWidth(s)
}

func DisplayPos(line string, column int) int {
	if column <= 0 {
		return 0
	}

	if column >= GraphemeLength(line) {
		return runewidth.StringWidth(line)
	}

	targetStr := GraphemeSubString(line, 0, column)
	return runewidth.StringWidth(targetStr)
}

func DisplayRunesAt(line string, screenX int) (rune, []rune) {
	currentX := 0
	gr := uniseg.NewGraphemes(line)
	for gr.Next() {
		if currentX == screenX {
			cluster := gr.Str()
			runes := []rune(cluster)
			if len(runes) > 0 {
				mainRune := runes[0]
				var combining []rune
				if len(runes) > 1 {
					combining = runes[1:]
				}
				return mainRune, combining
			}
		}
		cluster := gr.Str()
		if len(cluster) > 0 {
			mainRune := []rune(cluster)[0]
			width := runewidth.RuneWidth(mainRune)
			if width == 2 {
				currentX += 2
			} else {
				currentX++
			}
		} else {
			currentX++
		}
	}
	return ' ', nil
}
