package graph

import (
	"fmt"
	"sort"
	"testing"

	"github.com/epsxy/flower/pkg/model"
	"github.com/stretchr/testify/require"
)

func Test_Gen(t *testing.T) {
	cases := map[string]struct {
		vertexes []string
		links    map[string]*model.EntityLink
		expected map[string][]string
	}{
		"nominal": {
			vertexes: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"},
			links: map[string]*model.EntityLink{
				"a_b": {
					Left: &model.Link{
						SourceName:      "a",
						DestinationName: "b",
						IsNullable:      false,
					},
					Right: nil,
				},
				"a_c": {
					Left: &model.Link{
						SourceName:      "a",
						DestinationName: "c",
						IsNullable:      true,
					},
					Right: nil,
				},
				"c_e": {
					Left: &model.Link{
						SourceName:      "c",
						DestinationName: "e",
						IsNullable:      true,
					},
					Right: &model.Link{
						SourceName:      "e",
						DestinationName: "c",
						IsNullable:      true,
					},
				},
				"e_g": {
					Left: nil,
					Right: &model.Link{
						SourceName:      "e",
						DestinationName: "g",
						IsNullable:      true,
					},
				},
				"e_h": {
					Left: nil,
					Right: &model.Link{
						SourceName:      "e",
						DestinationName: "h",
						IsNullable:      true,
					},
				},
				"i_j": {
					Left: &model.Link{
						SourceName:      "i",
						DestinationName: "j",
						IsNullable:      true,
					},
					Right: nil,
				},
				"i_k": {
					Left: &model.Link{
						SourceName:      "i",
						DestinationName: "k",
						IsNullable:      true,
					},
					Right: nil,
				},
				"i_l": {
					Right: &model.Link{
						SourceName:      "l",
						DestinationName: "i",
						IsNullable:      true,
					},
				},
			},
			expected: map[string][]string{
				"a": {"b", "c"},
				"b": {"a"},
				"c": {"a", "e"},
				"d": {},
				"e": {"c", "g", "h"},
				"f": {},
				"g": {"e"},
				"h": {"e"},
				"i": {"j", "k", "l"},
				"j": {"i"},
				"k": {"i"},
				"l": {"i"},
			},
		},
	}

	for _, c := range cases {
		res := Gen(c.vertexes, c.links)
		for i, e := range c.expected {
			// we don't care about the order
			sort.Strings(e)
			sort.Strings(res[i])
			require.Equal(t, e, res[i], fmt.Sprintf("graph built fail for key: %s", i))
		}
	}
}
