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

func NormalizedStruct(s Structure) Structure {
	faces := make([][]int, len(s.RestElements[1]))
	edgeGraph := ConstructEdgeGraph(s.RestElements[0])
	for i, face := range s.RestElements[1] {
		faces[i] = recoverFace(face, edgeGraph)
	}
	s.RestElements[1] = faces
	s.RestElements = slices.Delete(s.RestElements, 0, 1)
	return s
}

func recoverFace(face []int, edgeGraph map[int][]int) []int {
	//been :=
	return []int{}
}

func ConstructEdgeGraph(edges [][]int) map[int][]int {
	result := make(map[int][]int)
	for _, i := range edges {
		result[i[0]] = append(result[i[0]], i[1])
		result[i[1]] = append(result[i[1]], i[0])
	}
	return result
}
