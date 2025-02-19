package geom

// A LineString represents a single, unbroken line, linearly interpreted
// between zero or more control points.
type LineString struct {
	geom1
}

// NewLineString returns a new LineString with layout l and no control points.
func NewLineString(l Layout) *LineString {
	return NewLineStringFlat(l, nil)
}

// NewLineStringFlat returns a new LineString with layout l and control points
// flatCoords.
func NewLineStringFlat(layout Layout, flatCoords []float64) *LineString {
	g := new(LineString)
	g.Layout = layout
	g.Stride = layout.Stride()
	g.FlatCoords = flatCoords
	return g
}

// Area returns the area of g, i.e. zero.
func (g *LineString) Area() float64 {
	return 0
}

// Clone returns a copy of g that does not alias g.
func (g *LineString) Clone() *LineString {
	return deriveCloneLineString(g)
}

// Interpolate returns the index and delta of val in dimension dim.
func (g *LineString) Interpolate(val float64, dim int) (int, float64) {
	n := len(g.FlatCoords)
	if n == 0 {
		panic("geom: empty linestring")
	}
	if val <= g.FlatCoords[dim] {
		return 0, 0
	}
	if g.FlatCoords[n-g.Stride+dim] <= val {
		return (n - 1) / g.Stride, 0
	}
	low := 0
	high := n / g.Stride
	for low < high {
		mid := (low + high) / 2
		if val < g.FlatCoords[mid*g.Stride+dim] {
			high = mid
		} else {
			low = mid + 1
		}
	}
	low--
	val0 := g.FlatCoords[low*g.Stride+dim]
	if val == val0 {
		return low, 0
	}
	val1 := g.FlatCoords[(low+1)*g.Stride+dim]
	return low, (val - val0) / (val1 - val0)
}

// Length returns the length of g.
func (g *LineString) Length() float64 {
	return length1(g.FlatCoords, 0, len(g.FlatCoords), g.Stride)
}

// MustSetCoords is like SetCoords but it panics on any error.
func (g *LineString) MustSetCoords(coords []Coord) *LineString {
	Must(g.SetCoords(coords))
	return g
}

// SetCoords sets the coordinates of g.
func (g *LineString) SetCoords(coords []Coord) (*LineString, error) {
	if err := g.setCoords(coords); err != nil {
		return nil, err
	}
	return g, nil
}

// SetSRID sets the SRID of g.
func (g *LineString) SetSRID(srid int) *LineString {
	g.Srid = srid
	return g
}

// SubLineString returns a LineString from starts at index start and stops at
// index stop of g. The returned LineString aliases g.
func (g *LineString) SubLineString(start, stop int) *LineString {
	return NewLineStringFlat(g.Layout, g.FlatCoords[start*g.Stride:stop*g.Stride])
}

// Swap swaps the values of g and g2.
func (g *LineString) Swap(g2 *LineString) {
	*g, *g2 = *g2, *g
}
