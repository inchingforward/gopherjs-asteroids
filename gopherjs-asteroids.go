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

// SaveContext pushes the current state onto the drawing stack.
func SaveContext(ctx *dom.CanvasRenderingContext2D) {
	ctx.Call("save")
}

// RestoreContext restores the most recently saved canvas state
// from the drawing stack.
func RestoreContext(ctx *dom.CanvasRenderingContext2D) {
	ctx.Call("restore")
}

func (ship Ship) draw(ctx *dom.CanvasRenderingContext2D) {
	SaveContext(ctx)
	ctx.Translate(ship.x, ship.y)
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
	RestoreContext(ctx)
}

func update() {

}

func draw() {
	ctx.FillStyle = "white"
	ctx.FillRect(0, 0, canvas.Width, canvas.Height)

	ship.draw(ctx)
}

func gameLoop() {
	//js.Debugger()
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
