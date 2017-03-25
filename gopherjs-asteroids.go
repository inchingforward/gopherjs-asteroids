package main

import (
	"honnef.co/go/js/dom"
)

var (
	ship   = Ship{}
	canvas *dom.HTMLCanvasElement
	ctx    *dom.CanvasRenderingContext2D
)

// A Ship is the player.
type Ship struct {
	x int
	y int
}

func (ship Ship) draw(ctx *dom.CanvasRenderingContext2D) {
	ctx.StrokeStyle = "black"
	ctx.LineWidth = 1
	ctx.Translate(ship.x, ship.y)
	ctx.MoveTo(0, -1)
	ctx.MoveTo(0, -10)
	ctx.LineTo(-7, 10)
	ctx.LineTo(0, 8)
	ctx.LineTo(7, 10)
	ctx.LineTo(0, -10)
	ctx.Stroke()
}

func update() {

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
	doc := dom.GetWindow().Document()
	canvas = doc.GetElementByID("canvas").(*dom.HTMLCanvasElement)
	ctx = canvas.GetContext2d()

	ship.x = canvas.Width / 2
	ship.y = canvas.Height / 2

	gameLoop()
}
