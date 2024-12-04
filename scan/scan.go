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
