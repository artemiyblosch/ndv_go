package ndv

import (
	"math"
	"slices"
)

func Ortho(point Point) Point {
	return point[:len(point)-1]
}

func Iso(direction Point) Projection {
	return func(point Point) Point {
		coordinate := point[len(point)-1]
		point = point[:len(point)-1]

		for i := range point {
			if i > len(direction)-1 {
				continue
			}
			point[i] += coordinate * direction[i]
		}
		return point
	}
}

func TranslateFunc(direction Point) Projection {
	return func(point Point) Point {
		for i := range point {
			point[i] += direction[i]
		}
		return point
	}
}

func ScaleFunc(scale float64) Projection {
	return func(point Point) Point {
		for i := range point {
			point[i] *= scale
		}
		return point
	}
}

func Upcast(dim int) Projection {
	return func(point Point) Point {
		diff := dim - len(point)
		if diff <= 0 {
			return point
		}
		return slices.Grow(point, diff)[:dim]
	}
}

func Safe(projection Projection) Projection {
	return func(point Point) Point {
		if len(point) < 3 {
			return point
		}
		return projection(point)
	}
}

func Central(center Point) Projection {
	return func(point Point) Point {
		point_height := point[len(point)-1]
		center_height := center[len(center)-1]

		if point_height == center_height {
			for i := range point {
				point[i] = math.Copysign(math.Inf(1), float64(point[i]))
			}
			return point
		}

		point = point[:len(point)-1]
		if len(point) >= len(center) {
			center = append(
				make(Point, len(point)-len(center)+1),
				center...,
			)
		}

		for i := range point {
			point[i] = center_height * (point[i] - center[i]) / (center_height - point_height)
		}
		return point
	}
}

func RotateFunc(plane [2]int, angle float64, around Point) Projection {
	return func(point Point) Point {
		cos, sin := math.Cos(angle), math.Sin(angle)
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
