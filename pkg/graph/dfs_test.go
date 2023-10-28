package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Dfs(t *testing.T) {
	cases := map[string]struct {
		vertexes []string
		graph    map[string][]string
		expected [][]string
	}{
		"nominal": {
			vertexes: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"},
			graph: map[string][]string{
				"a": {"c", "e"},
				"b": {},
				"c": {"a", "e"},
				"d": {"i", "k", "h"},
				"e": {"a", "c"},
				"f": {"l"},
				"g": {},
				"h": {"d", "i", "k"},
				"i": {"d", "k", "h"},
				"j": {},
				"k": {"d", "i", "h"},
				"l": {"f"},
			},
			expected: [][]string{
				{"a", "c", "e"},
				{"b"},
				{"d", "h", "i", "k"},
				{"f", "l"},
				{"g"},
				{"j"},
			},
		},
	}

	for _, c := range cases {
		res := Dfs(c.vertexes, c.graph)
		for _, e := range c.expected {
			require.Contains(t, res, e)
		}
	}
}
