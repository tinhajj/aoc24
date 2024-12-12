package structure

type Point struct {
	Y int
	X int
}

type VertexInt struct {
	Point Point
	Val   int
}

type VertexStr struct {
	Point Point
	Val   string
}

func VertexMatrixStr(m [][]string) [][]*VertexStr {
	matrix := [][]*VertexStr{}

	for i, r := range m {
		row := []*VertexStr{}
		for j, str := range r {
			v := &VertexStr{
				Point: Point{Y: i, X: j},
				Val:   str,
			}
			row = append(row, v)
		}
		matrix = append(matrix, row)
	}

	return matrix
}

func VertexMatrixInt(m [][]int) [][]*VertexInt {
	matrix := [][]*VertexInt{}

	for i, r := range m {
		row := []*VertexInt{}
		for j, digit := range r {
			v := &VertexInt{
				Point: Point{Y: i, X: j},
				Val:   digit,
			}
			row = append(row, v)
		}
		matrix = append(matrix, row)
	}

	return matrix
}
