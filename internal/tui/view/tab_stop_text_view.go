package view

type TabStopTextView interface {
	WidthWithTabStop(textViewResolver TextViewResolver, textLayout *TextLayout, column int) int
}
