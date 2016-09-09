package main

import (
	"image/color"
	"math/rand"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
)

func draw() {
	rand.Seed(int64(0))
	n := 1
	bubbleData := randomTriples(n)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Bubbles"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	bs, err := plotter.NewBubbles(bubbleData, vg.Points(1), vg.Points(20))
	if err != nil {
		panic(err)
	}
	bs.Color = color.RGBA{R: 196, B: 128, A: 255}
	p.Add(bs)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "1234567.png"); err != nil {
		panic(err)
	}
}

// randomTriples returns some random x, y, z triples
// with some interesting kind of trend.
func randomTriples(n int) plotter.XYZs {
	data := make(plotter.XYZs, n)
	for i := range data {
		if i == 0 {
			data[i].X = rand.Float64()
		} else {
			data[i].X = data[i-1].X + 2*rand.Float64()
		}
		data[i].Y = data[i].X + 10*rand.Float64()
		data[i].Z = data[i].X
	}
	return data
}
