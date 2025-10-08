package view

func CompositeViewLengthTable(ctx Context, textLayout *TextLayout) ([]int, int) {
	var table []int
	total := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := ctx.Resolver.Resolve(childElement)

		l := childView.MoveLength(ctx, textLayout.Children[i])
		table = append(table, l)
		total += l
	}
	return table, total
}

func CompositeViewLengthSum(table []int, index int) int {
	v := 0
	for i := 0; i <= index; i++ {
		v += table[i]
	}
	return v
}

func CompositeViewIndex(table []int, viewLocalPos int) (Row int, Column int) {
	n := 0
	index := -1
	col := -1
	for i, l := range table {
		start := n
		if viewLocalPos >= start && viewLocalPos < n+l {
			index = i
			col = viewLocalPos - start
			break
		}
		n += l
	}
	return index, col
}

func CompositeMoveUp(ctx Context, textLayout *TextLayout, viewLocalPos int) int {
	table, _ := CompositeViewLengthTable(ctx, textLayout)
	index, col := CompositeViewIndex(table, viewLocalPos)

	textView := ctx.Resolver.Resolve(textLayout.Children[index].Element)
	newCol := textView.MoveUp(ctx, textLayout.Children[index], col)
	if newCol == -1 {
		if index == 0 {
			return -1
		}
		return CompositeViewLengthSum(table, index-1) - 1
	}
	return CompositeViewLengthSum(table, index-1) + newCol
}

func CompositeMoveDown(ctx Context, textLayout *TextLayout, viewLocalPos int) int {
	table, _ := CompositeViewLengthTable(ctx, textLayout)
	index, col := CompositeViewIndex(table, viewLocalPos)

	textView := ctx.Resolver.Resolve(textLayout.Children[index].Element)
	newCol := textView.MoveDown(ctx, textLayout.Children[index], col)
	if newCol == -1 {
		if index == len(table)-1 {
			return -1
		}
		return CompositeViewLengthSum(table, index)
	}
	return CompositeViewLengthSum(table, index-1) + newCol
}
