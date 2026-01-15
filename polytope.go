package ndv

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Polytope struct {
	Structure  Structure
	Colors     []color.NRGBA
	FaceColors []color.NRGBA
}

func (p Polytope) String() string {
	representation := "0:\n"

	for _, v := range p.Structure.Verticies {
		representation += fmt.Sprintf("%v ", v)
	}

	for i, v := range p.Structure.RestElements {
		representation += fmt.Sprintf("\n%v:\n%v\n", i, v)
	}

	return representation
}

func (p Polytope) Rotate(plane [2]int, angle float64, around Point) Polytope {
	return p.VertexMap(RotateFunc(plane, angle, around))
}

func (p Polytope) Translate(direction Point) Polytope {
	return p.VertexMap(TranslateFunc(direction))
}

func (p Polytope) Scale(scale float64) Polytope {
	return p.VertexMap(ScaleFunc(scale))
}
func (p Polytope) Upcast(dim int) Polytope {
	return p.VertexMap(Upcast(dim))
}

func (p Polytope) VertexMap(f Projection) Polytope {
	for i, v := range p.Structure.Verticies {
		p.Structure.Verticies[i] = f(v)
	}
	return p
}

func (p Polytope) Center() Point {
	center := make(Point, len(p.Structure.Verticies[0]))
	for i, c_coord := range center {
		for _, p_coord := range p.Structure.Verticies {
			c_coord += p_coord[i]
		}
		c_coord /= float64(len(p.Structure.Verticies))
	}
	return center
}

func (p Polytope) ComprisingsOf(facet_index, dimension int) []int {
	facets := []int{facet_index}
	for dim := len(p.Structure.RestElements) - 1; dim >= dimension; dim-- {
		ridges := make([]bool, len(p.Structure.RestElements[dim-1]))
		for _, facet := range facets {
			for _, i := range p.Structure.RestElements[dim][facet] {
				ridges[int(i)] = true
			}
		}
		facets = []int{}
		for i, in := range ridges {
			if !in {
				continue
			}
			facets = append(facets, i)
		}
	}

	return facets
}

func (p *Polytope) UpdateFaceColors() {
	p.FaceColors = make([]color.NRGBA, len(p.Structure.RestElements[0]))
	var faces []int
	var faceColor color.NRGBA

	for facet_index := range p.Structure.RestElements[len(p.Structure.RestElements)-1] {
		faces = p.ComprisingsOf(facet_index, 1)
		faceColor = p.Colors[facet_index]
		for _, i := range faces {
			p.FaceColors[i] = combineColors(p.FaceColors[i], faceColor)
		}
	}
}

func combineColors(c1, c2 color.NRGBA) color.NRGBA {
	if c1.A == 0 {
		return c2
	}
	if c2.A == 0 {
		return c1
	}

	var newColor color.NRGBA
	newColor.A = 125
	newColor.R = (c1.R + c2.R) / 2
	newColor.G = (c1.G + c2.G) / 2
	newColor.B = (c1.B + c2.B) / 2
	return newColor
}

func (p Polytope) Draw(world World) {

	var face [][]float64
	for face_index, struct_face := range p.Structure.RestElements[0] {
		face = [][]float64{}
		for _, p_index := range struct_face {
			face = append(face, p.Structure.Verticies[int(p_index)])
		}

		world.DrawFace(face, p.FaceColors[face_index])
	}
}

func scan(scanner *bufio.Scanner) {
	scanner.Scan()
	for scanner.Text() == "" || scanner.Text()[0] == '#' {
		scanner.Scan()
	}
}

func spilttedText(scanner *bufio.Scanner) []string {
	return strings.Split(strings.TrimSpace(scanner.Text()), " ")
}

func FromStructure(structure Structure) Polytope {
	return Polytope{
		Structure:  structure,
		Colors:     make([]color.NRGBA, len(structure.RestElements[len(structure.RestElements)-1])),
		FaceColors: make([]color.NRGBA, len(structure.RestElements[len(structure.RestElements)-1])),
	}
}

func ImportOFF(path string) (Polytope, error) {
	file, err := os.Open(path)
	if err != nil {
		return Polytope{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	scanner.Scan()
	countsStr := spilttedText(scanner)
	counts := make([]int, len(countsStr))

	for i, v := range countsStr {
		counts[i], err = strconv.Atoi(v)
		if err != nil {
			return Polytope{}, err
		}
	}
	counts = slices.Delete(counts, 2, 3)

	structure := Structure{}
	structure.RestElements = make([][][]int, len(counts)-1)
	colors := make([]color.NRGBA, counts[len(counts)-1])

	for i := 0; i < counts[0]; i++ {
		scan(scanner)
		point := spilttedText(scanner)

		vertex := make(Point, len(point))
		for c_index, coord := range point {
			c, err := strconv.ParseFloat(coord, 64)
			if err != nil {
				return Polytope{}, err
			}
			vertex[c_index] = c
		}
		structure.Verticies = append(structure.Verticies, vertex)
	}

	var alpha uint8
	for dim, j := range counts[1:] {
		for i := 0; i < j; i++ {
			scan(scanner)
			point := spilttedText(scanner)

			l, err := strconv.Atoi(point[0])
			if err != nil {
				return Polytope{}, err
			}

			element := make([]int, l)
			for c_index, coord := range point[1 : l+1] {
				c, err := strconv.ParseInt(coord, 10, 64)
				if err != nil {
					return Polytope{}, err
				}
				element[c_index] = int(c)
			}
			structure.RestElements[dim] = append(structure.RestElements[dim], element)

			if dim != len(counts)-2 {
				continue
			}

			col := make([]uint8, 4)
			for col_i, col_c := range point[l+1:] {
				col_ci, err := strconv.ParseUint(col_c, 10, 0)
				if err != nil {
					return Polytope{}, err
				}

				col[col_i] = uint8(col_ci)
			}
			if col[3] == 0 && col[0] == col[1] && col[1] == col[2] && col[2] == col[3] {
				colors[i] = color.NRGBA{}
				continue
			}

			alpha = 125
			if col[3] != 0 {
				alpha = col[3]
			}
			colors[i] = color.NRGBA{col[0], col[1], col[2], alpha}
		}
	}
	p := Polytope{Structure: structure, Colors: colors}
	p.UpdateFaceColors()
	return p, nil
}

func (p Polytope) ExportOFF(path string) {
	dim := len(p.Structure.Verticies[0])
	counts := GetCounts(p.Structure)
	counts = slices.Insert(counts, 2, len(GetEdges(p)))

	write := fmt.Sprint(counts)
	write = write[1 : len(write)-1]
	file_contents := []byte(fmt.Sprintf("%vOFF\n%v\n\n", dim, write))

	for _, elem := range p.Structure.Verticies {
		write = fmt.Sprint(elem)
		write = write[1:len(write)-1] + "\n"
		file_contents = slices.Concat(file_contents, []byte(write))
	}

	for _, elems := range p.Structure.RestElements {
		file_contents = append(file_contents, '\n')
		for _, elem := range elems {
			write = fmt.Sprint(elem)
			write = write[1 : len(write)-1]
			write = fmt.Sprintln(len(elem), write)
			file_contents = slices.Concat(file_contents, []byte(write))
		}
	}

	os.WriteFile(path, file_contents, 0644)

}
