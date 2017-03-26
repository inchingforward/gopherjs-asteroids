package main

import (
	"math"

	"github.com/oskca/gopherjs-canvas"
	"github.com/oskca/gopherjs-dom"
)

// A Ship is the player.
type Ship struct {
	x, y, dx, dy, dir float64
}

// KeysPressed contains the current keys being pressed.
type KeysPressed struct {
	left, right, thrust bool
}

const maxSpeed = 6.0

var (
	ship         *Ship
	cnvs         *canvas.Canvas
	ctx          *canvas.Context2D
	keysPressed  *KeysPressed
	canvasWidth  float64
	canvasHeight float64
)

func clamp(lo, hi, x float64) float64 {
	return math.Max(lo, math.Min(hi, x))
}

func (ship Ship) draw(ctx *canvas.Context2D) {
	ctx.Save()
	ctx.Translate(ship.x, ship.y)
	ctx.Rotate(ship.dir)
	ctx.StrokeStyle = "black"
	ctx.LineWidth = 1
	ctx.BeginPath()
	ctx.MoveTo(0, -1)
	ctx.MoveTo(0, -10)
	ctx.LineTo(-7, 10)
	ctx.LineTo(0, 8)
	ctx.LineTo(7, 10)
	ctx.LineTo(0, -10)
	ctx.Stroke()
	ctx.Restore()
}

func updateShip() {
	leftAdjust := 0.0
	if keysPressed.left {
		leftAdjust = (-math.Pi / 16.0)
	}

	rightAdjust := 0.0
	if keysPressed.right {
		rightAdjust = (math.Pi / 16.0)
	}

	ship.dir = ship.dir + leftAdjust + rightAdjust

	thrust := 0.0
	if keysPressed.thrust {
		thrust = 0.5
	}

	dx := clamp(-maxSpeed, maxSpeed, ship.dx+thrust*math.Cos(ship.dir-math.Pi/2.0))
	dy := clamp(-maxSpeed, maxSpeed, ship.dy+thrust*math.Sin(ship.dir-math.Pi/2.0))

	ship.dx = dx
	ship.dy = dy

	ship.x = math.Mod(canvasWidth+ship.x+ship.dx, canvasWidth)
	ship.y = math.Mod(canvasHeight+ship.y+ship.dy, canvasHeight)
}

func update() {
	updateShip()
}

func draw() {
	ctx.FillStyle = "white"
	ctx.FillRect(0.0, 0.0, canvasWidth, canvasHeight)

	ship.draw(ctx)
}

func gameLoop() {
	update()
	draw()
	dom.Window().Call("setTimeout", gameLoop, 1000/30)
}

func main() {
	window := dom.Window()
	doc := window.Document
	cnvs := canvas.New(doc.GetElementById("canvas").Object)
	canvasWidth = float64(cnvs.Width)
	canvasHeight = float64(cnvs.Height)
	ctx = cnvs.GetContext2D()

	keysPressed = &KeysPressed{false, false, false}

	window.AddEventListener("keydown", func(event *dom.Event) {
		keyCode := event.KeyCode

		keysPressed.left = keyCode == 65
		keysPressed.right = keyCode == 83
		keysPressed.thrust = keyCode == 75
	})

	window.AddEventListener("keyup", func(event *dom.Event) {
		keyCode := event.KeyCode

		switch keyCode {
		case 65:
			keysPressed.left = false
		case 83:
			keysPressed.right = false
		case 75:
			keysPressed.thrust = false
		}
	})

	ship = &Ship{canvasWidth / 2.0, canvasHeight / 2.0, 0.0, 0.0, 0.0}

	gameLoop()
}
