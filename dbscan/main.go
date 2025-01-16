package main

import (
	"algos/drawer"
	"algos/tools"
	"fmt"
	"log"
	"math"
	"slices"

	"math/rand"

	"gonum.org/v1/plot/plotter"
)

type Point struct {
	N    int
	X, Y float64
}

// Distance calculates the Euclidean distance between two points.
func (p Point) Distance(q Point) float64 {
	return math.Sqrt(math.Pow(p.X-q.X, 2) + math.Pow(p.Y-q.Y, 2))
}

// DBSCAN performs the DBSCAN clustering algorithm.
func DBSCAN(points []Point, eps float64, minPts int) (clusters [][]Point) {
	visited := make([]bool, len(points))
	clusterID := 0
	clusters = make([][]Point, 0)

	for i := range points {
		if visited[i] {
			continue // Already visited
		}
		visited[i] = true

		neighbors := regionQuery(points, points[i], eps)
		if len(neighbors) < minPts {
			continue // Mark as noise, but we don't store it
		}

		clusterID++
		currentCluster := []Point{points[i]}
		clusters = append(clusters, currentCluster)

		for _, n := range neighbors {
			if !visited[n] {
				visited[n] = true
				newNeighbors := regionQuery(points, points[n], eps)
				if len(newNeighbors) >= minPts {
					neighbors = append(neighbors, newNeighbors...)
				}
			}
			if !contains(clusters[clusterID-1], points[n]) {
				clusters[clusterID-1] = append(clusters[clusterID-1], points[n])
			}
		}
	}

	return clusters
}

// regionQuery finds all points within the eps distance of the given point.
func regionQuery(points []Point, center Point, eps float64) []int {
	var neighbors []int
	for i := range points {
		if center.Distance(points[i]) <= eps {
			neighbors = append(neighbors, i)
		}
	}
	return neighbors
}

// contains checks if a point is in the cluster.
func contains(cluster []Point, p Point) bool {
	for _, point := range cluster {
		if point.X == p.X && point.Y == p.Y {
			return true
		}
	}
	return false
}

func main() {
	points := []Point{
		{0, 330, 1950}, // 0
		{1, 370, 2020}, // 1
		{2, 330, 1990}, // 2
		{3, 350, 2000}, // 3
		{4, 350, 2030}, // 4

		{5, 2040, 340}, // 5
		{6, 2050, 350}, // 6
		{7, 2020, 320}, // 7
		{8, 1960, 350}, // 8

		{9, 1200, 1800},  // 9
		{10, 1270, 1750}, // 10
		{11, 1340, 1730}, // 11
		{12, 1410, 1700}, // 12
		{13, 1494, 1680}, // 13
		{14, 1578, 1662}, // 14
		{15, 1600, 1640}, // 15
		{16, 1620, 1620}, // 16
		{17, 1640, 1600}, // 17
		{18, 1662, 1578}, // 18
		{19, 1680, 1494}, // 19
		{20, 1700, 1410}, // 20
		{21, 1730, 1340}, // 21
		{22, 1750, 1270}, // 22
		{23, 1800, 1200}, // 23

		{24, 300, 300}, // 24
		{25, 312, 290}, // 25
		{26, 302, 315}, // 26
		{27, 278, 255}, // 27
		{28, 700, 700}, // 28
		{29, 726, 702}, // 29
		{30, 666, 653}, // 30
		{31, 612, 623}, // 31
		{32, 400, 500}, // 32
		{33, 434, 561}, // 33
		{34, 322, 433}, // 34
		{35, 402, 441}, // 35
		{36, 355, 412}, // 36
		{37, 100, 700}, // 37
		{38, 32, 615},  // 38
		{39, 125, 670}, // 39
	}

	// points = genPoints(400, 3)

	wad := wadCalc(points)
	fmt.Printf("wad: %f\n", wad)

	minAvrWeightDst := minDistCalc(points)
	fmt.Printf("min dst: %f\n", minAvrWeightDst)

	// var eps float64 = minAvrWeightDst * math.Pi
	var eps float64 = minAvrWeightDst * 3
	fmt.Printf("eps: %f\n", eps)
	minPts := 3
	// minPts := len(points) / 10

	clusters := DBSCAN(points, eps, minPts)

	// for i, cluster := range clusters {
	// 	fmt.Printf("Cluster %d: ", i+1)
	// 	for _, point := range cluster {
	// 		fmt.Printf("%d(%.0f, %.0f) ", point.N, point.X, point.Y)
	// 	}
	// 	fmt.Println()
	// }
	for i, cluster := range clusters {
		fmt.Printf("Cluster %d: ", i+1)
		for _, point := range cluster {
			fmt.Printf("%d ", point.N)
		}
		fmt.Println()
	}

	clusterToMerge := convertToMerge(clusters)
	// Объединение кластеров
	mergedClusters2 := tools.MergeClusters(clusterToMerge)
	mergedClusters1 := tools.MergeClusters(mergedClusters2)
	mergedClusters := tools.MergeClusters(mergedClusters1)

	clusters2 := convertFromMerge(mergedClusters, points)
	fmt.Println()
	for i, cluster2 := range clusters2 {
		fmt.Printf("Cluster %d: ", i+1)
		for _, point := range cluster2 {
			fmt.Printf("%d ", point.N)
		}
		fmt.Println()
	}

	fmt.Printf("min dst: %f\n", minAvrWeightDst)
	fmt.Printf("eps: %f\n", eps)
	fmt.Printf("минимальное кол-о точек в кластере: %d\n", minPts)

	clstrsArray, total := convertToXYsArray(clusters2)

	fmt.Printf("Всего распределенных точек: %d\n", total)
	fmt.Printf("Всего нераспределенных точек: %d\n", len(points)-total)

	err := drawer.PlotClasters("outClustersDBSCAN.png", clstrsArray)
	if err != nil {
		log.Fatal(err.Error())
	}

	xys := convertToXYs(points)
	drawer.PlotData("outPlotLabels.png", xys)
}

