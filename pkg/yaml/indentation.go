package yamlsanitize

import "strings"

func detectMixedIndentationRows(lines []string) (int, []int) {
	indentCounts := map[int]int{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		spaces := leadingSpaces(line)
		if spaces > 0 {
			indentCounts[spaces]++
		}
	}

	if len(indentCounts) == 0 {
		return 0, nil
	}

	unit := dominantIndentUnit(indentCounts)
	if unit <= 0 {
		return 0, nil
	}

	offenders := make([]int, 0)
	for i, line := range lines {
		spaces := leadingSpaces(line)
		if spaces > 0 && spaces%unit != 0 {
			offenders = append(offenders, i)
		}
	}

	return unit, offenders
}

func dominantIndentUnit(indentCounts map[int]int) int {
	gcdAll := 0
	for width := range indentCounts {
		gcdAll = gcd(gcdAll, width)
	}
	if gcdAll <= 0 {
		gcdAll = 2
	}

	unit := gcdAll
	if unit == 1 {
		bestWidth, bestCount := 2, 0
		for width, count := range indentCounts {
			if count > bestCount {
				bestWidth, bestCount = width, count
			}
		}
		unit = bestWidth
	}
	return unit
}

func leadingSpaces(line string) int {
	spaces := 0
	for _, ch := range line {
		if ch == ' ' {
			spaces++
			continue
		}
		break
	}
	return spaces
}
