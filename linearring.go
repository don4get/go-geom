package geom

// A LinearRing is a linear ring.
type LinearRing struct {
	Geom1
}

// NewLinearRing returns a new LinearRing with no coordinates.
func NewLinearRing(layout Layout) *LinearRing {
	return NewLinearRingFlat(layout, nil)
}

// NewLinearRingFlat returns a new LinearRing with the given flat coordinates.
func NewLinearRingFlat(layout Layout, flatCoords []float64) *LinearRing {
	g := new(LinearRing)
	g.Layout = layout
	g.Stride = layout.Stride()
	g.FlatCoords = flatCoords
	return g
}

// Area returns the area.
func (g *LinearRing) Area() float64 {
	return doubleArea1(g.FlatCoords, 0, len(g.FlatCoords), g.Stride) / 2
}

// Clone returns a deep copy.
func (g *LinearRing) Clone() *LinearRing {
	return deriveCloneLinearRing(g)
}

// Length returns the length of the perimeter.
func (g *LinearRing) Length() float64 {
	return length1(g.FlatCoords, 0, len(g.FlatCoords), g.Stride)
}

// MustSetCoords sets the coordinates and panics if there is any error.
func (g *LinearRing) MustSetCoords(coords []Coord) *LinearRing {
	Must(g.SetCoords(coords))
	return g
}

// SetCoords sets the coordinates.
func (g *LinearRing) SetCoords(coords []Coord) (*LinearRing, error) {
	if err := g.setCoords(coords); err != nil {
		return nil, err
	}
	return g, nil
}

// SetSRID sets the SRID of g.
func (g *LinearRing) SetSRID(srid int) *LinearRing {
	g.Srid = srid
	return g
}

// Swap swaps the values of g and g2.
func (g *LinearRing) Swap(g2 *LinearRing) {
	*g, *g2 = *g2, *g
}
