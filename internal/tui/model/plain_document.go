package model

// PlainDocument is wrapper of Buffer.
// track a current cursor position.
type PlainDocument struct {
	buffer  Buffer
	version uint

	cache        []Element
	cacheVersion uint
}

// Init is initialize Buffer.
func (doc *PlainDocument) Init() {
	doc.Clear()
}

func (doc *PlainDocument) Read(r Range) Segment {
	return BufferSegment{
		buffer: doc.buffer,
		r:      r,
	}
}

func (doc *PlainDocument) doRender() []Element {
	elements := []Element{}

	// Fallback to plain text rendering
	for i := 0; i < doc.buffer.GetLineCount(); i++ {
		line := doc.buffer.GetLineAt(i)
		lineStr := line.GetContent()

		elements = append(elements, &PlainElement{
			Range: Range{
				StartPosition: Position{
					Row:    i,
					Column: 0,
				},
				EndPosition: Position{
					Row:    i,
					Column: len(lineStr),
				},
			},
		})
	}
	return elements
}

func (doc *PlainDocument) Render() []Element {
	if doc.cacheVersion == 0 || (doc.cacheVersion != doc.version) {
		doc.cache = doc.doRender()
	}
	doc.cacheVersion = doc.version
	return doc.cache
}

func (doc *PlainDocument) InsertString(row int, bytePos int, s string) {
	doc.buffer.InsertString(row, bytePos, s)
	doc.version++
}

func (doc *PlainDocument) Remove(row int, bytePos int, byteLen int) {
	doc.buffer.RemoveString(row, bytePos, byteLen)
	doc.version++
}

// Clear is initialize Buffer.
func (doc *PlainDocument) Clear() {
	doc.buffer = Buffer{}
	doc.buffer.Init()
	doc.version++
}

func (doc *PlainDocument) GetLineBytes(index int) int {
	return len(doc.buffer.GetLineAt(index).GetContent())
}

func (doc *PlainDocument) GetLineCount() int {
	return doc.buffer.GetLineCount()
}

func (doc *PlainDocument) GetVersion() uint {
	return doc.version
}

func (doc *PlainDocument) GetLineString(lineIndex int) string {
	return doc.buffer.GetLineAt(lineIndex).GetContent()
}
