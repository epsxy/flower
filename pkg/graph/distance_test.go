package graph

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/epsxy/flower/pkg/model"
	"github.com/epsxy/flower/pkg/utils"
	"github.com/stretchr/testify/require"
)

// TODO: implement missing tests
// func Test_Split(t *testing.T) {

// }

func Test_ReArrangePartitions(t *testing.T) {
	cases := map[string]struct {
		p1        []string
		p2        []string
		graph     map[string][]string
		affinity  map[string]map[string]float64
		expectedX []string
		expectedY []string
	}{
		"nominal": {
			p1: []string{
				"table_a1", "table_b1", "table_c1", "table_d1", "table_e1", "table_f1", "table_g1", "table_h1", "table_i1", "table_j1", "table_k1", "table_l1",
			},
			p2: []string{
				"table_a2", "table_b2", "table_c2", "table_d2", "table_e2", "table_f2", "table_g2", "table_h2", "table_i2", "table_j2", "table_k2", "table_l2",
			},
			// MermaidJs chart:
			// flowchart LR
			// table_a1((A1))
			// table_a2((A2))
			// table_b1((B1))
			// table_b2((B2))
			// table_c1((C1))
			// table_c2((C2))
			// table_d1((D1))
			// table_d2((D2))
			// table_e1((E1))
			// table_e2((E2))
			// table_f1((F1))
			// table_f2((F2))
			// table_g1((G1))
			// table_g2((G2))
			// table_h1((H1))
			// table_h2((H2))
			// table_i1((I1))
			// table_i2((I2))
			// table_j1((J1))
			// table_j2((J2))
			// table_k1((K1))
			// table_k2((K2))
			// table_l1((L1))
			// table_l2((L2))
			// table_a1 --> table_b1 --> table_a2 --> table_a1
			// table_a1 --> table_c1 --> table_c2 --> table_d1
			// table_a1 --> table_d1
			// table_b2 --> table_a1
			// table_d2 --> table_b2
			// table_e1 --> table_b2
			// table_e2 --> table_b2
			// table_f1 --> table_a1
			// table_e2 --> table_f1
			// table_f2 --> table_g1 --> table_g2
			// table_h2 --> table_h1
			// table_i2 --> table_h1
			// table_j2 --> table_h1
			// table_k2 --> table_h1
			// table_h1 --> table_g1
			// table_g2 --> table_d1
			// table_i1 --> table_h2
			// table_j1 --> table_e2
			// table_j1 --> table_k2
			// table_k1 --> table_j1
			// table_l1 --> table_j1
			// table_l2 --> table_j1
			graph: map[string][]string{
				"table_a1": {"table_b1", "table_a2", "table_c1", "table_d1", "table_f1", "table_b2"}, // OK
				"table_a2": {"table_b1", "table_a1"},                                                 // OK
				"table_b1": {"table_a1", "table_a2"},                                                 // OK
				"table_b2": {"table_d2", "table_e1", "table_e2", "table_a1"},                         // OK
				"table_c1": {"table_a1", "table_c2"},                                                 // OK
				"table_c2": {"table_c1", "table_d1"},                                                 // OK
				"table_d1": {"table_c2", "table_a1", "table_g2"},                                     // OK
				"table_d2": {"table_b2"},                                                             // OK
				"table_e1": {"table_b2"},                                                             // OK
				"table_e2": {"table_b2", "table_f1", "table_j1"},                                     // OK
				"table_f1": {"table_e2", "table_a1"},                                                 // OK
				"table_f2": {"table_g1"},                                                             // OK
				"table_g1": {"table_h1", "table_f2", "table_g2"},                                     // OK
				"table_g2": {"table_g1", "table_d1"},                                                 // OK
				"table_h1": {"table_k2", "table_j2", "table_i2", "table_h2", "table_g1"},             // OK
				"table_h2": {"table_h1", "table_i1"},                                                 // OK
				"table_i1": {"table_h2"},                                                             // OK
				"table_i2": {"table_h1"},                                                             // OK
				"table_j1": {"table_k1", "table_l1", "table_l2", "table_k2", "table_e2"},             // OK
				"table_j2": {"table_h1"},                                                             // OK
				"table_k1": {"table_j1"},                                                             // OK
				"table_k2": {"table_j1", "table_h1"},                                                 // OK
				"table_l1": {"table_j1"},                                                             // OK
				"table_l2": {"table_j1"},                                                             // OK
			},
			expectedX: []string{"table_c2", "table_g2", "table_h2", "table_k2", "table_i2", "table_f2", "table_g1", "table_h1", "table_l2", "table_j1", "table_d1", "table_l1"},
			expectedY: []string{"table_f1", "table_a2", "table_b2", "table_b1", "table_a1", "table_d2", "table_e1", "table_e2", "table_c1", "table_k1", "table_j2", "table_i1"},
		},
	}
	for name, c := range cases {
		x, y := ReArrangePartitions(c.p1, c.p2, c.graph, _buildFakeAffinityMap(c.p1, c.p2))
		require.Equal(t, x, c.expectedX, name)
		require.Equal(t, y, c.expectedY, name)
		// verify no value from the original partitions was lost
		for _, v := range c.p1 {
			containsX := utils.ArrayContains(x, v)
			containsY := utils.ArrayContains(y, v)
			require.True(t, containsX != containsY, fmt.Sprintf("element `%s` was in no array or in both arrays", v))
		}
		for _, v := range c.p2 {
			containsX := utils.ArrayContains(x, v)
			containsY := utils.ArrayContains(y, v)
			require.True(t, containsX != containsY, fmt.Sprintf("element `%s` was in no array or in both arrays", v))
		}
	}
}

