package main

import (
	"algos/drawer"
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"gonum.org/v1/plot/plotter"
)

var points = plotter.XYs{
	plotter.XY{X: 350, Y: 2000},

	plotter.XY{X: 2050, Y: 350},

	plotter.XY{X: 300, Y: 300},
	plotter.XY{X: 312, Y: 290},
	plotter.XY{X: 302, Y: 315},
	plotter.XY{X: 278, Y: 255},
	plotter.XY{X: 267, Y: 333},

	plotter.XY{X: 700, Y: 700},
	plotter.XY{X: 726, Y: 702},
	plotter.XY{X: 666, Y: 653},
	plotter.XY{X: 612, Y: 623},

	plotter.XY{X: 400, Y: 500},
	plotter.XY{X: 434, Y: 561},
	plotter.XY{X: 322, Y: 433},
	plotter.XY{X: 402, Y: 441},
	plotter.XY{X: 355, Y: 412},

	plotter.XY{X: 100, Y: 700},
	plotter.XY{X: 32, Y: 615},
	plotter.XY{X: 125, Y: 670},
}

type dest struct {
	from int
	to   int
	dest float64
}

type byDest []dest

func (d byDest) Len() int {
	return len(d)
}

func (d byDest) Less(i, j int) bool {
	return d[i].dest < d[j].dest
}

func (d byDest) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type byFrom []dest

func (d byFrom) Len() int {
	return len(d)
}

func (d byFrom) Less(i, j int) bool {
	return d[i].from < d[j].from
}

func (d byFrom) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func main() {
	var dst byDest
	for i := 0; i < len(points)-1; i++ {
		for j := i + 1; j < len(points); j++ {
			dst = append(dst, dest{i, j, distance(points[i], points[j])})
			fmt.Printf(
				"%d-%d: %f\n",
				dst[len(dst)-1].from,
				dst[len(dst)-1].to,
				dst[len(dst)-1].dest,
			)
		}
	}
	//	fr := 11
	// point, ds := maxDest(points, points[fr])
	//	fmt.Printf("\nmax dist from %d to %d: %f\n", fr, point, ds)

	printMaxDest(points)
	sort.Sort(dst)

	fmt.Printf("\nSorted data:\n")
	for _, v := range dst {
		fmt.Printf(
			"%d-%d: %f\n",
			v.from,
			v.to,
			v.dest,
		)
	}

	mediana := dst[len(dst)/2].dest
	fmt.Printf("\nmediana: %f\n", mediana)

	dstMed := byFrom(dst[:len(dst)/2])
	sort.Sort(dstMed)
	fmt.Printf("\nMedian data:\n")
	for _, v := range dstMed {
		fmt.Printf(
			"%d-%d: %f\n",
			v.from,
			v.to,
			v.dest,
		)
	}
	// store:=make(map[int][]dest)

	d := convert(points)
	// km := kmeans.New()
	km, err := kmeans.NewWithOptions(0.001, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	clstrs, err := km.Partition(d, 4)
	if err != nil {
		log.Fatal(err.Error())
	}

	var clstrsArray []plotter.XYs
	for _, c := range clstrs {
		fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
		fmt.Printf("Matching data points: %+v\n\n", c.Observations)
		var temp plotter.XYs
		for i := 0; i < len(c.Observations); i++ {
			temp = append(temp, plotter.XY{
				X: c.Observations[i].Coordinates()[0],
				Y: c.Observations[i].Coordinates()[1],
			})
		}
		clstrsArray = append(clstrsArray, temp)
	}
	err = drawer.PlotClasters("outClusters.png", clstrsArray)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = drawer.PlotData("outPlots.png", points)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = drawer.PlotPolygon("outpol.png", points)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func convert(data plotter.XYs) clusters.Observations {
	var d clusters.Observations
	for _, val := range data {
		d = append(d, clusters.Coordinates{
			val.X,
			val.Y,
		})
	}
	return d
}

func distance(one, two plotter.XY) float64 {
	x := math.Abs(one.X - two.X)
	y := math.Abs(one.Y - two.Y)
	dest := math.Sqrt(x + x + y*y)
	return dest
}

func printMaxDest(ptrs plotter.XYs) {
	fmt.Println()
	for i, p := range ptrs {
		var max float64 = 0
		var to = 0
		for j, d := range ptrs {
			if distance(p, d) > max {
				max = distance(p, d)
				to = j
			}
		}
		fmt.Printf("Max distance from %d: to %d at %f\n", i, to, max)
	}
}

func maxDest(ptrs plotter.XYs, pt plotter.XY) (int, float64) {
	var max float64 = 0
	var to = 0
	for j, d := range ptrs {
		if distance(pt, d) > max {
			max = distance(pt, d)
			to = j
		}
	}
	return to, max
}

// func findStart(points plotter.XYs, clastNumb int) {
// 	var nextPt int
// 	startPoint := rand.Intn(len(points))
// 	for i := 0; i < clastNumb; i++ {
// 		nextPt, dist := maxDest(points, points[startPoint])
// 		for j, v := range points {
// 			oned := distance(points[startPoint], points[j])
// 			twod := distance(points[nextPt], points[j])
// 		}
// 	}
// }
