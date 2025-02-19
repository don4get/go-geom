package geom

// A MultiPoint is a collection of Points.
type MultiPoint struct {
	// To represent an MultiPoint that allows EMPTY elements, e.g.
	// MULTIPOINT ( EMPTY, POINT(1.0 1.0), EMPTY), we have to allow
	// record ends. If there is an empty point, ends[i] == ends[i-1].
	Geom2
}

// NewMultiPoint returns a new, empty, MultiPoint.
func NewMultiPoint(layout Layout) *MultiPoint {
	return NewMultiPointFlat(layout, nil)
}

// NewMultiPointFlatOption represents an option that can be passed into
// NewMultiPointFlat.
type NewMultiPointFlatOption func(*MultiPoint)

// NewMultiPointFlatOptionWithEnds allows passing ends to NewMultiPointFlat,
// which allows the representation of empty points.
func NewMultiPointFlatOptionWithEnds(ends []int) NewMultiPointFlatOption {
	return func(mp *MultiPoint) {
		mp.Ends = ends
	}
}

// NewMultiPointFlat returns a new MultiPoint with the given flat coordinates.
// Assumes no points are empty by default. Use `NewMultiPointFlatOptionWithEnds`
// to specify empty points.
func NewMultiPointFlat(
	layout Layout, flatCoords []float64, opts ...NewMultiPointFlatOption,
) *MultiPoint {
	g := new(MultiPoint)
	g.Layout = layout
	g.Stride = layout.Stride()
	g.FlatCoords = flatCoords
	for _, opt := range opts {
		opt(g)
	}
	// If no ends are provided, assume all points are non empty.
	if g.Ends == nil && len(g.FlatCoords) > 0 {
		numCoords := 0
		if g.Stride > 0 {
			numCoords = len(flatCoords) / g.Stride
		}
		g.Ends = make([]int, numCoords)
		for i := range numCoords {
			g.Ends[i] = (i + 1) * g.Stride
		}
	}
	return g
}

// Area returns the area of g, i.e. zero.
func (g *MultiPoint) Area() float64 {
	return 0
}

// Clone returns a deep copy.
func (g *MultiPoint) Clone() *MultiPoint {
	return deriveCloneMultiPoint(g)
}

// Length returns zero.
func (g *MultiPoint) Length() float64 {
	return 0
}

// MustSetCoords sets the coordinates and panics on any error.
func (g *MultiPoint) MustSetCoords(coords []Coord) *MultiPoint {
	Must(g.SetCoords(coords))
	return g
}

// Coord returns the ith coord of g.
func (g *MultiPoint) Coord(i int) Coord {
	before := 0
	if i > 0 {
		before = g.Ends[i-1]
	}
	if g.Ends[i] == before {
		return nil
	}
	return g.FlatCoords[before:g.Ends[i]]
}

// SetCoords sets the coordinates.
func (g *MultiPoint) SetCoords(coords []Coord) (*MultiPoint, error) {
	g.FlatCoords = nil
	g.Ends = nil
	for _, c := range coords {
		if c != nil {
			var err error
			g.FlatCoords, err = deflate0(g.FlatCoords, c, g.Stride)
			if err != nil {
				return nil, err
			}
		}
		g.Ends = append(g.Ends, len(g.FlatCoords))
	}
	return g, nil
}

// Coords unpacks and returns all of g's coordinates.
func (g *MultiPoint) Coords() []Coord {
	coords1 := make([]Coord, len(g.Ends))
	offset := 0
	prevEnd := 0
	for i, end := range g.Ends {
		if end != prevEnd {
			coords1[i] = inflate0(g.FlatCoords, offset, offset+g.Stride, g.Stride)
			offset += g.Stride
		}
		prevEnd = end
	}
	return coords1
}

// NumCoords returns the number of coordinates in g.
func (g *MultiPoint) NumCoords() int {
	return len(g.Ends)
}

// SetSRID sets the SRID of g.
func (g *MultiPoint) SetSRID(srid int) *MultiPoint {
	g.Srid = srid
	return g
}

// NumPoints returns the number of Points.
func (g *MultiPoint) NumPoints() int {
	return len(g.Ends)
}

// Point returns the ith Point.
func (g *MultiPoint) Point(i int) *Point {
	coord := g.Coord(i)
	if coord == nil {
		return NewPointEmpty(g.Layout)
	}
	return NewPointFlat(g.Layout, coord)
}

// Push appends a point.
func (g *MultiPoint) Push(p *Point) error {
	if p.Layout != g.Layout {
		return ErrLayoutMismatch{Got: p.Layout, Want: g.Layout}
	}
	if !p.IsEmpty() {
		g.FlatCoords = append(g.FlatCoords, p.FlatCoords...)
	}
	g.Ends = append(g.Ends, len(g.FlatCoords))
	return nil
}

// Swap swaps the values of g and g2.
func (g *MultiPoint) Swap(g2 *MultiPoint) {
	*g, *g2 = *g2, *g
}
