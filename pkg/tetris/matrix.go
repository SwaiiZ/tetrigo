package tetris

import (
	"errors"
	"fmt"
)

const (
	ScoreWeightHoles      = 10
	ScoreWeightUnfillable = 30
	ScoreWeightBumpiness  = 3
	ScoreWeightHeight     = 2
	ScoreWeightLines      = 100
)

// Matrix represents the board of cells on which the game is played.
type Matrix [][]byte

// DefaultMatrix creates a new Matrix with a height of 40 and a width of 10.
func DefaultMatrix() Matrix {
	m, err := NewMatrix(40, 10)
	if err != nil {
		panic(fmt.Errorf("failed to create default matrix: %w", err))
	}
	return m
}

// NewMatrix creates a new Matrix with the given height and width.
// It returns an error if the height is less than 20 to allow for a buffer zone of 20 lines.
func NewMatrix(height, width int) (Matrix, error) {
	if height <= 20 {
		return nil, errors.New("matrix height must be greater than 20 to allow for a buffer zone of 20 lines")
	}

	matrix := make(Matrix, height)
	for i := range matrix {
		matrix[i] = make([]byte, width)
	}
	return matrix, nil
}

// GetHeight returns the height of the Matrix.
func (m *Matrix) GetHeight() int {
	return len(*m)
}

