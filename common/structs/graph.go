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

package dijkstrastructs

// GraphObject interface defines the minimum requirements for an object to be considered a graph by Go-Dijkstra.
// It must implement three functionalities:
// getting the successors for a given node,
// getting the predecessors for a given node (this functionality is not required for the standard Dijkstra algorithm and can be a stub)
// and returning the non-negative edge weight associated to two nodes.
type GraphObject interface {
	SuccessorsForNode(node string) []Connection    // get successors for node
	PredecessorsFromNode(node string) []Connection // get predecessors for node
	EdgeWeight(n1, n2 string) int                  // get edge weight
}
