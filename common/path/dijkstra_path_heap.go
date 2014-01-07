package dijkstrapath

// Support structure to allow the usage of heap containers with DijkstraPaths.
// DijkstraPathQueue implements the heap.Interface interface
type DijkstraPathQueue []DijkstraPath

func (pq DijkstraPathQueue) Len() int {
	return len(pq)
}

func (pq DijkstraPathQueue) Less(i, j int) bool {
	return pq[i].Weight < pq[j].Weight
}

func (pq DijkstraPathQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *DijkstraPathQueue) Push(x interface{}) {
	*pq = append(*pq, x.(DijkstraPath))
}

func (pq *DijkstraPathQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}
