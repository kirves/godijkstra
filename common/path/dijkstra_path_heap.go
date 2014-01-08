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
