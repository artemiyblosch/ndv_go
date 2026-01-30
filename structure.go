package ndv

import (
	"fmt"
	"slices"
)

type Structure struct {
	Verticies    []Point
	RestElements [][][]int
}

func (s Structure) String() string {
	result := fmt.Sprintf("Verticies:\n%v\n\n\nRest:\n", s.Verticies)

	for _, i := range s.RestElements {
		result += fmt.Sprintf("%v\n\n", i)
	}
	return result
}

func CanonizedStruct(p Polytope) Structure {
	result := p.Structure
	result.RestElements = slices.Insert(result.RestElements, 0, [][]int{})

	tmpFace := []int{}

	for i, face := range result.RestElements[1] {
		tmpFace = []int{}
		for j := range face[:len(face)-1] {
			tmpFace = append(tmpFace, addEdge(&result.RestElements[0], face[j], face[j+1]))
		}
		tmpFace = append(tmpFace, addEdge(&result.RestElements[0], face[0], face[len(face)-1]))
		result.RestElements[1][i] = tmpFace
	}

	if len(p.Structure.RestElements) == 2 && len(p.Structure.RestElements[1]) == 1 {
		result.RestElements = slices.Delete(result.RestElements, 1, 2)
	}
	return result
}

func (p Structure) Rotate(plane [2]int, angle float64, around Point) Structure {
	return p.VertexMap(RotateFunc(plane, angle, around))
}

func (p Structure) Translate(direction Point) Structure {
	return p.VertexMap(TranslateFunc(direction))
}

func (p Structure) Scale(scale float64) Structure {
	return p.VertexMap(ScaleFunc(scale))
}
func (p Structure) Upcast(dim int) Structure {
	return p.VertexMap(Upcast(dim))
}

func (p Structure) VertexMap(f Projection) Structure {
	for i, v := range p.Verticies {
		p.Verticies[i] = f(v)
	}
	return p
}
func NormalizedStruct(s Structure) Structure {
	faces := make([][]int, len(s.RestElements[1]))

	for i, face := range s.RestElements[1] {
		faces[i] = recoverFace(face, s.RestElements[0])
	}

	s.RestElements[1] = faces
	s.RestElements = slices.Delete(s.RestElements, 0, 1)
	return s
}

func recoverFace(face []int, edges [][]int) []int {
	face_edges := make([][]int, len(face))
	for i := range face {
		face_edges[i] = edges[face[i]]
	}

	been := []int{face_edges[0][0]}
	edge_graph := ConstructEdgeGraph(face_edges)
	done := false
	for !done {
		done = true
		for _, neighbour := range edge_graph[been[len(been)-1]] {
			if slices.Contains(been, neighbour) {
				continue
			}
			been = append(been, neighbour)
			done = false
			break
		}
	}
	return been
}

func ConstructEdgeGraph(edges [][]int) map[int][]int {
	result := make(map[int][]int)
	for _, i := range edges {
		result[i[0]] = append(result[i[0]], i[1])
		result[i[1]] = append(result[i[1]], i[0])
	}
	return result
}
