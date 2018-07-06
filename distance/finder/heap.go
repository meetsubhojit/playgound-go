package finder

type MinHeapCoordinate []Coordinate

func (h MinHeapCoordinate) Len() int           { return len(h) }
func (h MinHeapCoordinate) Less(i, j int) bool { return h[i].Distance < h[j].Distance }
func (h MinHeapCoordinate) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeapCoordinate) Push(x interface{}) {
	*h = append(*h, x.(Coordinate))
}
func (h *MinHeapCoordinate) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type MaxHeapCoordinate struct {
	MinHeapCoordinate
}

func (h MaxHeapCoordinate) Less(i, j int) bool {
	return h.MinHeapCoordinate[i].Distance > h.MinHeapCoordinate[j].Distance
}
