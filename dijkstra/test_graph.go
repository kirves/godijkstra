/*
Copyright 2013 Alessandro Frossi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dijkstra

import (
	"github.com/kirves/godijkstra/common/structs"
)

type testGraph struct {
	nodes        map[string]interface{}
	edges        map[string]map[string]interface{}
	reverseEdges map[string]map[string]interface{}
}

func newTestGraph() *testGraph {
	return &testGraph{make(map[string]interface{}), make(map[string]map[string]interface{}), make(map[string]map[string]interface{})}
}

func (t *testGraph) SuccessorsForNode(node string) []dijkstrastructs.Connection {
	ret := make([]dijkstrastructs.Connection, len(t.edges[node]))
	i := 0
	for k, _ := range t.edges[node] {
		ret[i] = dijkstrastructs.Connection{k, t.EdgeWeight(node, k)}
		i++
	}
	return ret
}

func (t *testGraph) PredecessorsFromNode(node string) []dijkstrastructs.Connection {
	ret := make([]dijkstrastructs.Connection, len(t.reverseEdges[node]))
	i := 0
	for k, _ := range t.reverseEdges[node] {
		ret[i] = dijkstrastructs.Connection{k, t.EdgeWeight(k, node)}
		i++
	}
	return ret
}

func (t *testGraph) EdgeWeight(n1, n2 string) int {
	return 1
}
