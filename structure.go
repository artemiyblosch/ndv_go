package ndv

import "slices"

type Structure[T Number] struct {
	Verticies    []Point[T]
	RestElements [][][]int
}

func CanonizedStruct[T Number](p Polytope[T]) Structure[T] {
	result := Structure[T]{}
	copy(result.Verticies, p.Structure.Verticies)
	copy(result.RestElements, p.Structure.RestElements)
	result.RestElements = slices.Insert(result.RestElements, 0, [][]int{})

	tmpFace := []int{}
	for i, face := range result.RestElements[1] {
		tmpFace = []int{}
		for j := range face[:len(face)-1] {
			tmpFace = append(tmpFace, addEdge(result, face, j, j+1))
		}
		tmpFace = append(tmpFace, addEdge(result, face, 0, len(face)-1))
		result.RestElements[1][i] = tmpFace
	}

	if len(p.Structure.RestElements) == 2 && len(p.Structure.RestElements[1]) == 1 {
		result.RestElements = slices.Delete(result.RestElements, 1, 2)
	}
	return result
}

func NormalizedStruct[T Number](structure Structure[T]) Structure[T] {
	result := Structure[T]{}
	copy(result.Verticies, structure.Verticies)
	copy(result.RestElements, structure.RestElements)
	if len(result.RestElements) == 1 {
		result = AddMaximal(result)
	}

	tmpFace := []int{}
	for i, face := range result.RestElements[1] {
		tmpFace = recoverFace(face, result.RestElements[1])
		result.RestElements[1][i] = tmpFace
	}

	result.RestElements = slices.Delete(result.RestElements, 0, 1)

	return result
}

func recoverFace[T Number](face Point[T], edges [][]int) []int {
	result := []int{}

	faceGraph := make(map[int][]int)
	for _, edge_index := range face {
		begin, end := edges[int(edge_index)][0], edges[int(edge_index)][1]
		faceGraph[begin] = append(faceGraph[begin], end)
		faceGraph[end] = append(faceGraph[end], begin)
	}

	curVertex := edges[int(face[0])][0]
	been := []int{}
	found := true

	for found {
		been = append(been, curVertex)
		result = append(result, curVertex)
		found = false
		for _, vertexAdj := range faceGraph[curVertex] {
			for _, beenVertex := range been {
				if beenVertex == vertexAdj {
					goto Been
				}
			}
			curVertex = vertexAdj
			found = true
		Been:
		}
	}

	return result
}
