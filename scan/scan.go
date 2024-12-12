package scan

import "strconv"

// Numbers scans for contiguous sections of numbers in a string
func Numbers(input string) []int {
	nums := []int{}
	var chunk string

	parse := func() {
		if chunk == "-" || len(chunk) < 1 {
			return
		}
		number, err := strconv.Atoi(chunk)
		if err != nil {
			panic(err)
		}
		nums = append(nums, number)
		chunk = ""
	}

	for i := 0; i <= len(input); i++ {
		if i == len(input) {
			parse()
			continue
		}

		r := input[i]
		if (r < '0' || r > '9') && r != '-' {
			parse()
			continue
		}

		chunk += string(r)
	}

	return nums
}

// DigitMatrix scans for a rectangle of single digit numbers
func DigitMatrix(input []string, charTransform func(s string) int) (matrix [][]int, height int, width int) {
	matrix = [][]int{}

	for i, line := range input {
		if line == "" {
			continue
		}
		matrix = append(matrix, []int{})
		for _, c := range line {
			number, err := strconv.Atoi(string(c))
			if err != nil {
				if charTransform != nil {
					r := charTransform(string(c))
					matrix[i] = append(matrix[i], r)
					continue
				}
				panic(err)
			}
			matrix[i] = append(matrix[i], number)
		}
	}

	return matrix, len(matrix), len(matrix[0])
}

// RuneMatrix scans for a rectangle of single runes
func RuneMatrix(input []string) (matrix [][]string, height int, width int) {
	matrix = [][]string{}

	for i, line := range input {
		if line == "" {
			continue
		}
		matrix = append(matrix, []string{})
		for _, c := range line {
			matrix[i] = append(matrix[i], string(c))
		}
	}

	return matrix, len(matrix), len(matrix[0])
}
