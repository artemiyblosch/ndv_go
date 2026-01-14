package ndv

import "golang.org/x/exp/constraints"

type Point = []float64

type Projection = func(Point) Point

type Number interface {
	constraints.Float | constraints.Integer
}
