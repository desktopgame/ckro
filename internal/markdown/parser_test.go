package markdown_test

import (
	"testing"

	"github.com/desktopgame/ckro/internal/markdown"
	"github.com/stretchr/testify/assert"
)

func TestList01(t *testing.T) {
	text := `
* item1
* item2
* item3
`
	doc := markdown.Parse(text)
	assert.Equal(t, len(doc.Blocks), 1)

	root := doc.Blocks[0]
	if lb, ok := root.(*markdown.ListBlock); ok {
		assert.Equal(t, len(lb.Items), 3)
	} else {
		assert.True(t, false)
	}
}

func TestList02(t *testing.T) {
	text := `
* item1
  * item2
    * item3
`
	doc := markdown.Parse(text)
	assert.Equal(t, len(doc.Blocks), 1)

	root := doc.Blocks[0]
	if lb, ok := root.(*markdown.ListBlock); ok {
		assert.Equal(t, len(lb.Items), 1)
	} else {
		assert.True(t, false)
	}
}

func TestList03(t *testing.T) {
	text := `
* item1
  * item2
* item3
  * item4
`
	doc := markdown.Parse(text)
	assert.Equal(t, len(doc.Blocks), 1)

	root := doc.Blocks[0]
	if lb, ok := root.(*markdown.ListBlock); ok {
		assert.Equal(t, len(lb.Items), 2)
		assert.Equal(t, len(lb.Items[0].Blocks), 1)
		assert.Equal(t, len(lb.Items[1].Blocks), 1)
	} else {
		assert.True(t, false)
	}
}

func TestList04(t *testing.T) {
	text := `
* item1
  * item2
    * item3
      * item4
    * item5
  * item6
`
	doc := markdown.Parse(text)
	assert.Equal(t, len(doc.Blocks), 1)

	root := doc.Blocks[0]
	if lb, ok := root.(*markdown.ListBlock); ok {
		assert.Equal(t, len(lb.Items), 1)
		assert.Equal(t, len(lb.Items[0].Blocks), 1)

		if tmp, ok := lb.Items[0].Blocks[0].(*markdown.ListBlock); ok {
			assert.Equal(t, len(tmp.Items), 2)
		} else {
			assert.True(t, false)
		}
		if tmp, ok := lb.Items[0].Blocks[0].(*markdown.ListBlock); ok {
			assert.Equal(t, len(tmp.Items[0].Blocks), 1)
		} else {
			assert.True(t, false)
		}
	} else {
		assert.True(t, false)
	}
}

func TestList05(t *testing.T) {
	text := `
* item1
  * item2
    * item3
      * item4
    * item5
    * item6
  * item7
  * item8
`
	doc := markdown.Parse(text)
	assert.Equal(t, len(doc.Blocks), 1)

	root := doc.Blocks[0]
	if lb, ok := root.(*markdown.ListBlock); ok {
		assert.Equal(t, len(lb.Items), 1)
		assert.Equal(t, len(lb.Items[0].Blocks), 1)

		if tmp, ok := lb.Items[0].Blocks[0].(*markdown.ListBlock); ok {
			assert.Equal(t, len(tmp.Items), 3)
		} else {
			assert.True(t, false)
		}
		if tmp, ok := lb.Items[0].Blocks[0].(*markdown.ListBlock); ok {
			assert.Equal(t, len(tmp.Items[0].Blocks), 1)
		} else {
			assert.True(t, false)
		}
	} else {
		assert.True(t, false)
	}
}

func TestList06(t *testing.T) {
	text := `
# Test

* Item1
* Item2
* Item3

- Item1
- Item2
- Item3

- Unchecked
- Checked

* Item1
  * Item2
  * Item3

* Item1
  * Item2
    * Item3
	`
	doc := markdown.Parse(text)
	// Should have 1 heading + 5 list blocks = 6 blocks total
	assert.Equal(t, len(doc.Blocks), 6)

	// First block should be a heading
	heading := doc.Blocks[0]
	_, ok := heading.(*markdown.Heading)
	assert.True(t, ok)

	// Remaining blocks should be lists
	for i := 1; i < len(doc.Blocks); i++ {
		_, ok := doc.Blocks[i].(*markdown.ListBlock)
		assert.True(t, ok)
	}
}

func TestList07(t *testing.T) {
	text := `
# Test

* [ ] item1
  * [ ] item2
  * [x] item3
    * [] item4
* [ ] item5
* [x] item6
	`
	doc := markdown.Parse(text)
	assert.Equal(t, len(doc.Blocks), 2)

	_, ok := doc.Blocks[0].(*markdown.Heading)
	assert.True(t, ok)

	if list, ok := doc.Blocks[1].(*markdown.ListBlock); ok {
		assert.Equal(t, len(list.Items), 3)
		assert.Equal(t, len(list.Items[0].Blocks), 1)

		if list2, ok := list.Items[0].Blocks[0].(*markdown.ListBlock); ok {
			assert.Equal(t, len(list2.Items), 2)
		}
	} else {
		assert.True(t, false)
	}
}

func TestTable01(t *testing.T) {
	text := `
|H1|H2|H3|
|---|---|---|
|co1|co2|co3|
	`
	doc := markdown.Parse(text)
	assert.Equal(t, len(doc.Blocks), 1)

	if tb, ok := doc.Blocks[0].(*markdown.Table); ok {
		assert.Equal(t, len(tb.Headers), 3)
		assert.Equal(t, tb.Headers[0], "H1")
		assert.Equal(t, tb.Headers[1], "H2")
		assert.Equal(t, tb.Headers[2], "H3")
		assert.Equal(t, len(tb.Aligns), 3)
		assert.Equal(t, tb.Aligns[0], "left")
		assert.Equal(t, tb.Aligns[1], "left")
		assert.Equal(t, tb.Aligns[2], "left")
	} else {
		assert.True(t, false)
	}
}

func TestTable02(t *testing.T) {
	text := `
|H1|H2|H3|
|---:|:---:|:---|
|co1|co2|co3|
	`
	doc := markdown.Parse(text)
	assert.Equal(t, len(doc.Blocks), 1)

	if tb, ok := doc.Blocks[0].(*markdown.Table); ok {
		assert.Equal(t, len(tb.Headers), 3)
		assert.Equal(t, tb.Headers[0], "H1")
		assert.Equal(t, tb.Headers[1], "H2")
		assert.Equal(t, tb.Headers[2], "H3")
		assert.Equal(t, len(tb.Aligns), 3)
		assert.Equal(t, tb.Aligns[0], "right")
		assert.Equal(t, tb.Aligns[1], "center")
		assert.Equal(t, tb.Aligns[2], "left")
	} else {
		assert.True(t, false)
	}
}
