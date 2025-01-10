package main

import (
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	p := plot.New()
	p.Title.Text = "Labels"
	p.X.Min = -10
	p.X.Max = +10
	p.Y.Min = 0
	p.Y.Max = +20

	labels, err := plotter.NewLabels(plotter.XYLabels{
		XYs: []plotter.XY{
			{X: -5, Y: 5},
			{X: +5, Y: 5},
			{X: +5, Y: 15},
			{X: -5, Y: 15},
		},
		Labels: []string{"A", "B", "C", "D"},
	})
	if err != nil {
		log.Fatalf("could not creates labels plotter: %+v", err)
	}

	p.Add(labels)

	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "labels.png")
	if err != nil {
		log.Fatalf("could save plot: %+v", err)
	}
}