func Test_Weight(t *testing.T) {
	cases := map[string]struct {
		partition [][]string
		graph     map[string][]string
		affinity  map[string]map[string]float64
		expected  float64
	}{
		"nominal": {
			partition: [][]string{
				{"table_a", "table_a_tags"},
				{"table_b", "table_b_tags"},
				{"table_a_table_b"},
			},
			graph: map[string][]string{
				"table_a":         {"table_a_tags", "table_a_table_b"},
				"table_b":         {"table_b_tags", "table_a_table_b"},
				"table_a_table_b": {"table_a", "table_b"},
				"table_a_tags":    {"table_a"},
				"table_b_tags":    {"table_b"},
			},
			affinity: map[string]map[string]float64{
				"table_a": {
					"table_a":         float64(0),
					"table_b":         float64(0.86),
					"table_a_table_b": float64(1),
					"table_a_tags":    float64(1),
					"table_b_tags":    float64(0.86),
				},
				"table_b": {
					"table_a":         float64(0.86),
					"table_b":         float64(0),
					"table_a_table_b": float64(1),
					"table_a_tags":    float64(0.86),
					"table_b_tags":    float64(1),
				},
				"table_a_table_b": {
					"table_a":         float64(0.47),
					"table_b":         float64(0.47),
					"table_a_table_b": float64(0),
					"table_a_tags":    float64(0.67),
					"table_b_tags":    float64(0.47),
				},
				"table_a_tags": {
					"table_a":         float64(0.58),
					"table_b":         float64(0.5),
					"table_a_table_b": float64(0.83),
					"table_a_tags":    float64(0),
					"table_b_tags":    float64(0.5),
				},
				"table_b_tags": {
					"table_a":         float64(0.5),
					"table_b":         float64(0.58),
					"table_a_table_b": float64(0.58),
					"table_a_tags":    float64(0.5),
					"table_b_tags":    float64(0),
				},
			},
			expected: float64(6),
		},
	}
	for name, c := range cases {
		require.Equal(t, Weight(c.partition, c.graph, c.affinity), c.expected, name)
	}
}

func Test_partitionWeight(t *testing.T) {
	cases := map[string]struct {
		partition []string
		graph     map[string][]string
		affinity  map[string]map[string]float64
		expected  float64
	}{
		"nominal": {
			partition: []string{"table_a", "table_b", "table_a_table_b", "table_a_tags", "table_b_tags"},
			graph: map[string][]string{
				"table_a":         {"table_a_tags", "table_a_table_b"},
				"table_b":         {"table_b_tags", "table_a_table_b"},
				"table_a_table_b": {"table_a", "table_b"},
				"table_a_tags":    {"table_a"},
				"table_b_tags":    {"table_b"},
			},
			affinity: map[string]map[string]float64{
				"table_a": {
					"table_a":         float64(0),
					"table_b":         float64(0.86),
					"table_a_table_b": float64(1),
					"table_a_tags":    float64(1),
					"table_b_tags":    float64(0.86),
				},
				"table_b": {
					"table_a":         float64(0.86),
					"table_b":         float64(0),
					"table_a_table_b": float64(1),
					"table_a_tags":    float64(0.86),
					"table_b_tags":    float64(1),
				},
				"table_a_table_b": {
					"table_a":         float64(0.47),
					"table_b":         float64(0.47),
					"table_a_table_b": float64(0),
					"table_a_tags":    float64(0.67),
					"table_b_tags":    float64(0.47),
				},
				"table_a_tags": {
					"table_a":         float64(0.58),
					"table_b":         float64(0.5),
					"table_a_table_b": float64(0.83),
					"table_a_tags":    float64(0),
					"table_b_tags":    float64(0.5),
				},
				"table_b_tags": {
					"table_a":         float64(0.5),
					"table_b":         float64(0.58),
					"table_a_table_b": float64(0.58),
					"table_a_tags":    float64(0.5),
					"table_b_tags":    float64(0),
				},
			},
			expected: float64(11.03),
		},
	}
	for name, c := range cases {
		require.Equal(t, partitionWeight(c.partition, c.graph, c.affinity), c.expected, name)
	}
}

