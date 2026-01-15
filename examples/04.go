package main

import (
	"fmt"
	"ndv"
)

func main() {
	p, err := ndv.ImportOFF("./examples/Cube.off")
	if err != nil {
		fmt.Println(err)
		return
	}
	p_struct := ndv.CanonizedStruct(p.Upcast(4))
	p_struct = p_struct.Pyramidify([]float64{0., 0., 0., 0.5})

	fmt.Println(p_struct)
	p = ndv.FromStructure(ndv.NormalizedStruct(p_struct))
	fmt.Println(p)
	p.ExportOFF("./examples/Cp.off")
}
