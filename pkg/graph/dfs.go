package graph

import (
	"sort"

	"github.com/epsxy/flower/pkg/utils"
)

func Dfs(vertexes []string, graph map[string][]string) [][]string {
	var partitions [][]string = [][]string{}
	var visited map[string]bool = map[string]bool{}
	for _, v := range vertexes {
		if !utils.Array2DContains(partitions, v) {
			c := _dfs_impl(v, graph, visited, []string{})
			sort.Strings(c)
			partitions = append(partitions, c)
		}
	}
	return partitions
}

func _dfs_impl(vertex string, graph map[string][]string, visited map[string]bool, current []string) []string {
	visited[vertex] = true
	current = utils.AppendWithoutDuplicates(current, vertex)
	for _, v := range graph[vertex] {
		if !visited[v] {
			current = _dfs_impl(v, graph, visited, current)
		}
	}
	return current
}
