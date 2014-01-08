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

import (
	"fmt"
	"github.com/kirves/godijkstra/common/structs"
)

var (
	cs    dijkstrastructs.CandidateSolution
	cs2   dijkstrastructs.CandidateSolution
	cs3   dijkstrastructs.CandidateSolution
	path  DijkstraPath
	path2 DijkstraPath
	path3 DijkstraPath
)

func init() {
	// Forward: START -> A -> B
	// Backward: B -> C -> D -> END
	forwNodes := []string{"START", "A", "B"}
	backNodes := []string{"B", "C", "D", "END"}
	var parent *dijkstrastructs.DijkstraCandidate = nil
	w := 0
	for _, v := range forwNodes {
		dc := dijkstrastructs.DijkstraCandidate{Node: v, Parent: parent, Weight: w}
		parent = &dc
		w++
	}
	cs.ForwCandidate = parent
	parent = nil
	w = 0
	for i := len(backNodes) - 1; i >= 0; i-- {
		dc := dijkstrastructs.DijkstraCandidate{Node: backNodes[i], Parent: parent, Weight: w}
		parent = &dc
		w++
	}
	cs.BackCandidate = parent

	// Forward: START -> A -> B
	// Backward: B -> E -> F -> G -> END
	forwNodes = []string{"START", "A", "B"}
	backNodes = []string{"B", "E", "F", "G", "END"}
	parent = nil
	w = 0
	for _, v := range forwNodes {
		dc := dijkstrastructs.DijkstraCandidate{Node: v, Parent: parent, Weight: w}
		parent = &dc
		w++
	}
	cs2.ForwCandidate = parent
	parent = nil
	w = 0
	for i := len(backNodes) - 1; i >= 0; i-- {
		dc := dijkstrastructs.DijkstraCandidate{Node: backNodes[i], Parent: parent, Weight: w}
		parent = &dc
		w++
	}
	cs2.BackCandidate = parent

	// Forward: A -> B -> J -> K
	// Backward: K -> L -> END
	forwNodes = []string{"B", "J", "K"}
	backNodes = []string{"K", "L", "END"}
	parent = nil
	w = 0
	for _, v := range forwNodes {
		dc := dijkstrastructs.DijkstraCandidate{Node: v, Parent: parent, Weight: w}
		parent = &dc
		w++
	}
	cs3.ForwCandidate = parent
	parent = nil
	w = 0
	for i := len(backNodes) - 1; i >= 0; i-- {
		dc := dijkstrastructs.DijkstraCandidate{Node: backNodes[i], Parent: parent, Weight: w}
		parent = &dc
		w++
	}
	cs3.BackCandidate = parent
}

func ExampleDijkstraPath_path() {
	path = ConvertToDijkstraPath(cs, "START", "END")
	for _, v := range path.Path {
		fmt.Printf("%s-", v.Node)
	}
	fmt.Printf("\n")

	path2 = ConvertToDijkstraPath(cs2, "START", "END")
	for _, v := range path2.Path {
		fmt.Printf("%s-", v.Node)
	}
	fmt.Printf("\n")

	path3 = ConvertToDijkstraPath(cs3, "A", "END")
	for _, v := range path3.Path {
		fmt.Printf("%s-", v.Node)
	}
	fmt.Printf("\n")

	// Output:
	// START-A-B-C-D-END-
	// START-A-B-E-F-G-END-
	// B-J-K-L-END-
}

func ExampleDijkstraPath_weight() {
	fmt.Println(path.Weight)
	fmt.Println(path2.Weight)
	fmt.Println(path3.Weight)

	// Output:
	// 5
	// 6
	// 4
}

func ExampleDijkstraPath_rootPaths() {
	subPaths := path.rootPaths()
	for _, v := range subPaths {
		for _, n := range v.Path {
			fmt.Printf("%s-", n.Node)
		}
		fmt.Printf(" :: w = %d\n", v.Weight)
	}

	// Output:
	// START- :: w = 0
	// START-A- :: w = 1
	// START-A-B- :: w = 2
	// START-A-B-C- :: w = 3
	// START-A-B-C-D- :: w = 4
}

func ExampleDijkstraPath_rootIncluded() {
	subPaths1 := path.rootPaths()
	for _, v := range subPaths1 {
		for _, n := range v.Path {
			fmt.Printf("%s-", n.Node)
		}
		fmt.Printf(" :: %t\n", path2.includesPath(v))
	}

	// Output:
	// START- :: true
	// START-A- :: true
	// START-A-B- :: true
	// START-A-B-C- :: false
	// START-A-B-C-D- :: false
}

func ExampleDijkstraPath_edgesToRemove() {
	subPaths2 := path2.rootPaths()
	for _, v := range subPaths2 {
		edges := make([][]string, 0)
		edges = append(edges, path.outgoingEdgeForSubPath(v))
		edges = append(edges, path2.outgoingEdgeForSubPath(v))
		for _, n := range v.Path {
			fmt.Printf("%s-", n.Node)
		}
		fmt.Printf(" :: %s\n", edges)
	}

	// Output:
	// START- :: [[START A] [START A]]
	// START-A- :: [[A B] [A B]]
	// START-A-B- :: [[B C] [B E]]
	// START-A-B-E- :: [[] [E F]]
	// START-A-B-E-F- :: [[] [F G]]
	// START-A-B-E-F-G- :: [[] [G END]]
}

func ExampleDijkstraPath_merge() {
	subPaths2 := path2.rootPaths()
	sp := subPaths2[2]
	for _, v := range sp.Path {
		fmt.Printf("%s-", v.Node)
	}
	fmt.Printf("\n")
	for _, v := range path2.Path {
		fmt.Printf("%s-", v.Node)
	}
	fmt.Printf("\n")
	ret := sp.mergeWith(path3)
	for _, v := range ret.Path {
		fmt.Printf("%s-", v.Node)
	}
	fmt.Printf("\n")
	fmt.Println("W =", ret.Weight)
	for _, v := range path2.Path {
		fmt.Printf("%s-", v.Node)
	}

	// Output:
	// START-A-B-
	// START-A-B-E-F-G-END-
	// START-A-B-J-K-L-END-
	// W = 6
	// START-A-B-E-F-G-END-
}

func ExampleDijkstraIterator_value() {
	it := path.rootPathIterator()
	for it.next() {
		fmt.Printf("%v\n", it.path().Path)
	}

	// Output:
	// [{START 0}]
	// [{START 0} {A 1}]
	// [{START 0} {A 1} {B 2}]
	// [{START 0} {A 1} {B 2} {C 3}]
	// [{START 0} {A 1} {B 2} {C 3} {D 4}]
}
