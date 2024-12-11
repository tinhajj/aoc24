package structure

type Point struct {
	Y int
	X int
}

type VertexInt struct {
	Point Point
	Val   int
}

func VertexMatrixInt(digiMatrix [][]int) [][]*VertexInt {
	vertMatrix := [][]*VertexInt{}

	for i, digiRow := range digiMatrix {
		vertRow := []*VertexInt{}
		for j, digi := range digiRow {
			v := &VertexInt{
				Point: Point{Y: i, X: j},
				Val:   digi,
			}
			vertRow = append(vertRow, v)
		}
		vertMatrix = append(vertMatrix, vertRow)
	}

	return vertMatrix
}
