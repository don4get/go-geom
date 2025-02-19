package geom

// A MultiLineString is a collection of LineStrings.
type MultiLineString struct {
	Geom2
}

// NewMultiLineString returns a new MultiLineString with no LineStrings.
func NewMultiLineString(layout Layout) *MultiLineString {
	return NewMultiLineStringFlat(layout, nil, nil)
}

// NewMultiLineStringFlat returns a new MultiLineString with the given flat coordinates.
func NewMultiLineStringFlat(layout Layout, flatCoords []float64, ends []int) *MultiLineString {
	g := new(MultiLineString)
	g.Layout = layout
	g.Stride = layout.Stride()
	g.FlatCoords = flatCoords
	g.Ends = ends
	return g
}

// Area returns the area of g, i.e. 0.
func (g *MultiLineString) Area() float64 {
	return 0
}

// Clone returns a deep copy.
func (g *MultiLineString) Clone() *MultiLineString {
	return deriveCloneMultiLineString(g)
}

// Length returns the sum of the length of the LineStrings.
func (g *MultiLineString) Length() float64 {
	return length2(g.FlatCoords, 0, g.Ends, g.Stride)
}

// LineString returns the ith LineString.
func (g *MultiLineString) LineString(i int) *LineString {
	offset := 0
	if i > 0 {
		offset = g.Ends[i-1]
	}
	if offset == g.Ends[i] {
		return NewLineString(g.Layout)
	}
	return NewLineStringFlat(g.Layout, g.FlatCoords[offset:g.Ends[i]])
}

// MustSetCoords sets the coordinates and panics on any error.
func (g *MultiLineString) MustSetCoords(coords [][]Coord) *MultiLineString {
	Must(g.SetCoords(coords))
	return g
}

// NumLineStrings returns the number of LineStrings.
func (g *MultiLineString) NumLineStrings() int {
	return len(g.Ends)
}

// Push appends a LineString.
func (g *MultiLineString) Push(ls *LineString) error {
	if ls.Layout != g.Layout {
		return ErrLayoutMismatch{Got: ls.Layout, Want: g.Layout}
	}
	g.FlatCoords = append(g.FlatCoords, ls.FlatCoords...)
	g.Ends = append(g.Ends, len(g.FlatCoords))
	return nil
}

// SetCoords sets the coordinates.
func (g *MultiLineString) SetCoords(coords [][]Coord) (*MultiLineString, error) {
	if err := g.setCoords(coords); err != nil {
		return nil, err
	}
	return g, nil
}

// SetSRID sets the SRID of g.
func (g *MultiLineString) SetSRID(srid int) *MultiLineString {
	g.Srid = srid
	return g
}

// Swap swaps the values of g and g2.
func (g *MultiLineString) Swap(g2 *MultiLineString) {
	*g, *g2 = *g2, *g
}
