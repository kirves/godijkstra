package dijkstra

import (
	"github.com/kirves/godijkstra/common/structs"
)

// DijkstraQueue is a collecition of DijkstraCandidate elements.
// It implements the heap.Interface interface to be used as a heap.
type DijkstraQueue []*dijkstrastructs.DijkstraCandidate

func (pq DijkstraQueue) Len() int {
	return len(pq)
}

func (pq DijkstraQueue) Less(i, j int) bool {
	return pq[i].Weight < pq[j].Weight
}

func (pq DijkstraQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *DijkstraQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*dijkstrastructs.DijkstraCandidate))
}

func (pq *DijkstraQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}
