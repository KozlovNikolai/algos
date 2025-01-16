package drawer

import (
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

func PlotClasters(path string, clstrsArray []plotter.XYs) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create %s: %v", path, err)
	}

	p := plot.New()

	for _, clst := range clstrsArray {
		sc, err := plotter.NewScatter(clst)
		if err != nil {
			return fmt.Errorf("could not create scatter: %v", err)
		}
		sc.GlyphStyle.Shape = draw.BoxGlyph{}

		sc.Color = color.RGBA{
			R: uint8(rand.Uint32() / 4),
			G: uint8(rand.Uint32() / 4),
			B: uint8(rand.Uint32() / 4),
			A: 255,
		}
		p.Add(sc)
	}

	wt, err := p.WriterTo(512, 512, "png")
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}
	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf("could not write to %s: %v", path, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close %s: %v", path, err)
	}
	return nil
}

func PlotData(path string, xys plotter.XYs) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create %s: %v", path, err)
	}

	p := plot.New()

	sp, err := plotter.NewScatter(xys)
	if err != nil {
		return fmt.Errorf("could not create scatter: %v", err)
	}
	sp.GlyphStyle.Shape = draw.CrossGlyph{}
	sp.Color = color.RGBA{R: 255, A: 255}
	p.Add(sp)

	var lbs = make([]string, len(xys))
	for i := range xys {
		lbs[i] = strconv.Itoa(i)
	}
	labels, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    xys,
		Labels: lbs,
	})
	if err != nil {
		log.Fatalf("could not creates labels plotter: %+v", err)
	}
	p.Add(labels)

	wt, err := p.WriterTo(512, 512, "png")
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}
	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf("could not write to %s: %v", path, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close %s: %v", path, err)
	}
	return nil
}

func PlotPolygon(path string, xyer plotter.XYer) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create %s: %v", path, err)
	}

	p := plot.New()

	s, err := plotter.NewPolygon(xyer)
	if err != nil {
		return fmt.Errorf("could not create scatter: %v", err)
	}

	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)

	wt, err := p.WriterTo(256, 256, "png")
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}
	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf("could not write to %s: %v", path, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close %s: %v", path, err)
	}
	return nil
}
