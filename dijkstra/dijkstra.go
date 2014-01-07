// Package Dijkstra provides an implementation of Dijkstra Algorithm to find the shortest path in directed graph.
//
// The Dijkstra Algorithm traverses a graph object implementing the dijkstra.GraphObject interface to find the shortest path;
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
func Dijkstra(graph GraphObject, startNode, endNode string, searchType int) (dijkstrapath.DijkstraPath, bool) {

	// SETUP ================================
	firstParent := NewDijkstraCandidate(startNode, nil, 0)
	lastParent := NewDijkstraCandidate(endNode, nil, 0)
	startSet := []*dijkstrastructs.DijkstraCandidate{firstParent}
	endSet := []*dijkstrastructs.DijkstraCandidate{lastParent}
	// ======================================

	var cs dijkstrastructs.CandidateSolution
	var valid bool
	switch searchType {
	case VANILLA:
		cs, valid = computeBiDirDijkstra(graph, startSet, endSet, dijkstrastructs.EmptyUnusableEdgeMap())
	case BIDIR:
		cs, valid = computeBiDirDijkstra(graph, startSet, endSet, dijkstrastructs.EmptyUnusableEdgeMap())
	default:
		valid = false
	}
	if !valid {
		return dijkstrapath.DijkstraPath{}, false
	}
	return dijkstrapath.ConvertToDijkstraPath(cs, startNode, endNode), true
}

// func (p *Planner) Yen(
// 	startNode, endNode common.NodePoint,
// 	mwd, mws float64,
// 	mlc int,
// 	penaltyFunction func() int,
// 	numSol int) []DijkstraPath {

// 	if numSol <= 0 {
// 		numSol = 1
// 	}

// 	// SETUP ================================
// 	start := time.Now()
// 	startConns, endConns := p.graph.GetConnectionsForStartingAndEndingNodes(startNode, endNode, mwd)
// 	firstParent := p.NewDijkstraCandidate(planningcommon.STARTING_NODE_NAME, planningcommon.STARTING_LINE_NAME, nil, 0, mws, penaltyFunction())
// 	lastParent := p.NewDijkstraCandidate(planningcommon.ENDING_NODE_NAME, planningcommon.ENDING_LINE_NAME, nil, 0, mws, penaltyFunction())
// 	startSet := make([]*DijkstraCandidate, len(startConns))
// 	endSet := make([]*DijkstraCandidate, len(endConns))
// 	// create initial path set
// 	i := 0
// 	for _, c := range startConns {
// 		cand := p.NewDijkstraCandidate(c.EndNode, c.EndLine, firstParent, c.Distance, mws, penaltyFunction())
// 		startSet[i] = cand
// 		i++
// 	}

// 	i = 0
// 	for _, c := range endConns {
// 		cand := p.NewDijkstraCandidate(c.EndNode, c.EndLine, lastParent, c.Distance, mws, penaltyFunction())
// 		endSet[i] = cand
// 		i++
// 	}
// 	p.logger.Warning("SETUP TIME: %.5fms", float64(time.Since(start).Nanoseconds())/1000000.0)
// 	// p.logger.Debug("START SET: %#v", startSet)
// 	// p.logger.Debug("END SET: %#v", endSet)
// 	// =======================================

// 	// FIRST SOLUTION ========================
// 	var cs CandidateSolution
// 	var valid bool
// 	start = time.Now()
// 	cs, valid = p.computeBiDirDijkstra(startSet, endSet, mwd, mws, mlc, planningcommon.EmptyUnusableEdgeMap(), penaltyFunction)
// 	if !valid {
// 		return make([]DijkstraPath, 0)
// 	}
// 	dp := convertToDijkstraPath(cs, startNode, endNode, mwd, mws)
// 	p.logger.Debug("DP SIGNATURE: %s", dp.signature)

// 	// add to heap
// 	finalList := make([]DijkstraPath, 0)
// 	foundPaths := make(map[string]interface{})

// 	candidateHeap := &DijkstraPathQueue{}
// 	heap.Init(candidateHeap)
// 	heap.Push(candidateHeap, dp)

// 	var cdp DijkstraPath
// 	for candidateHeap.Len() > 0 {
// 		cdp = heap.Pop(candidateHeap).(DijkstraPath)

// 		// SANITY CHECK
// 		if !cdp.sanityCheck() {
// 			p.logger.Critical("DijkstraPath is not valid :: %#v", cdp)
// 			continue
// 		}
// 		// ====================

// 		// add to finals and check if we're done
// 		cdpSign := p.lineSignatureForDijsktraPath(cdp)
// 		if _, ok := foundPaths[cdpSign]; ok {
// 			continue
// 		} else {
// 			foundPaths[cdpSign] = struct{}{}
// 		}
// 		p.logger.Debug("PUSHING FINAL: %#v", cdp)
// 		finalList = append(finalList, cdp)
// 		p.logger.Warning("GOT SOLUTION IN: %.5fms", float64(time.Since(start).Nanoseconds())/1000000.0)
// 		start = time.Now()

// 		if len(finalList) >= numSol {
// 			break
// 		}

// 		for _, rp := range cdp.rootPaths() {
// 			bannedEdges := planningcommon.EmptyUnusableEdgeMap()
// 			for _, path := range finalList {
// 				be := path.outgoingEdgeForPath(rp)
// 				if be != nil {
// 					var startLine string
// 					switch be[0] {
// 					case planningcommon.STARTING_LINE_NAME:
// 						startLine = planningcommon.STARTING_LINE_NAME
// 					case planningcommon.ENDING_LINE_NAME:
// 						startLine = planningcommon.ENDING_LINE_NAME
// 					default:
// 						startLine = p.graph.GetLineInfo(be[0]).No
// 					}

