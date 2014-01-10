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

// Package Dijkstra provides an implementation of Dijkstra Algorithm to find the shortest path in directed graph.
//
// The Dijkstra Algorithm traverses a graph object implementing the dijkstrastruct.GraphObject interface to find the shortest path;
// the only limitation is that all the edges' weights must be non-negative.
// The returned path, an instance of DijkstraPath struct, is a loopless path going from the starting node to the destination;
// it can be computed using either the "vanilla" Dijkstra algorithm or a bidirectional search algorithm.
package dijkstra

import (
	"container/heap"
	"github.com/kirves/godijkstra/common/path"
	"github.com/kirves/godijkstra/common/structs"
)

const (
	VANILLA = iota // Use "vanilla" Dijkstra algorithm
	BIDIR          // Use bi-directional search algorithm
)

func newDijkstraCandidate(node string, parent *dijkstrastructs.DijkstraCandidate, w int) *dijkstrastructs.DijkstraCandidate {
	return &dijkstrastructs.DijkstraCandidate{node, parent, w}
}

// func (cs CandidateSolution) IsEqualTo(sol CandidateSolution) bool {
// 	if cs.length != sol.length || cs.forwCandidate.node != sol.forwCandidate.node || cs.backCandidate.node != sol.backCandidate.node {
// 		return false
// 	}
// 	return true
// }

// Dijkstra returns the shortest path within the provided graph object that goes from startNode to endNode nodes.
// searchType parameter defines the type of algorithm to use.
func SearchPath(graph dijkstrastructs.GraphObject, startNode, endNode string, searchType int) (dijkstrapath.DijkstraPath, bool) {
	switch searchType {
	case VANILLA:
		return Dijkstra(graph, startNode, endNode, dijkstrastructs.EmptyUnusableEdgeMap())
	case BIDIR:
		return BiDirDijkstra(graph, startNode, endNode, dijkstrastructs.EmptyUnusableEdgeMap())
	default:
		return dijkstrapath.DijkstraPath{}, false
	}
}

func Dijkstra(graph dijkstrastructs.GraphObject, startNode, endNode string, bannedEdges dijkstrastructs.UnusableEdgeMap) (dijkstrapath.DijkstraPath, bool) {
	// SETUP ================================
	firstParent := newDijkstraCandidate(startNode, nil, 0)
	startSet := []*dijkstrastructs.DijkstraCandidate{firstParent}
	// ======================================
	cs, valid := computeVanillaDijkstra(graph, startSet, endNode, bannedEdges)
	if !valid {
		return dijkstrapath.DijkstraPath{}, false
	}
	return dijkstrapath.ConvertToDijkstraPath(cs, startNode, endNode), true
}

func BiDirDijkstra(graph dijkstrastructs.GraphObject, startNode, endNode string, bannedEdges dijkstrastructs.UnusableEdgeMap) (dijkstrapath.DijkstraPath, bool) {
	// SETUP ================================
	firstParent := newDijkstraCandidate(startNode, nil, 0)
	lastParent := newDijkstraCandidate(endNode, nil, 0)
	startSet := []*dijkstrastructs.DijkstraCandidate{firstParent}
	endSet := []*dijkstrastructs.DijkstraCandidate{lastParent}
	// ======================================
	cs, valid := computeBiDirDijkstra(graph, startSet, endSet, dijkstrastructs.EmptyUnusableEdgeMap())
	if !valid {
		return dijkstrapath.DijkstraPath{}, false
	}
	return dijkstrapath.ConvertToDijkstraPath(cs, startNode, endNode), true
}

func computeVanillaDijkstra(
	graph dijkstrastructs.GraphObject,
	startSet []*dijkstrastructs.DijkstraCandidate,
	endNode string,
	bannedEdges dijkstrastructs.UnusableEdgeMap) (dijkstrastructs.CandidateSolution, bool) {

	candidateSolution := dijkstrastructs.CandidateSolution{0, nil, nil}
	var succs []dijkstrastructs.Connection
	visitedNodesF := make(map[string]*dijkstrastructs.DijkstraCandidate)

	openListF := &DijkstraQueue{}
	heap.Init(openListF)

	// create initial path set
	for _, c := range startSet {
		heap.Push(openListF, c)
	}

	for openListF.Len() > 0 {

		// get candidates
		forwCandidate := heap.Pop(openListF).(*dijkstrastructs.DijkstraCandidate)

		// check if we reached termination
		if forwCandidate.Node == endNode {
			return dijkstrastructs.CandidateSolution{forwCandidate.Weight, forwCandidate, &dijkstrastructs.DijkstraCandidate{endNode, nil, 0}}, true
		}

		if _, ok := visitedNodesF[forwCandidate.Node]; ok {
			continue
		} else {
			visitedNodesF[forwCandidate.Node] = forwCandidate
		}

		succs = successorsForPath(graph, forwCandidate, bannedEdges)

		// for each successors
		for _, s := range succs {
			if _, ok := visitedNodesF[s.Destination]; ok {
				continue
			}
			newPath := newDijkstraCandidate(s.Destination, forwCandidate, forwCandidate.Weight+s.Weight)
			// duplicate and add step
			heap.Push(openListF, newPath)
		}
	}
	return candidateSolution, false
}

