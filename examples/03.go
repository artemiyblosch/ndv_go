package main

import (
	"fmt"
	"image/color"
	"ndv"

	"fyne.io/fyne/v2/app"
)

func e03() {
	a := app.New()
	w := a.NewWindow("Preivew")

	world := ndv.NewWorld(
		[]ndv.Projection{
			ndv.Central([]float64{3.}),
			ndv.Central([]float64{3.}),
		},
		color.NRGBA{255, 255, 255, 255},
		[2][2]float64{{-2, 2}, {-2, 2}},
		[2]int{500, 500},
	)

	p, err := ndv.ImportOFF("./examples/It.off")
	//p = p.Centrate()
	if err != nil {
		fmt.Println(err)
		return
	}
	preview := ndv.Preview{
		Window:    w,
		Polytope:  p,
		Framerate: 60,
		World:     world,
	}
	preview.Start()
	//p.Draw(world)
	//w.SetContent(canvas.NewImageFromImage(world.Image))
	w.ShowAndRun()
}
