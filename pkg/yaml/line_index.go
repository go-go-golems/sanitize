package yamlsanitize

import "sort"

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

func (li lineIndex) rowAtByte(byteOffset uint) int {
	if len(li.starts) == 0 {
		return 0
	}

	offset := int(byteOffset)
	if offset < 0 {
		offset = 0
	}

	row := sort.Search(len(li.starts), func(i int) bool {
		return li.starts[i] > offset
	}) - 1
	if row < 0 {
		return 0
	}
	return row
}

func (li lineIndex) rowColAtByte(byteOffset uint) (int, int) {
	row := li.rowAtByte(byteOffset)
	if row >= len(li.starts) {
		return row, 0
	}

	col := int(byteOffset) - li.starts[row]
	if col < 0 {
		col = 0
	}
	return row, col
}

func lineIndexAtByte(src string, byteOffset uint) int {
	return newLineIndex(src).rowAtByte(byteOffset)
}