func (m *Matrix) columnHeight(x int) int {
	for y := 0; y < len(*m); y++ {
		if (*m)[y][x] != 0 {
			return y
		}
	}
	return len(*m) // vide = hauteur max
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetSkyline returns the skyline; the highest row that the player can see.
func (m *Matrix) GetSkyline() int {
	return len(*m) - 20
}

func (m *Matrix) CountFullLines() int {
	count := 0
	for y := 0; y < len(*m); y++ {
		full := true
		for x := 0; x < len((*m)[y]); x++ {
			if (*m)[y][x] == 0 {
				full = false
				break
			}
		}
		if full {
			count++
		}
	}
	return count
}

func (m *Matrix) GetMaxHeight() int {
	mat := *m
	maxHeight := 0
	for mY := 0; mY < len(mat); mY++ {
		for mX := 0; mX < len(mat[mY]); mX++ {
			if mat[mY][mX] != 0 {
				maxHeight = mY
				break
			}
		}
		if maxHeight != 0 {
			break
		}
	}
	return maxHeight
}

func (m *Matrix) GetBumpiness() int {
	bumpiness := 0
	width := len((*m)[0])

	for x := 0; x < width-1; x++ {
		h1 := m.columnHeight(x)
		h2 := m.columnHeight(x + 1)
		bumpiness += abs(h1 - h2)
	}
	return bumpiness
}

func (m *Matrix) CountHoles() int {
	holes := 0
	for x := 0; x < len((*m)[0]); x++ {
		foundBlock := false
		for y := 0; y < len(*m); y++ {
			if (*m)[y][x] != 0 {
				foundBlock = true
			} else if foundBlock {
				holes++
			}
		}
	}
	return holes
}

func (m *Matrix) CountUnfillableHoles() int {
	unfillable := 0
	width := len((*m)[0])
	height := len(*m)

	// Pour chaque colonne
	for x := 0; x < width; x++ {
		blockSeen := false
		for y := 0; y < height; y++ {
			cell := (*m)[y][x]
			if cell != 0 {
				blockSeen = true
			} else if blockSeen {
				// Trou détecté
				leftHigher := x == 0 || m.columnHeight(x-1) > y
				rightHigher := x == width-1 || m.columnHeight(x+1) > y

				if leftHigher && rightHigher {
					unfillable++
				}
			}
		}
	}
	return unfillable
}

func (m *Matrix) EvaluateScore() int {
	lines := m.CountFullLines()
	holes := m.CountHoles()
	unfillable := m.CountUnfillableHoles()
	bump := m.GetBumpiness()
	height := m.GetMaxHeight()

	score := 0
	score += holes * ScoreWeightHoles
	score += unfillable * ScoreWeightUnfillable
	score += bump * ScoreWeightBumpiness
	score += height * ScoreWeightHeight
	score -= lines * ScoreWeightLines

	return score
}

func (m *Matrix) CanPlace(tCells [][]bool, offsetX, offsetY int) bool {
	for y := 0; y < len(tCells); y++ {
		for x := 0; x < len(tCells[y]); x++ {
			if !tCells[y][x] {
				continue
			}

			boardY := offsetY + y
			boardX := offsetX + x

			// Hors limites
			if boardY < 0 || boardY >= len(*m) || boardX < 0 || boardX >= len((*m)[0]) {
				return false
			}

			// Collision avec une cellule occupée
			if (*m)[boardY][boardX] != 0 {
				return false
			}
		}
	}
	return true
}

// GetVisible returns the Matrix without the buffer zone at the top (ie. the visible portion of the Matrix).
func (m *Matrix) GetVisible() Matrix {
	return (*m)[20:]
}

func (m *Matrix) DeepCopy() *Matrix {
	duplicate := make(Matrix, len(*m))
	for i := range *m {
		duplicate[i] = make([]byte, len((*m)[i]))
		copy(duplicate[i], (*m)[i])
	}
	return &duplicate
}

func (m *Matrix) isLineComplete(row int) bool {
	for _, cell := range (*m)[row] {
		if isCellEmpty(cell) {
			return false
		}
	}
	return true
}

func (m *Matrix) removeLine(row int) {
	(*m)[0] = make([]byte, len((*m)[0]))
	for i := row; i > 0; i-- {
		(*m)[i] = (*m)[i-1]
	}
}

// RemoveTetrimino removes the given Tetrimino from the Matrix.
// It returns an error if the Tetrimino is out of bounds or if the Tetrimino is not found in the Matrix.
func (m *Matrix) RemoveTetrimino(tet *Tetrimino) error {
	isExpectedValue := func(cell byte) bool {
		return cell == tet.Value
	}

	return m.modifyCell(tet.Cells, tet.Position, 0, isExpectedValue)
}

// AddTetrimino adds the given Tetrimino to the Matrix.
// It returns an error if the Tetrimino is out of bounds or if the Tetrimino overlaps with an occupied mino.
func (m *Matrix) AddTetrimino(tet *Tetrimino) error {
	return m.modifyCell(tet.Cells, tet.Position, tet.Value, isCellEmpty)
}

func (m *Matrix) modifyCell(minos [][]bool, pos Coordinate, newValue byte, isExpectedValue func(byte) bool) error {
	for row := range minos {
		for col := range minos[row] {
			if !minos[row][col] {
				continue
			}
			minoAbsRow := row + pos.Y
			minoAbsCol := col + pos.X

			if minoAbsRow >= len(*m) || minoAbsRow < 0 {
				return fmt.Errorf("row %d is out of bounds", minoAbsRow)
			}
			if minoAbsCol >= len((*m)[row]) || minoAbsCol < 0 {
				return fmt.Errorf("col %d is out of bounds", minoAbsCol)
			}

			minoValue := (*m)[minoAbsRow][minoAbsCol]
			if !isExpectedValue(minoValue) {
				// TODO: Perhaps there is a better way to do this:
				// Add in ghost minos is an exception. Occasionally the ghost mino will be
				// placed on top of a mino (eg. when playing at the skyline).
				if newValue != 'G' {
					return fmt.Errorf("mino at row %d, col %d is '%s' (byte value %v) not the expected value",
						minoAbsRow, minoAbsCol, string(minoValue), minoValue)
				}
			}
			(*m)[minoAbsRow][minoAbsCol] = newValue
		}
	}
	return nil
}

// RemoveCompletedLines checks each row that the given Tetrimino occupies and
// removes any completed lines from the Matrix.
// It returns an Action to be used for calculating the score.
func (m *Matrix) RemoveCompletedLines(tet *Tetrimino) Action {
	lines := 0
	for row := range tet.Cells {
		if m.isLineComplete(tet.Position.Y + row) {
			m.removeLine(tet.Position.Y + row)
			lines++
		}
	}

	switch lines {
	case 0:
		return Actions.None
	case 1:
		return Actions.Single
	case 2:
		return Actions.Double
	case 3:
		return Actions.Triple
	case 4:
		return Actions.Tetris
	}
	return Actions.Unknown
}

func (m *Matrix) isOutOfBoundsHorizontally(col int) bool {
	return col < 0 || col >= len((*m)[0])
}

func (m *Matrix) isOutOfBoundsVertically(row int) bool {
	return row < 0 || row >= len(*m)
}

func isCellEmpty(cell byte) bool {
	return cell == 0 || cell == 'G'
}

func (m *Matrix) canPlaceInCell(row, col int) bool {
	if m.isOutOfBoundsHorizontally(col) {
		return false
	}
	if m.isOutOfBoundsVertically(row) {
		return false
	}
	if !isCellEmpty((*m)[row][col]) {
		return false
	}
	return true
}
