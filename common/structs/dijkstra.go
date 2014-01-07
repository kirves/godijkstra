// Package DijkstraStructs contains support structures for the Go-Dijkstra component.
package dijkstrastructs

// DijkstraCandidate represent a node analyzed during the run of the Dijkstra Algorithm.
type DijkstraCandidate struct {
	Node   string             // Name of the analyzed node
	Parent *DijkstraCandidate // Parent in the candidate path (used for backtracking)
	Weight int                // Weight of the path so far
}

// CandidateSolution is a possibile complete path from source to destination in the provided graph.
type CandidateSolution struct {
	Length        int                // Weight of the solution
	ForwCandidate *DijkstraCandidate // Last DijkstraCandidate found in forward search
	BackCandidate *DijkstraCandidate // Last DijkstraCandidate found in backward search
}

// Connection is an outgoing edge from a given node to node Destination, having weight Weight
type Connection struct {
	Destination string // Destination node
	Weight      int    // Edge weight
}

// UnusableEdgeMap is a list of "banned" edges used by the deviation algorithm
type UnusableEdgeMap map[string]map[string]interface{}

// EmptyUnnusableEdgeMap creates an empty UnusableEdgeMap
func EmptyUnusableEdgeMap() UnusableEdgeMap {
	return make(map[string]map[string]interface{})
}
