package jsonsanitize

type lineIndex struct {
	starts []int
}

func newLineIndex(src string) lineIndex {
	starts := []int{0}
	for i := 0; i < len(src); i++ {
		if src[i] == '\n' && i+1 <= len(src) {
			starts = append(starts, i+1)
		}
	}
	return lineIndex{starts: starts}
}
