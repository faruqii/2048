package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SCREEN_WIDTH  = 400
	SCREEN_HEIGHT = 480
	TILE_SIZE     = 80
	GRID_SIZE     = 4
)

type Game struct {
	board [GRID_SIZE][GRID_SIZE]int
	score int
}

var (
	tileColors = map[int]color.RGBA{
		0:    {205, 193, 180, 255}, // Empty
		2:    {238, 228, 218, 255},
		4:    {237, 224, 200, 255},
		8:    {242, 177, 121, 255},
		16:   {245, 149, 99, 255},
		32:   {246, 124, 95, 255},
		64:   {246, 94, 59, 255},
		128:  {237, 207, 114, 255},
		256:  {237, 204, 97, 255},
		512:  {237, 200, 80, 255},
		1024: {237, 197, 63, 255},
		2048: {237, 194, 46, 255},
	}
)

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	g := &Game{}
	g.AddRandomTile()
	g.AddRandomTile()
	return g
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.Move("up")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.Move("down")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.Move("left")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.Move("right")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	screen.Fill(color.RGBA{238, 228, 218, 255})

	// Draw tiles
	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			x := float64(j*TILE_SIZE + 20)
			y := float64(i*TILE_SIZE + 100)
			value := g.board[i][j]

			// Draw tile
			tile := ebiten.NewImage(TILE_SIZE-10, TILE_SIZE-10)
			tile.Fill(tileColors[value])
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x+5, y+5)
			screen.DrawImage(tile, op)

			// Draw number
			if value != 0 {
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", value), int(x)+30, int(y)+30)
			}
		}
	}

	// Draw score
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.score))

	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Use WASD or Arrow Keys to move", 20, 20)
	if g.IsGameOver() {
		ebitenutil.DebugPrintAt(screen, "Game Over! Press R to restart", 120, SCREEN_HEIGHT/2)
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			*g = *NewGame()
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func (g *Game) AddRandomTile() {
	var emptyTiles [][2]int
	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			if g.board[i][j] == 0 {
				emptyTiles = append(emptyTiles, [2]int{i, j})
			}
		}
	}
	if len(emptyTiles) > 0 {
		pos := emptyTiles[rand.Intn(len(emptyTiles))]
		if rand.Float32() < 0.9 {
			g.board[pos[0]][pos[1]] = 2
		} else {
			g.board[pos[0]][pos[1]] = 4
		}
	}
}

func (g *Game) Move(direction string) {
	original := g.board
	switch direction {
	case "up":
		g.moveUp()
	case "down":
		g.moveDown()
	case "left":
		g.moveLeft()
	case "right":
		g.moveRight()
	}
	if !g.boardsEqual(original, g.board) {
		g.AddRandomTile()
	}
}

// Movement functions (same as console version)
func (g *Game) moveLeft() {
	for i := 0; i < GRID_SIZE; i++ {
		row := g.board[i]
		newRow := [GRID_SIZE]int{}
		pos := 0
		for j := 0; j < GRID_SIZE; j++ {
			if row[j] != 0 {
				if pos > 0 && newRow[pos-1] == row[j] {
					newRow[pos-1] *= 2
					g.score += newRow[pos-1]
				} else {
					newRow[pos] = row[j]
					pos++
				}
			}
		}
		g.board[i] = newRow
	}
}

func (g *Game) moveRight() {
	g.rotate180()
	g.moveLeft()
	g.rotate180()
}

func (g *Game) moveUp() {
	for col := 0; col < GRID_SIZE; col++ {
		var newCol []int
		for row := 0; row < GRID_SIZE; row++ {
			if g.board[row][col] != 0 {
				newCol = append(newCol, g.board[row][col]) // Collect non-zero values
			}
		}

		// Merge adjacent same values
		for i := 0; i < len(newCol)-1; i++ {
			if newCol[i] == newCol[i+1] {
				newCol[i] *= 2
				newCol[i+1] = 0 // Mark as merged
			}
		}

		// Shift again after merging
		finalCol := []int{}
		for _, num := range newCol {
			if num != 0 {
				finalCol = append(finalCol, num)
			}
		}

		// Fill the column with zeros
		for len(finalCol) < GRID_SIZE {
			finalCol = append(finalCol, 0)
		}

		// Copy back to board
		for row := 0; row < GRID_SIZE; row++ {
			g.board[row][col] = finalCol[row]
		}
	}
}

func (g *Game) moveDown() {
	for col := 0; col < GRID_SIZE; col++ {
		var newCol []int
		for row := GRID_SIZE - 1; row >= 0; row-- {
			if g.board[row][col] != 0 {
				newCol = append(newCol, g.board[row][col]) // Collect non-zero values
			}
		}

		// Merge adjacent same values
		for i := 0; i < len(newCol)-1; i++ {
			if newCol[i] == newCol[i+1] {
				newCol[i] *= 2
				newCol[i+1] = 0 // Mark as merged
			}
		}

		// Shift again after merging
		finalCol := []int{}
		for _, num := range newCol {
			if num != 0 {
				finalCol = append(finalCol, num)
			}
		}

		// Fill the column with zeros
		for len(finalCol) < GRID_SIZE {
			finalCol = append(finalCol, 0)
		}

		// Copy back to board (from bottom to top)
		for row, i := GRID_SIZE-1, 0; row >= 0; row, i = row-1, i+1 {
			g.board[row][col] = finalCol[i]
		}
	}
}

func (g *Game) rotate90() {
	newBoard := [GRID_SIZE][GRID_SIZE]int{}
	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			newBoard[j][GRID_SIZE-1-i] = g.board[i][j]
		}
	}
	g.board = newBoard
}

func (g *Game) rotate180() {
	g.rotate90()
	g.rotate90()
}

func (g *Game) boardsEqual(b1, b2 [GRID_SIZE][GRID_SIZE]int) bool {
	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			if b1[i][j] != b2[i][j] {
				return false
			}
		}
	}
	return true
}

func (g *Game) IsGameOver() bool {
	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			if g.board[i][j] == 0 {
				return false
			}
		}
	}
	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE-1; j++ {
			if g.board[i][j] == g.board[i][j+1] {
				return false
			}
		}
	}
	for j := 0; j < GRID_SIZE; j++ {
		for i := 0; i < GRID_SIZE-1; i++ {
			if g.board[i][j] == g.board[i+1][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("2048")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