func Test_buildAffinityMatrix(t *testing.T) {
	cases := map[string]struct {
		vertexes []string
		norm     model.DistanceNorm
		expected map[string]map[string]float64
	}{
		"substring: full match": {
			vertexes: []string{"table_a", "table_b", "table_a_table_b", "table_a_tags", "table_b_tags"},
			norm:     model.DistanceNormLevenshtein,
			expected: map[string]map[string]float64{
				"table_a": {
					"table_a":         float64(0),
					"table_b":         float64(0.5),
					"table_a_table_b": float64(0.13),
					"table_a_tags":    float64(0.2),
					"table_b_tags":    float64(0.2),
				},
				"table_b": {
					"table_a":         float64(0.5),
					"table_a_table_b": float64(0.13),
					"table_a_tags":    float64(0.14),
					"table_b":         float64(0),
					"table_b_tags":    float64(0.2),
				},
				"table_a_table_b": {
					"table_a":         float64(0.13),
					"table_b":         float64(0.13),
					"table_a_table_b": float64(0),
					"table_a_tags":    float64(0.14),
					"table_b_tags":    float64(0.11),
				},
				"table_a_tags": {
					"table_a":         float64(0.2),
					"table_b":         float64(0.14),
					"table_a_table_b": float64(0.14),
					"table_a_tags":    float64(0),
					"table_b_tags":    float64(0.5),
				},
				"table_b_tags": {
					"table_a":         float64(0.2),
					"table_a_table_b": float64(0.11),
					"table_a_tags":    float64(0.5),
					"table_b":         float64(0.2),
					"table_b_tags":    float64(0),
				},
			},
		},
	}
	for name, c := range cases {
		require.Equal(t, c.expected, buildAffinityMatrix(c.vertexes, c.norm), name)
	}
}

func Test_wordWeight(t *testing.T) {
	cases := map[string]struct {
		word1    string
		word2    string
		norm     model.DistanceNorm
		expected float64
	}{
		"substring: full match": {
			word1:    "here_is_a_string",
			word2:    "here_is_a_string_suffixed",
			norm:     model.DistanceNormSubstring,
			expected: float64(0.82),
		},
		"substring: reversed full match": {
			word1:    "here_is_a_string_suffixed",
			word2:    "here_is_a_string",
			norm:     model.DistanceNormSubstring,
			expected: float64(0.82),
		},
		"substring: partial mid-word match": {
			word1:    "here_is_a_string",
			word2:    "hey_is_ok",
			norm:     model.DistanceNormSubstring,
			expected: float64(0.35),
		},
		"substring: no match": {
			word1:    "here_is_a_string",
			word2:    "zzzzzzz",
			norm:     model.DistanceNormSubstring,
			expected: float64(0),
		},
		"levenshtein: full match": {
			word1:    "here_is_ok_nok",
			word2:    "here_is_ko",
			norm:     model.DistanceNormLevenshtein,
			expected: float64(0.25),
		},
		"levenshtein: reversed full match": {
			word1:    "here_is_ko",
			word2:    "here_is_ok_nok",
			norm:     model.DistanceNormLevenshtein,
			expected: float64(0.25),
		},
		"levenshtein: partial mid-word match": {
			word1:    "here_is_ok",
			word2:    "here_ar_ok",
			norm:     model.DistanceNormLevenshtein,
			expected: float64(0.25),
		},
		"levenshtein: no match": {
			word1:    "here_is_ok",
			word2:    "zzzzzzzzzz",
			norm:     model.DistanceNormLevenshtein,
			expected: float64(0.05),
		},
	}

	for name, c := range cases {
		require.Equal(t, wordWeight(c.word1, c.word2, c.norm), c.expected, name)
	}
}

func _buildFakeAffinityMap(p1 []string, p2 []string) map[string]map[string]float64 {
	res := map[string]map[string]float64{}
	for i, v := range p1 {
		res[v] = map[string]float64{}
		for j, w := range p2 {
			seed, err := strconv.ParseInt(fmt.Sprintf("%d%d", i, j), 10, 64)
			if err != nil {
				panic("failed to seed random float")
			}
			r := rand.New(rand.NewSource(seed))
			res[v][w] = r.Float64()
		}
	}
	return res
}
