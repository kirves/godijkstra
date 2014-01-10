package yen

import (
	"github.com/kirves/godijkstra/common/path"
	"github.com/kirves/godijkstra/common/structs"
	"github.com/kirves/godijkstra/dijkstra"
	"testing"
	"time"
)

var (
	graph *testGraph
)

func init() {
	nodes := map[string]interface{}{
		"A": struct{}{},
		"B": struct{}{},
		"C": struct{}{},
		"D": struct{}{},
		"E": struct{}{},
		"F": struct{}{},
		"G": struct{}{},
		"S": struct{}{},
		"T": struct{}{},
		"U": struct{}{},
	}
	edges := map[string]map[string]interface{}{
		"A": map[string]interface{}{
			"B": struct{}{},
			"C": struct{}{},
		},
		"B": map[string]interface{}{
			"D": struct{}{},
		},
		"C": map[string]interface{}{
			"E": struct{}{},
			"G": struct{}{},
		},
		"D": map[string]interface{}{
			"C": struct{}{},
		},
		"E": map[string]interface{}{
			"F": struct{}{},
			"G": struct{}{},
		},
		"F": map[string]interface{}{
			"G": struct{}{},
		},
		"G": map[string]interface{}{
			"T": struct{}{},
		},
		"S": map[string]interface{}{
			"A": struct{}{},
			// "B": struct{}{},
		},
	}
	graph = newTestGraph()
	graph.nodes = nodes
	// revEdges := make(map[string]map[string]interface{})
	graph.edges = edges
	for k1, v := range edges {
		for k2, _ := range v {
			if _, ok := graph.reverseEdges[k2]; !ok {
				graph.reverseEdges[k2] = make(map[string]interface{})
			}
			graph.reverseEdges[k2][k1] = struct{}{}
		}
	}
}

func TestMultiplePaths(t *testing.T) {
	paths := Yen(graph, "S", "T", 4, dijkstra.Dijkstra)
	if len(paths) == 0 {
		t.Fatal("Didn't find any paths.")
	}
	expPath := [][]string{
		[]string{"S", "A", "C", "G", "T"},
		[]string{"S", "A", "C", "E", "G", "T"},
		[]string{"S", "A", "B", "D", "C", "G", "T"},
		[]string{"S", "A", "C", "E", "F", "G", "T"},
	}
	for k, p := range paths {
		for i, v := range p.Path {
			if v.Node != expPath[k][i] {
				t.Fatalf("Wrong path (%d).\n", k)
			}
		}
	}
}

func yenWrapper(
	graph dijkstrastructs.GraphObject,
	startNode, endNode string,
	k int,
	searchFunc func(dijkstrastructs.GraphObject, string, string, dijkstrastructs.UnusableEdgeMap) (dijkstrapath.DijkstraPath, bool),
	ch chan []dijkstrapath.DijkstraPath) {
	ch <- Yen(graph, startNode, endNode, k, searchFunc)
}

func TestTermination(t *testing.T) {
	k := 0
	paths := make([]dijkstrapath.DijkstraPath, 0)
	ch := make(chan []dijkstrapath.DijkstraPath)
	for len(paths) >= k {
		k++
		go yenWrapper(graph, "S", "T", k, dijkstra.Dijkstra, ch)
		select {
		case paths = <-ch:
			continue
		case <-time.After(30 * time.Second):
			t.Fatalf("Path search algorithm did not terminate after 30 secsonds.\nSomething went wrong.")
		}
	}
}
