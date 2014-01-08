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
