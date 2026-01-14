package main

import (
	"image/color"
	"ndv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func e_01() {
	a := app.New()
	w := a.NewWindow("Preivew")

	world := ndv.NewWorld(
		[]ndv.Projection{ndv.Iso([]float64{1., 1.})},
		color.NRGBA{255, 255, 255, 255},
		[2][2]float64{{-5, 5}, {-5, 5}},
		[2]int{500, 500},
	)
	world.DrawFace([][]float64{
		{0.5, 0.5, 0.5},
		{0.5, 0.5, -0.5},
		{0.5, -0.5, -0.5},
		{0.5, -0.5, 0.5},
	}, color.NRGBA{0, 0, 0, 125})
	w.Resize(fyne.NewSize(500, 500))
	w.SetContent(canvas.NewImageFromImage(world.Image))
	w.ShowAndRun()
}
