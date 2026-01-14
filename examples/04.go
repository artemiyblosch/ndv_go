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
	p_struct = ndv.Pyramidify(p_struct, []float64{0., 0., 0., 0.5})

	p = ndv.FromStructure(ndv.NormalizedStruct(p_struct))
	p.ExportOFF("./examples/Cp.off")
}
