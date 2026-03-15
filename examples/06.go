package main

import (
	"fmt"
	"ndv"
)

func main() {
	p1 := ndv.CanonizedStruct(ndv.Polygon(4, 1))
	p2 := ndv.CanonizedStruct(ndv.Polygon(3, 1))
	fmt.Println(ndv.Product(p1, p2))
}
