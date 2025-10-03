package litemark

type Span struct {
	StartColumn int
	EndColumn   int
}

type Block struct {
	LineIndex int
	LineCount int
}

func (b *Block) BaseBlock() *Block {
	return b
}

type AbstractBlock interface {
	BaseBlock() *Block
}

type Heading struct {
	Block
	Span  Span
	Level int
}

type BlankLine struct {
	Block
}

type CodeBlock struct {
	Block
	Span Span
}

type FoldBlock struct {
	Block
}

type HorizontalLine struct {
	Block
}

type Inline struct {
	Spans []Span
}

func (i *Inline) BaseInline() *Inline {
	return i
}

type AbstractInline interface {
	BaseInline() *Inline
}

type Code struct {
	Inline
}

type Italic struct {
	Inline
}

type Bold struct {
	Inline
}

type Strike struct {
	Inline
}

type Link struct {
	Inline
}

type Image struct {
	Inline
}

type PlainText struct {
	Inline
}

type Text struct {
	Block
	Inlines []AbstractInline
}
