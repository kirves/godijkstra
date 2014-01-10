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

// Package Yen implements Yen's algorithm for k-shortest paths search in a graph.
//
// Yen's algorithm is made of two components: the shortest path search and the deviation from already found solutions.
// This package makes use of the DijkstraPath structure to implement the deviation algorithm, while leaves to the developer the
// choice for the search algorithm. This implementation of Yen's algorithm has been successfully tested using the
// dijstra bidirectional algorithm provided within the same package.
package yen

import (
	"container/heap"
	"github.com/kirves/godijkstra/common/path"
	"github.com/kirves/godijkstra/common/structs"
)

func Yen(
	graph dijkstrastructs.GraphObject,
	startNode, endNode string,
	k int,
	searchFunc func(dijkstrastructs.GraphObject, string, string, dijkstrastructs.UnusableEdgeMap) (dijkstrapath.DijkstraPath, bool)) []dijkstrapath.DijkstraPath {

	if k <= 0 {
		return make([]dijkstrapath.DijkstraPath, 0)
	}

	// FIRST SOLUTION ========================
	dp, valid := searchFunc(graph, startNode, endNode, dijkstrastructs.EmptyUnusableEdgeMap())
	if !valid {
		return make([]dijkstrapath.DijkstraPath, 0)
	}

	// add to heap
	finalList := make([]dijkstrapath.DijkstraPath, 0)
	// foundPaths := make(map[string]interface{})

	candidateHeap := &dijkstrapath.DijkstraPathQueue{}
	heap.Init(candidateHeap)
	heap.Push(candidateHeap, dp)

	var cdp dijkstrapath.DijkstraPath
	for candidateHeap.Len() > 0 {
		cdp = heap.Pop(candidateHeap).(dijkstrapath.DijkstraPath)

		// add to finals and check if we're done
		finalList = append(finalList, cdp)

		if len(finalList) >= k {
			break
		}

		for _, rp := range cdp.RootPaths() {
			bannedEdges := dijkstrastructs.EmptyUnusableEdgeMap()
			for _, path := range finalList {
				be := path.OutgoingEdgeForSubPath(rp)
				if be != nil {
					if _, ok := bannedEdges[be[0]]; !ok {
						bannedEdges[be[0]] = make(map[string]interface{})
					}
					bannedEdges[be[0]][be[1]] = struct{}{}
				}
			}

			// build start and end sets
			// 3 cases
			ln := rp.LastNode()
			dp, valid = searchFunc(graph, ln.Node, endNode, bannedEdges)
			if !valid {
				continue
			}
			dp = rp.MergeWith(dp)

			heap.Push(candidateHeap, dp)
		}

	}
	return finalList
}
