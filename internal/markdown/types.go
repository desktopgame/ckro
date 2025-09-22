package markdown

import "github.com/desktopgame/ckro/internal/optional"

type Node struct {
	StartRune int
	EndRune   int
	StartLine int
	EndLine   int
	Attrs     map[string]string
}

func (n *Node) BaseNode() *Node {
	return n
}

type AbstractNode interface {
	BaseNode() *Node
}

type Block struct {
	Node
}

func (b *Block) BaseBlock() *Block {
	return b
}

type AbstractBlock interface {
	BaseBlock() *Block
}

type Document struct {
	Block
	Blocks []AbstractBlock
	Defs   map[string]string
}

type Paragraph struct {
	Block
	Inlines []AbstractInline
}

type Heading struct {
	Block
	Level   int
	Inlines []AbstractInline
	Id      optional.Optional[string]
}

type CodeBlock struct {
	Block
	Info string
	Code string
}

type MathBlock struct {
	Block
	Tex string
}

type Blockquote struct {
	Block
	Blocks []AbstractBlock
}

type ThematicBreak struct {
	Block
}

type ListItem struct {
	Node
	Paragraph optional.Optional[Paragraph]
	Checked   optional.Optional[bool]
	Blocks    []AbstractBlock
}

type ListBlock struct {
	Block
	Ordered bool
	Start   int
	Tight   bool
	Items   []*ListItem
}

type TableRow struct {
	Columns []Paragraph
}

type Table struct {
	Block
	Headers []string
	Aligns  []string
	Rows    []TableRow
}

type Inline struct {
	Node
}

func (il *Inline) BaseInline() *Inline {
	return il
}

type AbstractInline interface {
	BaseInline() *Inline
}

type Text struct {
	Inline
	Text string
}

type Emph struct {
	Inline
	Children []AbstractInline
}

type Strong struct {
	Inline
	Children []AbstractInline
}

type Code struct {
	Inline
	Text string
}

type Math struct {
	Inline
	Tex string
}

type Strike struct {
	Inline
	Children []AbstractInline
}

type Tag struct {
	Inline
	Text string
}

type Link struct {
	Inline
	Children []AbstractInline
	Url      string
	Title    optional.Optional[string]
}

type Image struct {
	Inline
	Alt   string
	Url   string
	Title optional.Optional[string]
}

type SoftBreak struct {
	Inline
}

type HardBreak struct {
	Inline
}

type ExtInline struct {
	Inline
	Name     string
	Children []AbstractInline
	Data     map[string]string
}

type ExtBlock struct {
	Block
	Name   string
	Blocks []AbstractBlock
}
