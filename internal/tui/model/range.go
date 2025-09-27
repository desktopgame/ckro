package model

type Range struct {
	// 範囲の最初の文字の最初のバイト（これを含む）
	StartPosition Position
	// 範囲の最後の文字の最後のバイト（これを含まない）
	EndPosition Position
}
