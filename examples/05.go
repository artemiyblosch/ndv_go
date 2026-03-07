package main

import (
	"fmt"
	"ndv"
)

func e05() {
	p, err := ndv.ImportOFF("./examples/C10.off")
	if err != nil {
		fmt.Println(err)
		return
	}
	p_struct := ndv.CanonizedStruct(p)
	for i := 10; i < 11; i++ {
		d := make([]float64, 11)
		d[i] = 1
		p_struct = p_struct.Prismify(d)
	}
	p = ndv.FromStructure(ndv.NormalizedStruct(p_struct))
	//fmt.Println(p)
	p.ExportOFF("./examples/C11.off")
}
