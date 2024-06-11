package helpers

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// https://www.w3schools.com/css/css_colors_rgb.asp

func plodder() {
	RED := color.NRGBA{R: 255, A: 255}
	//GREEN := color.NRGBA{G: 255, A: 255}
	BLUE := color.NRGBA{B: 255, A: 255}
	GREY := color.RGBA{180, 180, 180, 255}
	BLACK := color.Black
	rnd := rand.New(rand.NewSource(1))

	// randomPoints returns some random x, y points
	// with some interesting kind of trend.
	randomPoints := func(n int) plotter.XYs {
		month := 1
		const (
			day  = 1
			hour = 1
			min  = 1
			sec  = 1
			nsec = 1
		)
		pts := make(plotter.XYs, n)
		for ix := range pts {
			date := time.Date(2024, time.Month(month+ix), day, hour, min, sec, nsec, time.UTC).Unix()
			pts[ix].X = float64(date)
			pts[ix].Y = 100.0 * rnd.Float64()
		}
		return pts
	}

	n := 6

	p := plot.New()
	p.Title.Text = "NV Polling"
	p.X.Tick.Marker = plot.TimeTicks{Format: "Jan 02"}
	p.Y.Label.Text = "Voter %"
	p.Add(plotter.NewGrid())

	data := randomPoints(n)
	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}
	line.Color = BLUE
	line.Width = 2
	points.Shape = draw.CircleGlyph{}
	points.Color = BLACK
	p.Add(line, points)

	data = randomPoints(n)
	line, points, err = plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}
	line.Color = RED
	line.Width = 2
	points.Shape = draw.CircleGlyph{}
	points.Color = BLACK
	p.Add(line, points)

	data = randomPoints(n)
	line, points, err = plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}
	line.Color = GREY
	line.Width = 2
	points.Shape = draw.CircleGlyph{}
	points.Color = BLACK
	p.Add(line, points)

	err = p.Save(10*vg.Centimeter, 5*vg.Centimeter, "xx.png")
	if err != nil {
		log.Panic(err)
	}
}
