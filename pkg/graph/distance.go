package graph

import (
	"github.com/epsxy/flower/pkg/model"
	"github.com/epsxy/flower/pkg/utils"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"gopkg.in/vmarkovtsev/go-lcss.v1"
)

// TODO: Remove naive implementation
func Optimize(vertexes []string, graph map[string][]string, distance model.DistanceNorm) [][]string {
	mid := len(vertexes) / 2
	left := vertexes[:mid]
	right := vertexes[mid:]

	affinityMatrix := _buildAffinityMatrix(vertexes, distance)

	var tempLeft []string = make([]string, len(left))
	var tempRight []string = make([]string, len(right))

	copy(tempLeft, left)
	copy(tempRight, right)

	var resLeft []string = make([]string, len(left))
	var resRight []string = make([]string, len(right))

	copy(resLeft, left)
	copy(resRight, right)

	currentWeight := Weight([][]string{left, right}, graph, affinityMatrix)
	for i, v := range left {
		for j, w := range right {
			tempLeft[i] = w
			tempRight[j] = v
			weight := Weight([][]string{tempLeft, tempRight}, graph, affinityMatrix)
			if weight > currentWeight {
				currentWeight = weight
				copy(resLeft, tempLeft)
				copy(resRight, tempRight)
			}
		}
	}
	return [][]string{resLeft, resRight}
}

func Weight(partitions [][]string, graph map[string][]string, affinityMatrix map[string]map[string]float64) float64 {
	count := float64(0)
	for _, partition := range partitions {
		count += _partitionWeight(partition, graph, affinityMatrix)
	}
	return count
}

// Naive weighting function:
// - If 2 neighbours are in the same partition we add 1
// - Otherwise 0
func _partitionWeight(partition []string, graph map[string][]string, affinity map[string]map[string]float64) float64 {
	count := float64(0)
	for i, v := range partition {
		for _, j := range graph[v] {
			if utils.ArrayContains(partition, j) {
				count++
			}
		}
		if i < len(partition)-2 {
			count += 2 * affinity[v][partition[i+1]]
		}
	}
	return count
}

func _buildAffinityMatrix(vertexes []string, distance model.DistanceNorm) map[string]map[string]float64 {
	res := map[string]map[string]float64{}
	for _, v := range vertexes {
		res[v] = map[string]float64{}
		for _, w := range vertexes {
			if v == w {
				res[v][w] = 0
			} else {
				res[v][w] = _wordWeight(v, w, distance)
			}
		}
	}
	return res
}

func _wordWeight(w1, w2 string, norm model.DistanceNorm) float64 {
	switch norm {
	case model.DistanceNormSubstring:
		return float64(len(w1)) / float64(len(string(lcss.LongestCommonSubstring([]byte(w1), []byte(w2)))))
	case model.DistanceNormLevenshtein:
		// TODO: add custom levenshtein options
		return float64(levenshtein.DistanceForStrings([]rune(w1), []rune(w2), levenshtein.DefaultOptions))
	default:
		return float64(len(w1)) / float64(len(string(lcss.LongestCommonSubstring([]byte(w1), []byte(w2)))))
	}
}
