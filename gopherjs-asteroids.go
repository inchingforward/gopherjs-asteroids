package main

import (
	"honnef.co/go/js/dom"
)

func main() {
	doc := dom.GetWindow().Document()
	canvas := doc.GetElementByID("canvas").(*dom.HTMLCanvasElement)
	ctx := canvas.GetContext2d()

	ctx.FillStyle = "black"
	ctx.StrokeStyle = "black"

	ctx.BeginPath()
	ctx.MoveTo(150, 150)
	ctx.LineTo(250, 250)
	ctx.LineTo(50, 250)
	ctx.LineTo(150, 150)
	ctx.Fill()
}
