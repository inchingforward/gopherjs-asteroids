package main

import (
	"math"

	"honnef.co/go/js/dom"
)

// A Ship is the player.
type Ship struct {
	x, y, dx, dy int
	dir          float64
}

// KeysPressed contains the current keys being pressed.
type KeysPressed struct {
	left, right, thrust bool
}

var (
	ship        *Ship
	canvas      *dom.HTMLCanvasElement
	ctx         *dom.CanvasRenderingContext2D
	keysPressed *KeysPressed
)

func (ship Ship) draw(ctx *dom.CanvasRenderingContext2D) {
	ctx.Save()
	ctx.Translate(ship.x, ship.y)

	ctx.Call("rotate", ship.dir)

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
	if keysPressed.right {
		ship.dir = ship.dir + math.Pi/16.0
	}

	if keysPressed.left {
		ship.dir = ship.dir + -math.Pi/16.0
	}
}

func update() {
	updateShip()
}

func draw() {
	ctx.FillStyle = "white"
	ctx.FillRect(0, 0, canvas.Width, canvas.Height)

	ship.draw(ctx)
}

func gameLoop() {
	update()
	draw()
	dom.GetWindow().SetTimeout(gameLoop, 1000/30)
}

func main() {
	window := dom.GetWindow()
	doc := window.Document()
	canvas = doc.GetElementByID("canvas").(*dom.HTMLCanvasElement)
	ctx = canvas.GetContext2d()

	keysPressed = &KeysPressed{false, false, false}

	window.AddEventListener("keydown", false, func(event dom.Event) {
		ke := event.(*dom.KeyboardEvent)

		keysPressed.left = ke.KeyCode == 65
		keysPressed.right = ke.KeyCode == 83
		keysPressed.thrust = ke.KeyCode == 75
	})

	window.AddEventListener("keyup", false, func(event dom.Event) {
		ke := event.(*dom.KeyboardEvent)

		switch ke.KeyCode {
		case 65:
			keysPressed.left = false
		case 83:
			keysPressed.right = false
		case 76:
			keysPressed.thrust = false
		}
	})

	ship = &Ship{canvas.Width / 2, canvas.Height / 2, 0, 0, 0}

	gameLoop()
}
