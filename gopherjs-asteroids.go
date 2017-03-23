package main

import (
	"honnef.co/go/js/dom"
)

func main() {
	doc := dom.GetWindow().Document()

	elem := doc.CreateElement("p")
	elem.SetInnerHTML("Hello, world!")

	doc.GetElementsByTagName("body")[0].AppendChild(elem)
}
