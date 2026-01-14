package main

import (
	"fmt"
	"image/color"
	"math"
	"ndv"
)

func e00() {
	structure := ndv.Structure{
		Verticies:    [][]float64{{0., 0.}, {0., 1.}, {1., 0.}, {1., 1.}},
		RestElements: [][][]int{{{0, 1, 2, 3}}},
	}

	p := ndv.Polytope{
		Structure: structure,
		Colors:    []color.NRGBA{{0, 0, 0, 255}},
	}

	p = p.Rotate([2]int{0, 1}, math.Pi/2, []float64{0., 0.})
	fmt.Printf("%v\n", p)
}
