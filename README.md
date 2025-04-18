# Tetrigo

- Mathieu ROCHER
- Bastien LE BAIL

![app demo](./docs/readme-demo.gif)

## Problématique
> Comment optimiser un solveur tétris ?

## Différentes versions
| Versions |                                                                   Branche/Tag                                                                   |
|:---------|:-----------------------------------------------------------------------------------------------------------------------------------------------:|
| v0       |                     [Lien v0 - Tag vAI_No_rotate](https://github.com/SwaiiZ/tetrigo/tree/vAI_No_rotate?tab=readme-ov-file)                      |
| v1       |                    [Lien v1 - Tag v1_rotating_ai](https://github.com/SwaiiZ/tetrigo/tree/v1_rotating_ai?tab=readme-ov-file)                     |
| v2       |                    [Lien v2 - Tag v2_rotating_ai](https://github.com/SwaiiZ/tetrigo/tree/v2_rotating_ai?tab=readme-ov-file)                     |
| v3       |                    [Lien v3 - Tag v3_rotating_ai](https://github.com/SwaiiZ/tetrigo/tree/v3_rotating_ai?tab=readme-ov-file)                     |
| v4       | [Lien v4 - Branch feature/ai_mode_sequence_no_opti](https://github.com/SwaiiZ/tetrigo/tree/feature/ai_mode_sequence_no_opti?tab=readme-ov-file) |
| v5       |    [Lien v5 - Branch feature/ai_mode_sequence_opti](https://github.com/SwaiiZ/tetrigo/tree/feature/ai_mode_sequence_opti?tab=readme-ov-file)    |
| v6       | [Lien v6 - Branch feature/ai_mode_sequence_opti_wg](https://github.com/SwaiiZ/tetrigo/tree/feature/ai_mode_sequence_opti_wg?tab=readme-ov-file) |

## Code benchmark
Dans un fichier ***ai_test.go*** dans ***internal/tui/views***

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
            // FindBestPlacement est pour les versions sans sequence
            // model.FindBestPlacement(matrix,tet)
            model.FindBestPlacementSequence(matrix,[]tetris.Tetrimino{ tet,tet1,tet2},2)
        }
    }