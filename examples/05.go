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
	p_struct := ndv.CanonizedStruct(p)
	p_struct = p_struct.Prismify([]float64{0., 0., 0., 1})

	p = ndv.FromStructure(ndv.NormalizedStruct(p_struct))
	//fmt.Println(p)
	p.ExportOFF("./examples/Ts.off")
}
