package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth          = 400
	screenHeight         = 400
	halfGameScreenWidth  = screenWidth / 2
	halfGameScreenHeight = screenHeight / 2
	blockLength          = 20
	unitsX               = screenWidth / blockLength
	unitsY               = screenHeight / blockLength
)

type Direction int

const (
	// right is 0, up is 1, left is 2 and down is 3
	right Direction = iota
	up
	left
	down
)

type Vec2D struct {
	x int
	y int
}

type Snake struct {
	head Vec2D
	body []Vec2D
	d    Direction
}

type Game struct {
	grid   [unitsX][unitsY]int
	player *Snake
	time   int
	food   Vec2D
	score  int
}

type GameObject interface {
	Update()
	Draw(screen *ebiten.Image)
}

func (p *Snake) input() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		p.d = up
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		p.d = down
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		p.d = left
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		p.d = right
	}
}

func (g *Game) Update() error {
	// update direction based on input
	g.player.input()
	// updating snake position every second
	if int(time.Now().Unix()) > g.time {
		g.time = int(time.Now().Unix())

		// updating position of body first
		for i := len(g.player.body) - 1; i > -1; i-- {
			if i != 0 {
				g.player.body[i] = g.player.body[i-1]
			} else {
				g.player.body[i] = g.player.head
			}
		}

		// updating position of head based on direction set
		switch g.player.d {
		case right:
			g.player.head.x++
		case up:
			g.player.head.y--
		case left:
			g.player.head.x--
		case down:
			g.player.head.y++
		}

		// detecting collisions with wall and body
		if g.player.head.x < 0 || g.player.head.x >= unitsX || g.player.head.y < 0 || g.player.head.y >= unitsY {
			// TODO: gracefully handle game over
			fmt.Printf("game over")
			// exit game
			os.Exit(0)
		}
		for _, e := range g.player.body {
			if g.player.head.x == e.x && g.player.head.y == e.y {
				// TODO: gracefully handle game over
				fmt.Printf("game over")
				os.Exit(0)
			}
		}

		// detecting collision with food
		if g.player.head.x == g.food.x && g.player.head.y == g.food.y {
			g.player.body = append(g.player.body, Vec2D{g.player.body[len(g.player.body)-1].x, g.player.body[len(g.player.body)-1].y})
			g.food = spawnfood(g)
			g.score = g.score + 1
			fmt.Println("score: ", g.score)
		}
	}

	g.time = int(time.Now().Unix())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.player.head.x), float64(g.player.head.y))

	// drawing the head
	vector.DrawFilledRect(screen,
		float32(g.player.head.x*blockLength),
		float32(g.player.head.y*blockLength),
		float32(blockLength),
		float32(blockLength),
		color.RGBA{65, 77, 68, 0}, true)

	// drawing the body
	for _, e := range g.player.body {
		vector.DrawFilledRect(screen,
			float32(e.x*blockLength),
			float32(e.y*blockLength),
			float32(blockLength),
			float32(blockLength),
			color.White, true)
	}

	// drawing the food
	vector.DrawFilledRect(screen, float32(g.food.x)*float32(blockLength), float32(g.food.y)*float32(blockLength), float32(blockLength), float32(blockLength), color.RGBA{255, 0, 0, 0}, true)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func newGame() *Game {
	// setting up the snake at the center of the screen
	head := Vec2D{unitsX / 2, unitsY / 2}
	body := []Vec2D{}
	body = append(body, Vec2D{head.x - 1, head.y})
	body = append(body, Vec2D{head.x - 2, head.y})

	game := &Game{
		grid:   [unitsX][unitsY]int{},
		player: &Snake{head: head, body: body, d: right},
	}

	food := spawnfood(game)
	game.food = food
	game.score = 0
	return game
}

// spawn food
func spawnfood(game *Game) Vec2D {
	food := Vec2D{}

	for {
		food = Vec2D{rand.Intn(unitsX), rand.Intn(unitsY)}
		if game.grid[food.x][food.y] == 0 {
			game.grid[food.x][food.y] = 1
			return food
		}
	}
}

func main() {
	// ebiten init setup
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowTitle("Snake")
	ebiten.SetFullscreen(false)

	// creating a new game object
	game := newGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
