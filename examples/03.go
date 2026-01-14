package main

import (
	"fmt"
	"image/color"
	"ndv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func main() {
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

	p, err := ndv.ImportOFF("./examples/Cp.off")
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Resize(fyne.NewSize(500, 500))

	framerate := 60.
	go func() {
		t := time.NewTicker(time.Duration(1000/framerate) * time.Millisecond)
		defer t.Stop()
		for range t.C {
			p.Rotate([2]int{1, 3}, 0.025, []float64{0., 0., 0., 0.})
			p.Rotate([2]int{0, 2}, 0.05, []float64{0., 0., 0., 0.})
			world.Clear()
			p.Draw(world)

			fyne.DoAndWait(func() { w.SetContent(canvas.NewImageFromImage(world.Image)) })
		}
	}()
	//p.Draw(world)
	//w.SetContent(canvas.NewImageFromImage(world.Image))
	w.ShowAndRun()
}
