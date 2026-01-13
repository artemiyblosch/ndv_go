package ndv

import "golang.org/x/exp/constraints"

type Point[T interface{}] = []T

type Projection[T interface{}] = func(Point[T]) Point[T]

type Number interface {
	constraints.Float | constraints.Integer
}
