package ndv

import (
	"maps"
	"math"
)

func GetEdges(p Polytope) [][]int {
	edges := Structure{}
	for _, face := range p.Structure.RestElements[0] {
		for j := range face[:len(face)-1] {
			addEdge(edges, face, j, j+1)
		}
		addEdge(edges, face, 0, len(face)-1)
	}
	return edges.RestElements[0]
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

func Pyramidify(structure Structure, point Point) Structure {
	if len(structure.RestElements[len(structure.RestElements)-1]) > 1 {
		structure = AddMaximal(structure)
	}
	counts := GetCounts(structure)
	structure.Verticies = append(structure.Verticies, point)

	var d, nd map[int]int
	d = make(map[int]int)

	for vIndex := 0; vIndex < counts[0]; vIndex++ {
		d[vIndex] = len(structure.RestElements[0])
		structure.RestElements[0] = append(structure.RestElements[0], []int{vIndex, counts[0]})
	}

	var facet, ridge []int
	for dim := 0; dim < len(counts)-2; dim++ {
		nd = make(map[int]int)
		for index := 0; index < counts[dim]; index++ {
			ridge = structure.RestElements[dim][index]
			facet = []int{index}

			for _, peak := range ridge {
				facet = append(facet, d[int(peak)])
			}
			nd[index] = len(structure.RestElements[dim+1])
			structure.RestElements[dim+2] = append(structure.RestElements[dim+2], facet)
		}
		d = maps.Clone(nd)
	}

	return structure
}

func GetCounts(structure Structure) []int {
	result := make([]int, len(structure.RestElements)+1)
	result[0] = len(structure.Verticies)
	for i, v := range structure.RestElements {
		result[i] = len(v)
	}
	return result
}

func addEdge(result Structure, face []int, begin, end int) int {
	if face[end] > face[begin] {
		begin, end = end, begin
	}
	I := index(result.RestElements[1], []int{face[begin], face[end]})
	if I == -1 {
		result.RestElements[0] = append(result.RestElements[0], []int{face[begin], face[end]})
		return len(result.RestElements[0]) - 1
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
