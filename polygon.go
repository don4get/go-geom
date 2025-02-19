package geom

// A Polygon represents a polygon as a collection of LinearRings. The first
// LinearRing is the outer boundary. Subsequent LinearRings are inner
// boundaries (holes).
type Polygon struct {
	geom2
}

// NewPolygon returns a new, empty, Polygon.
func NewPolygon(layout Layout) *Polygon {
	return NewPolygonFlat(layout, nil, nil)
}

// NewPolygonFlat returns a new Polygon with the given flat coordinates.
func NewPolygonFlat(layout Layout, flatCoords []float64, ends []int) *Polygon {
	g := new(Polygon)
	g.Layout = layout
	g.Stride = layout.Stride()
	g.FlatCoords = flatCoords
	g.Ends = ends
	return g
}

// Area returns the area.
func (g *Polygon) Area() float64 {
	return doubleArea2(g.FlatCoords, 0, g.Ends, g.Stride) / 2
}

// Clone returns a deep copy.
func (g *Polygon) Clone() *Polygon {
	return deriveClonePolygon(g)
}

// Length returns the perimter.
func (g *Polygon) Length() float64 {
	return length2(g.FlatCoords, 0, g.Ends, g.Stride)
}

// LinearRing returns the ith LinearRing.
func (g *Polygon) LinearRing(i int) *LinearRing {
	offset := 0
	if i > 0 {
		offset = g.Ends[i-1]
	}
	return NewLinearRingFlat(g.Layout, g.FlatCoords[offset:g.Ends[i]])
}

// MustSetCoords sets the coordinates and panics on any error.
func (g *Polygon) MustSetCoords(coords [][]Coord) *Polygon {
	Must(g.SetCoords(coords))
	return g
}

// NumLinearRings returns the number of LinearRings.
func (g *Polygon) NumLinearRings() int {
	return len(g.Ends)
}

// Push appends a LinearRing.
func (g *Polygon) Push(lr *LinearRing) error {
	if lr.Layout != g.Layout {
		return ErrLayoutMismatch{Got: lr.Layout, Want: g.Layout}
	}
	g.FlatCoords = append(g.FlatCoords, lr.FlatCoords...)
	g.Ends = append(g.Ends, len(g.FlatCoords))
	return nil
}

// SetCoords sets the coordinates.
func (g *Polygon) SetCoords(coords [][]Coord) (*Polygon, error) {
	if err := g.setCoords(coords); err != nil {
		return nil, err
	}
	return g, nil
}

// SetSRID sets the SRID of g.
func (g *Polygon) SetSRID(srid int) *Polygon {
	g.Srid = srid
	return g
}

// Swap swaps the values of g and g2.
func (g *Polygon) Swap(g2 *Polygon) {
	*g, *g2 = *g2, *g
}
