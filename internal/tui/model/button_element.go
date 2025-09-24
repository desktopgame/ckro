package model

type ButtonElement struct {
	StartPosition Position
	EndPosition   Position
}

func (b *ButtonElement) GetStyle() *Style {
	return nil
}

func (b *ButtonElement) GetText() string {
	return ""
}

func (b *ButtonElement) GetStartPosition() Position {
	return b.StartPosition
}

func (b *ButtonElement) GetEndPosition() Position {
	return b.EndPosition
}

func (b *ButtonElement) GetElement(index int) Element {
	return nil
}

func (b *ButtonElement) GetElementCount() int {
	return 0
}
