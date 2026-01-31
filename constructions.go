package ndv

import (
	"maps"
	"math"
	"slices"
)

func GetEdges(p Polytope) [][]int {
	return [][]int{}
}

func AddMaximal(p Structure) Structure {
	maximal := make([]int, len(p.RestElements[len(p.RestElements)-1]))
	for i := range p.RestElements[len(p.RestElements)-1] {
		maximal[i] = i
	}
	p.RestElements = append(p.RestElements, [][]int{maximal})
	return p
}

func Polygon(sides_num int, sides_denom int) Polytope {
	structure := Structure{}
	structure.Verticies = make([][]float64, sides_num)
	structure.RestElements[0] = make([][]int, sides_denom)

	angle := 2 * math.Pi / float64(sides_num)

	for i := range structure.Verticies {
		structure.Verticies[i] = []float64{math.Cos(angle * float64(i)), math.Sin(angle * float64(i))}
	}

	var vertex int
	var face []int
	for i := 0; i < sides_denom; i++ {
		face = []int{}
		vertex = i

		face = append(face, vertex)
		vertex += sides_denom
		vertex %= sides_num

		for vertex != i {
			face = append(face, vertex)
			vertex += sides_denom
			vertex %= sides_num
		}
	}

	return FromStructure(structure)
}

func (s Structure) Pyramidify(point Point) Structure {
	s = s.Upcast(len(point))
	counts := GetCounts(s)
	s.Verticies = append(s.Verticies, point)
	if len(s.RestElements[len(s.RestElements)-1]) != 1 {
		s = AddMaximal(s)
	}

	d := make(map[int]int)
	for i := 0; i < counts[0]; i++ {
		d[i] = len(s.RestElements[0])
		s.RestElements[0] = append(s.RestElements[0], []int{i, counts[0]})
	}
	for dim := 1; dim < len(counts); dim++ {
		nd := make(map[int]int)
		for element := 0; element < counts[dim]; element++ {
			tmpNext := []int{element}
			for _, ridge := range s.RestElements[dim-1][element] {
				tmpNext = append(tmpNext, d[ridge])
			}
			nd[element] = len(s.RestElements[dim])
			s.RestElements[dim] = append(s.RestElements[dim], tmpNext)
		}
		d = maps.Clone(nd)
	}

	return s
}

func shift(facet []int, shift int) []int {
	f := slices.Clone(facet)
	for i := range facet {
		f[i] = f[i] + shift
	}
	return f
}

func (s Structure) Prismify(direction Point) Structure {
	s = s.Upcast(len(direction))
	counts := GetCounts(s)

	for i := range s.Verticies {
		s.Verticies = append(s.Verticies, TranslateFunc(direction)(slices.Clone(s.Verticies[i])))
	}

	if len(s.RestElements[len(s.RestElements)-1]) != 1 {
		s = AddMaximal(s)
	}

	for i := range s.RestElements[0] {
		s.RestElements[0] = append(s.RestElements[0], shift(s.RestElements[0][i], counts[0]))
	}

	for dim := 1; dim < len(counts); dim++ {
		for i := range s.RestElements[dim] {
			s.RestElements[dim] = append(s.RestElements[dim], shift(s.RestElements[dim][i], counts[dim]))
		}
	}
	d := make(map[int]int)
	for i := 0; i < counts[0]; i++ {
		d[i] = len(s.RestElements[0])
		s.RestElements[0] = append(s.RestElements[0], []int{i, counts[0] + i})
	}
	for dim := 1; dim < len(counts); dim++ {
		nd := make(map[int]int)
		for element := 0; element < counts[dim]; element++ {
			tmpNext := []int{element, counts[dim] + element}
			for _, ridge := range s.RestElements[dim-1][element] {
				tmpNext = append(tmpNext, d[ridge])
			}
			nd[element] = len(s.RestElements[dim])
			s.RestElements[dim] = append(s.RestElements[dim], tmpNext)
		}
		d = maps.Clone(nd)
	}

	return s
}

func GetCounts(structure Structure) []int {
	result := make([]int, len(structure.RestElements)+1)
	result[0] = len(structure.Verticies)
	for i, v := range structure.RestElements {
		result[i+1] = len(v)
	}
	return result
}

func addEdge(edges *[][]int, begin, end int) int {
	if begin > end {
		begin, end = end, begin
	}
	I := index(*edges, []int{begin, end})
	if I == -1 {
		*edges = append(*edges, []int{begin, end})
		return len(*edges) - 1
	}
	return I
}

func index(elems [][]int, value []int) int {
	var eq bool
	for i, elem := range elems {
		eq = false
		if len(elem) != len(value) {
			continue
		}

		for j, c1 := range elem {
			if c1 != value[j] {
				eq = false
				break
			}
			eq = true
		}
		if eq {
			return i
		}
	}
	return -1
}
