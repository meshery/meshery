package utils

func levenshtein(str1, str2 []rune) int {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return column[s1len]
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

func minIntSlice(values []int) (int, int) {
	min := values[0]
	minIndex := 0
	for i, val := range values {
		if min > val {
			min = val
			minIndex = i
		}
	}
	return minIndex, min
}

func FindClosestArg(argPassed string, validArgs []string) string {
	argPassedR := []rune(argPassed)
	var values []int
	for _, arg := range validArgs {
		argR := []rune(arg)
		values = append(values, levenshtein(argPassedR, argR))
	}
	minIndex, _ := minIntSlice(values)
	return validArgs[minIndex]
}