func computeBiDirDijkstra(
	graph dijkstrastructs.GraphObject,
	startSet []*dijkstrastructs.DijkstraCandidate,
	endSet []*dijkstrastructs.DijkstraCandidate,
	bannedEdges dijkstrastructs.UnusableEdgeMap) (dijkstrastructs.CandidateSolution, bool) {

	candidateSolution := dijkstrastructs.CandidateSolution{0, nil, nil}
	skipForward := false
	var succs []dijkstrastructs.Connection
	visitedNodesF := make(map[string]*dijkstrastructs.DijkstraCandidate)
	visitedNodesB := make(map[string]*dijkstrastructs.DijkstraCandidate)

	openListF := &DijkstraQueue{}
	openListB := &DijkstraQueue{}
	heap.Init(openListF)
	heap.Init(openListB)

	// create initial path set
	for _, c := range startSet {
		heap.Push(openListF, c)
	}

	for _, c := range endSet {
		heap.Push(openListB, c)
	}

	for openListF.Len() > 0 && openListB.Len() > 0 {

		// get candidates
		forwCandidate := heap.Pop(openListF).(*dijkstrastructs.DijkstraCandidate)
		backCandidate := heap.Pop(openListB).(*dijkstrastructs.DijkstraCandidate)

		// check if we reached termination
		if candidateSolution.Length != 0 && forwCandidate.Weight+backCandidate.Weight >= candidateSolution.Length {
			break
		}

		// ***************************************************
		// forward search
		if _, ok := visitedNodesF[forwCandidate.Node]; ok {
			skipForward = true
		} else {
			visitedNodesF[forwCandidate.Node] = forwCandidate
			skipForward = false
		}

		if !skipForward {
			if v, ok := visitedNodesB[forwCandidate.Node]; ok {
				// found an explored backward path
				newWeight := forwCandidate.Weight + v.Weight
				if candidateSolution.Length == 0 || candidateSolution.Length > newWeight {
					// found new solution candidate
					candidateSolution.Length = newWeight
					candidateSolution.ForwCandidate = forwCandidate
					candidateSolution.BackCandidate = visitedNodesB[forwCandidate.Node]
				}
			}

			succs = successorsForPath(graph, forwCandidate, bannedEdges)

			// for each successors
			for _, s := range succs {
				if _, ok := visitedNodesF[s.Destination]; ok {
					continue
				}
				newPath := newDijkstraCandidate(s.Destination, forwCandidate, forwCandidate.Weight+s.Weight)
				// duplicate and add step
				heap.Push(openListF, newPath)
			}
		}
		// ****************************************************

		// ***************************************************
		// backward search
		if _, ok := visitedNodesB[backCandidate.Node]; ok {
			continue
		} else {
			visitedNodesB[backCandidate.Node] = backCandidate
		}

		if v, ok := visitedNodesF[backCandidate.Node]; ok {
			// found an explored backward path
			newWeight := backCandidate.Weight + v.Weight
			if candidateSolution.Length == 0 || candidateSolution.Length > newWeight {
				// found new solution candidate
				candidateSolution.Length = newWeight
				candidateSolution.ForwCandidate = visitedNodesF[backCandidate.Node]
				candidateSolution.BackCandidate = backCandidate
			}
		}

		succs = predecessorsForPath(graph, backCandidate, bannedEdges)

		// for each successors
		for _, s := range succs {
			if _, ok := visitedNodesB[s.Destination]; ok {
				continue
			}
			newPath := newDijkstraCandidate(s.Destination, backCandidate, backCandidate.Weight+s.Weight)
			heap.Push(openListB, newPath)
		}
		// ****************************************************
	}

	if candidateSolution.Length == 0 {
		return dijkstrastructs.CandidateSolution{}, false
	}
	return candidateSolution, true
}

func successorsForPath(graph dijkstrastructs.GraphObject, path *dijkstrastructs.DijkstraCandidate, bannedEdges dijkstrastructs.UnusableEdgeMap) []dijkstrastructs.Connection {
	tmp := graph.SuccessorsForNode(path.Node)
	ret := make([]dijkstrastructs.Connection, 0)
	for _, s := range tmp {
		if bannedEdges[path.Node][s.Destination] == nil {
			ret = append(ret, s)
		}
	}
	return ret
}

func predecessorsForPath(graph dijkstrastructs.GraphObject, path *dijkstrastructs.DijkstraCandidate, bannedEdges dijkstrastructs.UnusableEdgeMap) []dijkstrastructs.Connection {
	tmp := graph.PredecessorsFromNode(path.Node)
	ret := make([]dijkstrastructs.Connection, 0)
	for _, s := range tmp {
		if bannedEdges[s.Destination][path.Node] == nil {
			ret = append(ret, s)
		}
	}
	return ret
}
