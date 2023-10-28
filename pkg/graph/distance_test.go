package graph

import (
	"testing"

	"github.com/epsxy/flower/pkg/model"
	"github.com/stretchr/testify/require"
)

func Test_Split(t *testing.T) {
	// TODO
}

func Test_ReArrangePartitions(t *testing.T) {
	// TODO
}

func Test_Weight(t *testing.T) {
	// TODO
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
			expected: float64(10.53),
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
			norm:     model.DistanceNormSubstring,
			expected: map[string]map[string]float64{
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
		},
	}
	for name, c := range cases {
		require.Equal(t, buildAffinityMatrix(c.vertexes, c.norm), c.expected, name)
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
			expected: float64(1),
		},
		"substring: reversed full match": {
			word1:    "here_is_a_string_suffixed",
			word2:    "here_is_a_string",
			norm:     model.DistanceNormSubstring,
			expected: float64(0.64),
		},
		"substring: partial mid-word match": {
			word1:    "here_is_a_string",
			word2:    "hey_is_ok",
			norm:     model.DistanceNormSubstring,
			expected: float64(0.25),
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