// 					var endLine string
// 					switch be[1] {
// 					case planningcommon.STARTING_LINE_NAME:
// 						endLine = planningcommon.STARTING_LINE_NAME
// 					case planningcommon.ENDING_LINE_NAME:
// 						endLine = planningcommon.ENDING_LINE_NAME
// 					default:
// 						endLine = p.graph.GetLineInfo(be[1]).No
// 					}

// 					if _, ok := bannedEdges[startLine]; !ok {
// 						bannedEdges[startLine] = make(map[string]interface{})
// 					}
// 					bannedEdges[startLine][endLine] = struct{}{}
// 				}
// 			}
// 			p.logger.Debug("Banned edges: %#v", bannedEdges)

// 			// build start and end sets
// 			// 3 cases
// 			ln := rp.lastNode()
// 			switch ln.node {
// 			case planningcommon.STARTING_NODE_NAME:
// 				startSet = make([]*DijkstraCandidate, 0)
// 				endSet = make([]*DijkstraCandidate, 0)
// 				firstParent := p.NewDijkstraCandidate(planningcommon.STARTING_NODE_NAME, planningcommon.STARTING_LINE_NAME, nil, 0, mws, penaltyFunction())
// 				lastParent := p.NewDijkstraCandidate(planningcommon.ENDING_NODE_NAME, planningcommon.ENDING_LINE_NAME, nil, 0, mws, penaltyFunction())
// 				for _, c := range startConns {
// 					if bannedEdges[planningcommon.STARTING_LINE_NAME][p.graph.GetNodeInfo(c.EndLine, c.EndNode).LineNo] == nil {
// 						cand := p.NewDijkstraCandidate(c.EndNode, c.EndLine, firstParent, c.Distance, mws, penaltyFunction())
// 						startSet = append(startSet, cand)
// 					}
// 				}
// 				for _, c := range endConns {
// 					if bannedEdges[p.graph.GetNodeInfo(c.EndLine, c.EndNode).LineNo][planningcommon.ENDING_LINE_NAME] == nil {
// 						cand := p.NewDijkstraCandidate(c.EndNode, c.EndLine, lastParent, c.Distance, mws, penaltyFunction())
// 						endSet = append(endSet, cand)
// 					}
// 				}
// 			default:
// 				startSet = make([]*DijkstraCandidate, 1)
// 				endSet = make([]*DijkstraCandidate, 0)
// 				firstParent := p.NewDijkstraCandidate(ln.node, ln.line, nil, 0, mws, penaltyFunction())
// 				lastParent := p.NewDijkstraCandidate(planningcommon.ENDING_NODE_NAME, planningcommon.ENDING_LINE_NAME, nil, 0, mws, penaltyFunction())
// 				startSet[0] = firstParent
// 				for _, c := range endConns {
// 					if bannedEdges[p.graph.GetNodeInfo(c.EndLine, c.EndNode).LineNo][planningcommon.ENDING_LINE_NAME] == nil {
// 						cand := p.NewDijkstraCandidate(c.EndNode, c.EndLine, lastParent, c.Distance, mws, penaltyFunction())
// 						endSet = append(endSet, cand)
// 					}
// 				}
// 			}
// 			cs, valid = p.computeBiDirDijkstra(startSet, endSet, mwd, mws, mlc, bannedEdges, penaltyFunction)
// 			if !valid {
// 				continue
// 			}
// 			dp = convertToDijkstraPath(cs, startNode, endNode, mwd, mws)
// 			dp = rp.mergeWith(dp)

// 			heap.Push(candidateHeap, dp)
// 			p.logger.Debug("Found candidate: %#v", dp.signature)
// 		}

// 	}

// 	// ret := make([]*PlanningPath, len(finalList))
// 	for _, path := range finalList {
// 		// p.logger.Debug("Saving candidate: %#v", path)
// 		p.printSolution(path)
// 		// 	ret[j] = path.convertToPlanningPath(NewPlanningPath(startNode, endNode, mwd, mws, mlc, make([]string, 0), make([]string, 0), penaltyFunction))
// 		// 	// ret[j] = createPlanningPath(cs, NewPlanningPath(startNode, endNode, mwd, mws, mlc, make([]string, 0), make([]string, 0), penaltyFunction))
// 	}
// 	return finalList
// }

func computeBiDirDijkstra(
	graph GraphObject,
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

func successorsForPath(graph GraphObject, path *dijkstrastructs.DijkstraCandidate, bannedEdges dijkstrastructs.UnusableEdgeMap) []dijkstrastructs.Connection {
	tmp := graph.SuccessorsForNode(path.Node)
	ret := make([]dijkstrastructs.Connection, 0)
	for _, s := range tmp {
		if bannedEdges[path.Node][s.Destination] == nil {
			ret = append(ret, s)
		}
	}
	return ret
}

func predecessorsForPath(graph GraphObject, path *dijkstrastructs.DijkstraCandidate, bannedEdges dijkstrastructs.UnusableEdgeMap) []dijkstrastructs.Connection {
	tmp := graph.PredecessorsFromNode(path.Node)
	ret := make([]dijkstrastructs.Connection, 0)
	for _, s := range tmp {
		if bannedEdges[s.Destination][path.Node] == nil {
			ret = append(ret, s)
		}
	}
	return ret
}
