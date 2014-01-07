package dijkstra

import (
	"github.com/kirves/godijkstra/common/structs"
)

// GraphObject interface defines the minimum requirements for an object to be considered a graph by Go-Dijkstra.
// It must implement three functionalities:
// getting the successors for a given node,
// getting the predecessors for a given node (this functionality is not required for the standard Dijkstra algorithm and can be a stub)
// and returning the non-negative edge weight associated to two nodes.
type GraphObject interface {
	SuccessorsForNode(node string) []dijkstrastructs.Connection    // get successors for node
	PredecessorsFromNode(node string) []dijkstrastructs.Connection // get predecessors for node
	EdgeWeight(n1, n2 string) int                                  // get edge weight
}
