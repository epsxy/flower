package graph

import (
	"github.com/epsxy/flower/pkg/model"
	"github.com/epsxy/flower/pkg/utils"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"gopkg.in/vmarkovtsev/go-lcss.v1"
)

func Split(vertexes []string, graph map[string][]string, options *model.UMLTreeOptions) [][]string {
	var splitVertexes [][]string
	// trivial case, we can return the full vertexes list
	if len(vertexes) < options.MaxPartitionSize {
		return [][]string{vertexes}
	}

	currentRes := []string{}
	for _, v := range vertexes {
		if len(currentRes) == options.MaxPartitionSize {
			splitVertexes = append(splitVertexes, currentRes)
			currentRes = []string{}
		}
		currentRes = append(currentRes, v)
	}
	splitVertexes = append(splitVertexes, currentRes)

	// affinity matrix
	affinityMatrix := buildAffinityMatrix(vertexes, options.DistanceNorm)

	// optimize each partition
	for i := range splitVertexes {
		for j := range splitVertexes {
			if i >= j {
				continue
			}
			splitVertexes[i], splitVertexes[j] = ReArrangePartitions(splitVertexes[i], splitVertexes[j], graph, affinityMatrix, options)
		}
	}

	// return
	return splitVertexes
}

func ReArrangePartitions(p1, p2 []string, graph map[string][]string, affinityMatrix map[string]map[string]float64, options *model.UMLTreeOptions) ([]string, []string) {
	var tempP1 []string = make([]string, len(p1))
	var tempP2 []string = make([]string, len(p2))

	copy(tempP1, p1)
	copy(tempP2, p2)

	var resP1 []string = make([]string, len(p1))
	var resP2 []string = make([]string, len(p2))

	copy(resP1, p1)
	copy(resP2, p2)

	currentWeight := Weight([][]string{p1, p2}, graph, affinityMatrix, options)
	for i := range p1 {
		for j := range p2 {
			temp := tempP1[i]
			tempP1[i] = tempP2[j]
			tempP2[j] = temp
			weight := Weight([][]string{tempP1, tempP2}, graph, affinityMatrix, options)
			if weight > currentWeight {
				currentWeight = weight
				copy(resP1, tempP1)
				copy(resP2, tempP2)
			} else {
				tempP2[j] = tempP1[i]
				tempP1[i] = temp
			}
		}
	}
	return resP1, resP2
}

func Weight(partitions [][]string, graph map[string][]string, affinityMatrix map[string]map[string]float64, options *model.UMLTreeOptions) float64 {
	count := float64(0)
	for _, partition := range partitions {
		count += partitionWeight(partition, graph, affinityMatrix, options)
	}
	return count
}

// Naive weighting function:
// With k1, k2 defined in in the configuration (TODO)
// - If 2 neighbours are in the same partition we add k1 × 1 (default)
// - We also add k2 × 1 (default) the affinity ratio
// - Otherwise 0
func partitionWeight(partition []string, graph map[string][]string, affinity map[string]map[string]float64, options *model.UMLTreeOptions) float64 {
	count := float64(0)
	for i, v := range partition {
		for _, j := range graph[v] {
			if utils.ArrayContains(partition, j) {
				count += float64(options.WeightEdge) * 1
			}
		}
		if i < len(partition)-1 {
			count += float64(options.WeightDistance) * affinity[v][partition[i+1]]
		}
	}
	return count
}

func buildAffinityMatrix(vertexes []string, distance model.DistanceNorm) map[string]map[string]float64 {
	res := map[string]map[string]float64{}
	for _, v := range vertexes {
		res[v] = map[string]float64{}
		for _, w := range vertexes {
			if v == w {
				res[v][w] = 0
			} else {
				res[v][w] = wordWeight(v, w, distance)
			}
		}
	}
	return res
}

func wordWeight(w1, w2 string, norm model.DistanceNorm) float64 {
	var d float64
	switch norm {
	case model.DistanceNormSubstring:
		d1 := float64(len(string(lcss.LongestCommonSubstring([]byte(w1), []byte(w2))))) / float64(len(w1))
		d2 := float64(len(string(lcss.LongestCommonSubstring([]byte(w1), []byte(w2))))) / float64(len(w2))
		d = (d1 + d2) / 2
	case model.DistanceNormLevenshtein:
		// TODO: handle custom levenshtein options
		d = 1 / float64(levenshtein.DistanceForStrings([]rune(w1), []rune(w2), levenshtein.DefaultOptions))
	default:
		d1 := float64(len(string(lcss.LongestCommonSubstring([]byte(w1), []byte(w2))))) / float64(len(w1))
		d2 := float64(len(string(lcss.LongestCommonSubstring([]byte(w1), []byte(w2))))) / float64(len(w2))
		d = (d1 + d2) / 2
	}
	return utils.RoundFloat(d, 2)
}
