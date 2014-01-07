// Package DijkstraPath provides the path structures for the Go-Dijkstra component.
package dijkstrapath

import (
	"github.com/kirves/godijkstra/common/structs"
)

// Atmoic element of the path reperesented by the DijkstraPath struct
type DijkstraPathElement struct {
	Node   string // Node name
	Weight int    // Weight of the node (as computed by the Dijkstra algorithm)
}

// The DijkstraPath structure saves all the information about the found path between source and destination
// in the provided graph.
// It contains the succession of visited nodes, as well as the overall weight of the solution as returned by the Dijkstra algorithm
// and the names of starting and ending nodes for clarity purposes.
type DijkstraPath struct {
	Path      []DijkstraPathElement // Successions of DijkstraPathElements
	Weight    int                   // Weight of the solution
	StartNode string                // Name of the starting node
	EndNode   string                // Name of the target node
}

func newElementFromDijkstraCandidate(dc *dijkstrastructs.DijkstraCandidate) DijkstraPathElement {
	return DijkstraPathElement{dc.Node, dc.Weight}
}

// ConvertToDijkstraPath creates a DijkstraPath from the CandidateSolution instance return by a run of the DijkstraAlgorithm.
// It requires the starting and ending nodes for completeness' sake
func ConvertToDijkstraPath(cs dijkstrastructs.CandidateSolution, start, end string) DijkstraPath {
	ret := DijkstraPath{}

	tmp := appendForwardStepToDijkstraPath(cs.ForwCandidate, make([]DijkstraPathElement, 0))
	var dc *dijkstrastructs.DijkstraCandidate
	var parent DijkstraPathElement = newElementFromDijkstraCandidate(cs.ForwCandidate)
	realParent := cs.BackCandidate
	for dc = cs.BackCandidate.Parent; dc != nil; dc = dc.Parent {
		item := newElementFromDijkstraCandidate(dc)
		item.Weight = parent.Weight + (realParent.Weight - dc.Weight)
		tmp = append(tmp, item)
		parent = item
		realParent = dc
	}
	ret.Path = tmp
	ret.Weight = ret.computeWeight()
	ret.StartNode = start
	ret.EndNode = end
	return ret
}

func appendForwardStepToDijkstraPath(dc *dijkstrastructs.DijkstraCandidate, path []DijkstraPathElement) []DijkstraPathElement {
	if dc == nil {
		return path
	}
	tmp := appendForwardStepToDijkstraPath(dc.Parent, path)
	return append(tmp, newElementFromDijkstraCandidate(dc))
}

func (dp DijkstraPath) computeWeight() int {
	return dp.Path[len(dp.Path)-1].Weight
}

func (dp DijkstraPath) rootPaths() []DijkstraPath {
	ret := make([]DijkstraPath, len(dp.Path)-1)
	// var tmpParent *DijkstraCandidate = nil
	// var item *DijkstraCandidate
	for i := 0; i < len(dp.Path)-1; i++ {
		tmp := DijkstraPath{}
		tmp.Path = dp.Path[:i+1]
		tmp.Weight = tmp.computeWeight()
		tmp.StartNode = dp.StartNode
		tmp.EndNode = dp.EndNode
		ret[i] = tmp
	}
	return ret
}

func (dp DijkstraPath) lastNode() DijkstraPathElement {
	return dp.Path[len(dp.Path)-1]
}

func (dp DijkstraPath) includesPath(p DijkstraPath) bool {
	if len(dp.Path) < len(p.Path) {
		return false
	}
	for i, v := range p.Path {
		if dp.Path[i].Node != v.Node {
			return false
		}
	}
	return true
}

func (dp DijkstraPath) outgoingEdgeForSubPath(p DijkstraPath) []string {
	if !dp.includesPath(p) {
		return nil
	}
	edgeInd := len(p.Path)
	return []string{dp.Path[edgeInd-1].Node, dp.Path[edgeInd].Node}
}

func (dp DijkstraPath) mergeWith(p DijkstraPath) DijkstraPath {
	ret := DijkstraPath{}
	ret.Path = make([]DijkstraPathElement, len(dp.Path))
	for i, e := range dp.Path {
		ret.Path[i] = e
	}
	mergeParent := ret.Path[len(ret.Path)-1]
	realParent := p.Path[0]
	j := 0
	if dp.Path[len(dp.Path)-1].Node == p.Path[0].Node {
		j = 1
	}
	for ; j < len(p.Path); j++ {
		item := p.Path[j]
		item.Weight = mergeParent.Weight + (item.Weight - realParent.Weight)
		ret.Path = append(ret.Path, item)
		mergeParent = item
		realParent = p.Path[j]
	}
	ret.Weight = ret.computeWeight()
	ret.StartNode = dp.StartNode
	ret.EndNode = dp.EndNode
	return ret
}

func (dp DijkstraPath) rootPathIterator() dijkstraPathIterator {
	return dijkstraPathIterator{&dp, 0}
}
