package src

type StreamReader struct {
	Stream      []rune
	curPosition int
}

func NewStreamReader(stream string) *StreamReader {
	sr := &StreamReader{
		Stream:      []rune(stream),
		curPosition: -1,
	}
	return sr
}

func (sr *StreamReader) Read() rune {
	sr.curPosition += 1
	if sr.curPosition >= len(sr.Stream) {
		return rune('0')
	}
	return sr.Stream[rune(sr.curPosition)]
}

func (sr *StreamReader) Peek() rune {
	peekPos := sr.curPosition + 1
	if peekPos >= len(sr.Stream) {
		return 0
	}
	return sr.Stream[rune(sr.curPosition)]
}

func (sr *StreamReader) EndOfStream() bool {
	return sr.curPosition+1 >= len(sr.Stream)
}
