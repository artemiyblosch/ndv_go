package ndv

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/llgcode/draw2d/draw2dimg"
)

type World[T Number] struct {
	Projections []Projection[T]
	Background  color.NRGBA
	Domain      [2][2]float64
	Image       *image.RGBA
}

func NewWorld[T Number](
	projections []Projection[T],
	bg color.NRGBA,
	domain [2][2]float64,
	size [2]int,
) World[T] {
	img := image.NewRGBA(image.Rect(0, 0, size[0], size[1]))
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{0, 0}, draw.Src)

	return World[T]{projections, bg, domain, img}
}

func (world *World[T]) Clear() {
	img := image.NewRGBA(world.Image.Rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{world.Background}, image.Point{0, 0}, draw.Src)
	world.Image = img
}

func (world World[T]) Project(point Point[T]) Point[T] {
	result := make(Point[T], len(point))
	copy(result, point)

	for _, projection := range world.Projections {
		result = projection(result)
	}
	return result[:2]
}

func (world World[float64]) DomainToPixelCoords(point Point[float64]) [2]float64 {
	return [2]float64{
		(point[0] - float64(world.Domain[0][0])) /
			(float64(world.Domain[0][1]) - float64(world.Domain[0][0])) *
			float64(world.Image.Bounds().Dx()),

		(point[1] - float64(world.Domain[1][0])) /
			(float64(world.Domain[1][1]) - float64(world.Domain[1][0])) *
			float64(world.Image.Bounds().Dy()),
	}
}

func (world World[float64]) DrawFace(points [][]float64, faceColor color.NRGBA) {
	p_coords := make([][2]float64, len(points))
	for i, point := range points {
		p_coords[i] = world.DomainToPixelCoords(world.Project(point))
	}
	gc := draw2dimg.NewGraphicContext(world.Image)

	gc.SetFillColor(faceColor)
	faceColor.A = 255
	gc.SetStrokeColor(faceColor)
	gc.SetLineWidth(1)
	gc.MoveTo(float64(p_coords[0][0]), float64(p_coords[0][1]))
	for _, p := range p_coords {
		gc.LineTo(p[0], p[1])
	}
	gc.Close()
	gc.FillStroke()
}
