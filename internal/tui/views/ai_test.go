package views

import (
	"testing"

	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
)

func BenchmarkFindBestPlacement(b *testing.B) {
	tetriminos := []tetris.Tetrimino{
		{Value: 'I', Cells: [][]bool{{true, true, true, true}}},
		{Value: 'O', Cells: [][]bool{{true, true}, {true, true}}},
		{Value: 'T', Cells: [][]bool{{true, true, true}, {false, true, false}}},
		{Value: 'S', Cells: [][]bool{{false, true, true}, {true, true, false}}},
		{Value: 'Z', Cells: [][]bool{{true, true, false}, {false, true, true}}},
		{Value: 'J', Cells: [][]bool{{true, false, false}, {true, true, true}}},
		{Value: 'L', Cells: [][]bool{{false, false, true}, {true, true, true}}},
	}

	// Generate a few matrix states (simulate boards with some height)
	var matrices []tetris.Matrix
	for i := 0; i < 5; i++ {
		m := tetris.DefaultMatrix()

		// Add some stacked blocks in a column pattern
		for y := 35; y < 39; y++ {
			for x := 2; x < 4+i%3; x++ {
				m[y][x] = 'X'
			}
		}
		matrices = append(matrices, m)
	}

	model := AIModel{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Rotate through matrix/tetrimino scenarios
		matrix := matrices[i%len(matrices)]
		tet := tetriminos[i%len(tetriminos)]
		tet1 := tetriminos[(i+1)%len(tetriminos)]
		tet2 := tetriminos[(i+2)%len(tetriminos)]
		tet.Position = tetris.Coordinate{X: 0, Y: 0}

		model.FindBestPlacementSequence(matrix, []tetris.Tetrimino{tet, tet1, tet2}, 2)
	}
}
