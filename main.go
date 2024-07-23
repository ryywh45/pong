package main

import (
	"log"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"fmt"
)

const (
	screenWidth  = 640
	screenHeight = 480
	ballSpeed    = 3
	paddleSpeed  = 6
)

type BaseObj struct {
	X, Y, W, H int
}

type Paddle struct {
	BaseObj
}

type Ball struct {
	BaseObj
	dxdt int // x velocity per tick
	dydt int // y velocity per tick
}

type Game struct {
	Paddle
	Ball
	score int
	highScore int
}

func main() {
	ebiten.SetWindowTitle("PONG")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	paddle := Paddle{
		BaseObj: BaseObj{ X: 600, Y: 200, W: 15, H: 100 },
	}
	ball := Ball{
		BaseObj: BaseObj{ X: 0, Y: 0, W: 15, H: 15 },
		dxdt: ballSpeed,
		dydt: ballSpeed,
	}
	g := &Game{
		Paddle: paddle,
		Ball: ball,
	}

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image){
	vector.DrawFilledRect(screen,
		float32(g.Paddle.X), float32(g.Paddle.Y),
		float32(g.Paddle.W), float32(g.Paddle.H),
		color.White, false,
	)
	vector.DrawFilledRect(screen,
		float32(g.Ball.X), float32(g.Ball.Y),
		float32(g.Ball.W), float32(g.Ball.H),
		color.White, false,
	)
	
	scoreStr := fmt.Sprint("Score: ", g.score)
	text.Draw(screen, scoreStr, basicfont.Face7x13, 10, 10 ,color.White)

	highScoreStr := fmt.Sprint("High Score: ", g.highScore)
	text.Draw(screen, highScoreStr, basicfont.Face7x13, 10, 30, color.White)
}

func (g *Game) Update() error {
	g.Paddle.MoveOnKeyPressed()
	g.Ball.Move()
	g.CheckCollision()
	return nil
}

func (g *Game) Reset() {
	g.Ball.X = 0
	g.Ball.Y = 0
	g.score = 0
}

func (p *Paddle) MoveOnKeyPressed() {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyJ):
		p.Y += paddleSpeed
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyK):
		p.Y -= paddleSpeed
	}
}

func (b *Ball) Move() {
	b.X += b.dxdt
	b.Y += b.dydt
}

func (g *Game) CheckCollision() {
	// with wall
	switch {
	case g.Ball.X >= screenWidth:  // right wall
		g.Reset()
	case g.Ball.X <= 0:            // left wall
		g.Ball.dxdt = ballSpeed
	case g.Ball.Y <= 0:            // top wall
		g.Ball.dydt = ballSpeed
	case g.Ball.Y >= screenHeight: // bottom wall
		g.Ball.dydt = -ballSpeed
	}

	// with paddle
	if g.Ball.X >= g.Paddle.X && g.Ball.Y >= g.Paddle.Y && g.Ball.Y <= g.Paddle.Y + g.Paddle.H {
		g.Ball.dxdt = -g.Ball.dxdt
		g.score++
		if g.score > g.highScore {
			g.highScore = g.score
		}
	}
}
