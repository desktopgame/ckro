package model

type BufferSegment struct {
	buffer Buffer
	r      Range
}

func (bs BufferSegment) GetLine(lineIndex int) string {
	type Span struct {
		StartColumn int
		EndColumn   int
	}
	spans := []Span{}

	if bs.r.StartPosition.Row == bs.r.EndPosition.Row {
		spans = append(spans, Span{
			StartColumn: bs.r.StartPosition.Column,
			EndColumn:   bs.r.EndPosition.Column,
		})
	} else {
		spans = append(spans, Span{
			StartColumn: bs.r.StartPosition.Column,
			EndColumn:   len(bs.buffer.GetLineAt(bs.r.StartPosition.Row).GetContent()),
		})

		for i := bs.r.StartPosition.Row + 1; i < bs.r.EndPosition.Row; i++ {
			spans = append(spans, Span{
				StartColumn: 0,
				EndColumn:   len(bs.buffer.GetLineAt(i).GetContent()),
			})
		}

		spans = append(spans, Span{
			StartColumn: 0,
			EndColumn:   bs.r.EndPosition.Column,
		})
	}

	span := spans[lineIndex]
	line := bs.buffer.GetLineAt(bs.r.StartPosition.Row + lineIndex)
	return line.GetContent()[span.StartColumn:span.EndColumn]
}

func (bs BufferSegment) GetLineCount() int {
	startRow := bs.r.StartPosition.Row
	endRow := bs.r.EndPosition.Row
	return (endRow - startRow) + 1
}