func wadCalc(points []Point) float64 {
	var dstArr []float64
	for _, p1 := range points {
		for _, p2 := range points {
			if p1.N == p2.N {
				continue
			}
			dst := p1.Distance(p2)
			dstArr = append(dstArr, dst)
		}
	}
	slices.Sort(dstArr)
	return dstArr[len(dstArr)/2]
}

func minDistCalc(points []Point) float64 {
	var minDstArray []float64
	for _, p1 := range points {
		var minDst float64 = math.MaxFloat64
		for _, p2 := range points {
			if p1.N == p2.N {
				continue
			}
			dst := p1.Distance(p2)
			if dst < minDst {
				minDst = dst
			}
		}
		minDstArray = append(minDstArray, minDst)
	}
	slices.Sort(minDstArray)
	fmt.Printf("%.0f\n", minDstArray)
	return minDstArray[len(minDstArray)/2]
}

func convertFromMerge(mergedClusters []*tools.Cluster, points []Point) [][]Point {
	var result [][]Point
	for i, v := range mergedClusters {
		result = append(result, make([]Point, 0))
		for _, vv := range v.GetPoints() {
			result[i] = append(result[i], getPoint(vv, points))
		}
	}
	return result
}

func getPoint(id int, points []Point) Point {
	for i, v := range points {
		if v.N == id {
			return points[i]
		}
	}
	return Point{}
}

func convertToMerge(input [][]Point) []*tools.Cluster {
	var result []*tools.Cluster
	for _, v := range input {
		var temp []int
		for _, vv := range v {
			temp = append(temp, vv.N)
		}
		result = append(result, tools.NewCluster(temp))
	}

	return result
}

func convertToXYsArray(clstrs [][]Point) ([]plotter.XYs, int) {
	var clstrsArray []plotter.XYs
	total := 0
	for _, c := range clstrs {
		var temp plotter.XYs
		for _, p := range c {
			temp = append(temp, plotter.XY{
				X: p.X,
				Y: p.Y,
			})
			total++
		}
		clstrsArray = append(clstrsArray, temp)
	}
	return clstrsArray, total
}

func convertToXYs(pnts []Point) plotter.XYs {

	var temp plotter.XYs
	for _, p := range pnts {
		temp = append(temp, plotter.XY{
			X: p.X,
			Y: p.Y,
		})
	}
	return temp
}

func genPoints(quantity int, param int) []Point {
	cnt := 0
	var result []Point

	for i := 0; i < param; i++ {

		size := rand.Intn(500) + 500
		shift := rand.Intn(2000 - size)

		for cnt < (quantity/param)*(i+1) {
			x := float64(rand.Intn(size) + shift)
			y := float64(rand.Intn(size) + shift)
			result = append(result, Point{
				N: cnt,
				X: x,
				Y: y,
			})
			cnt++
		}
		fmt.Printf("i: %d\n", i)
	}
	fmt.Printf("result length: %d\n", len(result))
	return result
}
