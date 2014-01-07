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
