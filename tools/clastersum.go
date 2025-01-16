package tools

import (
	"sort"
)

// Cluster представляет собой структуру для хранения кластера.
type Cluster struct {
	Points map[int]struct{}
}

// NewCluster создает новый кластер.
func NewCluster(points []int) *Cluster {
	c := &Cluster{Points: make(map[int]struct{})}
	for _, p := range points {
		c.Points[p] = struct{}{}
	}
	return c
}

// Merge объединяет текущий кластер с другим кластером.
func (c *Cluster) Merge(other *Cluster) {
	for point := range other.Points {
		c.Points[point] = struct{}{}
	}
}

// HasIntersection проверяет, имеет ли кластер пересечения с другим кластером.
func (c *Cluster) HasIntersection(other *Cluster) bool {
	for point := range c.Points {
		if _, found := other.Points[point]; found {
			return true
		}
	}
	return false
}

// GetPoints возвращает отсортированный срез точек из кластера.
func (c *Cluster) GetPoints() []int {
	points := make([]int, 0, len(c.Points))
	for point := range c.Points {
		points = append(points, point)
	}
	sort.Ints(points)
	return points
}

// MergeClusters объединяет пересекающиеся кластеры.
func MergeClusters(clusters []*Cluster) []*Cluster {
	mergedClusters := make([]*Cluster, 0)

	for _, cluster := range clusters {
		merged := false
		for _, mCluster := range mergedClusters {
			if cluster.HasIntersection(mCluster) {
				mCluster.Merge(cluster)
				merged = true
				break
			}
		}
		if !merged {
			mergedClusters = append(mergedClusters, cluster)
		}
	}

	return mergedClusters
}
