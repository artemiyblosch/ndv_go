package ndv

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Preview struct {
	Window    fyne.Window
	World     World
	Polytope  Polytope
	Framerate int
}

func (p Preview) Start() {
	p.Polytope.UpdateFaceColors()
	p.Window.Resize(fyne.NewSize(float32(p.World.Image.Rect.Size().X), float32(p.World.Image.Rect.Size().Y)))
	dim := len(p.Polytope.Structure.Verticies[0])
	space := []int{0., 1., 2.}
	entering_space := -1
	angular_speed := 0.005 * 2 * math.Pi

	if dim == 2 {
		p.Window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
			switch key.Name {
			case "Up":
				p.Polytope.Rotate([2]int{space[0], space[1]}, angular_speed, []float64{})
			case "Down":
				p.Polytope.Rotate([2]int{space[0], space[1]}, -angular_speed, []float64{})
			}
		})
	} else {
		p.Window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
			switch key.Name {
			case "Up":
				p.Polytope.Rotate([2]int{space[1], space[2]}, angular_speed, []float64{})
			case "Down":
				p.Polytope.Rotate([2]int{space[1], space[2]}, -angular_speed, []float64{})
			case "Left":
				p.Polytope.Rotate([2]int{space[0], space[2]}, angular_speed, []float64{})
			case "Right":
				p.Polytope.Rotate([2]int{space[0], space[2]}, -angular_speed, []float64{})
			case "LeftShift":
				entering_space = 0
			case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
				if entering_space < 0 {
					return
				}
				num, err := strconv.ParseInt(string(key.Name), 10, 4)
				if err != nil {
					return
				}
				if int(num) >= dim {
					return
				}
				space[entering_space] = int(num)

				entering_space += 1
				if entering_space > 2 {
					fmt.Println(space)
					entering_space = -1
				}
			}
		})
	}

	go func() {
		t := time.NewTicker(time.Duration(1000/p.Framerate) * time.Millisecond)
		defer t.Stop()
		for range t.C {
			p.World.Clear()
			p.Polytope.Draw(p.World, true)

			fyne.DoAndWait(func() {
				p.Window.SetContent(canvas.NewImageFromImage(p.World.Image))
			})
		}
	}()
}
