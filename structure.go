package ndv

import (
	"fmt"
	"slices"
)

type Structure struct {
	Verticies    []Point
	RestElements [][][]int
}

func CanonizedStruct(p Polytope) Structure {
	result := p.Structure
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

func NormalizedStruct(structure Structure) Structure {
	result := structure

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

func recoverFace(face []int, edges [][]int) []int {
	result := []int{}
	fmt.Println(edges)

	faceGraph := make(map[int][]int)
	for _, edge_index := range face {
		begin, end := edges[int(edge_index)][0], edges[int(edge_index)][1]
		faceGraph[begin] = append(faceGraph[begin], end)
		faceGraph[end] = append(faceGraph[end], begin)
	}

	curVertex := edges[face[0]][0]
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
