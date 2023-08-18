package parsers

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func StringDiffScore(str1, str2 string) int {
	m, n := len(str1), len(str2)
	matrix := make([][]int, m+1)

	for i := 0; i <= m; i++ {
		matrix[i] = make([]int, n+1)
		matrix[i][0] = i
	}

	for j := 0; j <= n; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 0
			if str1[i-1] != str2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}

	return matrix[m][n]
}
