package main

import (
	"fmt"

	"honnef.co/go/js/dom"
)

func update() {
	fmt.Println("update")
}

func draw() {
	fmt.Println("draw")
}

func gameLoop() {
	update()
	draw()
	dom.GetWindow().SetTimeout(gameLoop, 1000/30)
}

func main() {
	gameLoop()
}
