package graph

import (
	"github.com/epsxy/flower/pkg/global"
	"github.com/epsxy/flower/pkg/model"
)

func Gen(vertexes []string, edges map[string]*model.EntityLink) map[string][]string {
	var datagraph map[string][]string = map[string][]string{}
	for _, vertex := range vertexes {
		datagraph[vertex] = []string{}
	}
	for _, link := range edges {
		if link.Left != nil {
			if !global.Contains(datagraph[link.Left.SourceName], link.Left.DestinationName) {
				datagraph[link.Left.SourceName] = append(datagraph[link.Left.SourceName], link.Left.DestinationName)
			}
			if !global.Contains(datagraph[link.Left.DestinationName], link.Left.SourceName) {
				datagraph[link.Left.DestinationName] = append(datagraph[link.Left.DestinationName], link.Left.SourceName)
			}
		} else if link.Right != nil {
			if !global.Contains(datagraph[link.Right.SourceName], link.Right.DestinationName) {
				datagraph[link.Right.SourceName] = append(datagraph[link.Right.SourceName], link.Right.DestinationName)
			}
			if !global.Contains(datagraph[link.Right.DestinationName], link.Right.SourceName) {
				datagraph[link.Right.DestinationName] = append(datagraph[link.Right.DestinationName], link.Right.SourceName)
			}
		}
	}
	return datagraph
}
