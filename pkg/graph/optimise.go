package graph

import "github.com/epsxy/flower/pkg/utils"

// TODO: Remove naive implementation
func Optimize(vertexes []string, input []string, graph map[string][]string) [][]string {
	mid := len(vertexes) / 2
	left := vertexes[:mid]
	right := vertexes[mid:]

	var tempLeft []string = make([]string, len(left))
	var tempRight []string = make([]string, len(right))

	copy(tempLeft, left)
	copy(tempRight, right)

	var resLeft []string = make([]string, len(left))
	var resRight []string = make([]string, len(right))

	copy(resLeft, left)
	copy(resRight, right)

	currentWeight := Weight([][]string{left, right}, graph)
	for i, v := range left {
		for j, w := range right {
			tempLeft[i] = w
			tempRight[j] = v
			weight := Weight([][]string{tempLeft, tempRight}, graph)
			if weight > currentWeight {
				currentWeight = weight
				copy(resLeft, tempLeft)
				copy(resRight, tempRight)
			}
		}
	}
	return [][]string{resLeft, resRight}
}

func Weight(partitions [][]string, graph map[string][]string) int {
	count := 0
	for _, partition := range partitions {
		count += _partitionWeight(partition, graph)
	}
	return count
}

// Naive weighting function:
// - If 2 neighbours are in the same partition we add 1
// - Otherwise 0
func _partitionWeight(partition []string, graph map[string][]string) int {
	count := 0
	for _, v := range partition {
		for _, j := range graph[v] {
			if utils.ArrayContains(partition, j) {
				count++
			}
		}
	}
	return count
}
