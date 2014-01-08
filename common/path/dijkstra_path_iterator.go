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

type dijkstraPathIterator struct {
	dp  *DijkstraPath
	ind int
}

func (dpi *dijkstraPathIterator) next() bool {
	dpi.ind++
	if dpi.ind >= len(dpi.dp.Path) {
		return false
	}
	return true
}

func (dpi *dijkstraPathIterator) path() DijkstraPath {
	tmp := DijkstraPath{}
	tmp.Path = dpi.dp.Path[:dpi.ind]
	tmp.Weight = tmp.computeWeight()
	tmp.StartNode = dpi.dp.StartNode
	tmp.EndNode = dpi.dp.EndNode
	return tmp
}
