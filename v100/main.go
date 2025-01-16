package main

import (
	"algos/drawer"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot/plotter"
)

type Point struct {
	X, Y float64
}

// Функция для расчета расстояния между двумя точками
func distance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}

// Функция для отброса выбросов
func filterOutliers(points []Point) []Point {
	var xs, ys []float64
	for _, p := range points {
		xs = append(xs, p.X)
		ys = append(ys, p.Y)
	}

	meanX, meanY := stat.Mean(xs, nil), stat.Mean(ys, nil)
	stdDevX, stdDevY := stat.StdDev(xs, nil), stat.StdDev(ys, nil)

	var filteredPoints []Point
	for _, p := range points {
		if math.Abs(p.X-meanX) <= 2*stdDevX && math.Abs(p.Y-meanY) <= 2*stdDevY {
			filteredPoints = append(filteredPoints, p)
		}
	}

	return filteredPoints
}

// Алгоритм K-средних
func kMeans(points []Point, k int) (map[int][]Point, []Point) {
	// Инициализация центроидов
	centroids := make([]Point, k)
	rand.Seed(time.Now().UnixNano())
	for i := range centroids {
		centroids[i] = points[rand.Intn(len(points))]
	}

	clusters := make(map[int][]Point)
	for {
		// Очистка кластеров
		for i := range clusters {
			clusters[i] = nil
		}

		// Присвоение точек к ближайшим центроидам
		for _, p := range points {
			closestIndex := -1
			closestDistance := math.MaxFloat64
			for i, c := range centroids {
				dist := distance(p, c)
				if dist < closestDistance {
					closestDistance = dist
					closestIndex = i
				}
			}
			clusters[closestIndex] = append(clusters[closestIndex], p)
		}

		// Обновление центроидов
		newCentroids := make([]Point, k)
		for i := 0; i < k; i++ {
			var xSum, ySum float64
			for _, p := range clusters[i] {
				xSum += p.X
				ySum += p.Y
			}
			if len(clusters[i]) > 0 {
				newCentroids[i] = Point{
					X: xSum / float64(len(clusters[i])),
					Y: ySum / float64(len(clusters[i])),
				}
			} else {
				newCentroids[i] = centroids[i] // если кластер пустой, оставить прежний центроид
			}
		}

		// Проверка на сходимость
		converged := true
		for i := range centroids {
			if centroids[i] != newCentroids[i] {
				converged = false
				break
			}
		}
		if converged {
			break
		}
		centroids = newCentroids
	}
	return clusters, centroids
}

func main() {
	// Исходные данные
	points := []Point{
		{350, 2000}, // выброс
		{2050, 350}, // выброс
		{300, 300},
		{312, 290},
		{302, 315},
		{278, 255},
		{700, 700},
		{726, 702},
		{666, 653},
		{612, 623},
		{400, 500},
		{434, 561},
		{322, 433},
		{402, 441},
		{355, 412},
		{100, 700},
		{32, 615},
		{125, 670},
	}

	// Отбор выбросов
	filteredPoints := filterOutliers(points)

	// Количество кластеров
	k := 4
	clusters, centroids := kMeans(filteredPoints, k)

	// Вывод результатов
	fmt.Printf("Центроиды кластеров:\n")
	for i, c := range centroids {
		fmt.Printf("Кластер %d: (%.2f, %.2f)\n", i, c.X, c.Y)
	}
	fmt.Printf("\nКластеры:\n")
	for i, cluster := range clusters {
		fmt.Printf("Кластер %d: ", i)
		for _, p := range cluster {
			fmt.Printf("(%.2f, %.2f) ", p.X, p.Y)
		}
		fmt.Println()
	}

	clstrsArray := convertToXYsArray(clusters)

	err := drawer.PlotClasters("outClustersV100.png", clstrsArray)
	if err != nil {
		log.Fatal(err.Error())
	}

	// err = drawer.PlotData("outPlots.png", points)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
}

func convertToXYsArray(clstrs map[int][]Point) []plotter.XYs {
	var clstrsArray []plotter.XYs
	for _, c := range clstrs {
		var temp plotter.XYs
		for _, p := range c {
			temp = append(temp, plotter.XY{
				X: p.X,
				Y: p.Y,
			})
		}
		clstrsArray = append(clstrsArray, temp)
	}
	return clstrsArray
}
