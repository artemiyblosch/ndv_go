package ndv

import (
	"math"
	"slices"
)

func Ortho[T Number](point Point[T]) Point[T] {
	return point[:len(point)-1]
}

func Iso[T Number](direction Point[T]) Projection[T] {
	return func(point Point[T]) Point[T] {
		coordinate := point[len(point)-1]
		point = point[:len(point)-1]

		for i := range point {
			point[i] += coordinate * direction[i]
		}
		return point
	}
}

func TranslateFunc[T Number](direction Point[T]) Projection[T] {
	return func(point Point[T]) Point[T] {
		for i := range point {
			point[i] += direction[i]
		}
		return point
	}
}

func ScaleFunc[T Number](scale T) Projection[T] {
	return func(point Point[T]) Point[T] {
		for i := range point {
			point[i] *= scale
		}
		return point
	}
}

func Upcast[T Number](dim int) Projection[T] {
	return func(point Point[T]) Point[T] {
		diff := dim - len(point)
		if diff <= 0 {
			return point
		}
		return slices.Grow(point, diff)[:dim]
	}
}

func Safe[T Number](projection Projection[T]) Projection[T] {
	return func(point Point[T]) Point[T] {
		if len(point) < 3 {
			return point
		}
		return projection(point)
	}
}

func Central[T Number](center Point[T]) Projection[T] {
	return func(point Point[T]) Point[T] {
		point_height := point[len(point)-1]
		center_height := center[len(center)-1]

		if point_height == center_height {
			for i := range point {
				point[i] = T(math.Copysign(math.Inf(1), float64(point[i])))
			}
			return point
		}

		point = point[:len(point)-1]
		if len(point) >= len(center) {
			center = append(
				make(Point[T], len(point)-len(center)+1),
				center...,
			)
		}

		for i := range point {
			point[i] = center_height * (point[i] - center[i]) / (center_height - point_height)
		}
		return point
	}
}

func RotateFunc[T Number](plane [2]int, angle float64, around Point[T]) Projection[T] {
	return func(point Point[T]) Point[T] {
		cos, sin := T(math.Cos(angle)), T(math.Sin(angle))
		x, y := plane[0], plane[1]

		for i, v := range around {
			point[i] -= v
		}

		point[x], point[y] = cos*point[x]+sin*point[y], cos*point[y]-sin*point[x]

		for i, v := range around {
			point[i] += v
		}

		return point
	}
}
