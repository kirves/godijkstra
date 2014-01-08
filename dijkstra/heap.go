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
