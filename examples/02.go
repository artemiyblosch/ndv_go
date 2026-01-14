package main

import (
	"fmt"
	"image/color"
	"ndv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func e_02() {
	a := app.New()
	w := a.NewWindow("Preivew")

	world := ndv.NewWorld(
		[]ndv.Projection{
			ndv.Central([]float64{1.}),
		},
		color.NRGBA{255, 255, 255, 255},
		[2][2]float64{{-5, 5}, {-5, 5}},
		[2]int{500, 500},
	)

	p, err := ndv.ImportOFF("./examples/Cube.off")
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Resize(fyne.NewSize(500, 500))
	p.Draw(world)
	w.SetContent(canvas.NewImageFromImage(world.Image))
	w.ShowAndRun()
}
