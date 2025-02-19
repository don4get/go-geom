package geom

// A MultiPolygon is a collection of Polygons.
type MultiPolygon struct {
	geom3
}

// NewMultiPolygon returns a new MultiPolygon with no Polygons.
func NewMultiPolygon(layout Layout) *MultiPolygon {
	return NewMultiPolygonFlat(layout, nil, nil)
}

// NewMultiPolygonFlat returns a new MultiPolygon with the given flat coordinates.
func NewMultiPolygonFlat(layout Layout, flatCoords []float64, endss [][]int) *MultiPolygon {
	g := new(MultiPolygon)
	g.Layout = layout
	g.Stride = layout.Stride()
	g.FlatCoords = flatCoords
	g.Endss = endss
	return g
}

// Area returns the sum of the area of the individual Polygons.
func (g *MultiPolygon) Area() float64 {
	return doubleArea3(g.FlatCoords, 0, g.Endss, g.Stride) / 2
}

// Clone returns a deep copy.
func (g *MultiPolygon) Clone() *MultiPolygon {
	return deriveCloneMultiPolygon(g)
}

// Length returns the sum of the perimeters of the Polygons.
func (g *MultiPolygon) Length() float64 {
	return length3(g.FlatCoords, 0, g.Endss, g.Stride)
}

// MustSetCoords sets the coordinates and panics on any error.
func (g *MultiPolygon) MustSetCoords(coords [][][]Coord) *MultiPolygon {
	Must(g.SetCoords(coords))
	return g
}

// NumPolygons returns the number of Polygons.
func (g *MultiPolygon) NumPolygons() int {
	return len(g.Endss)
}

// Polygon returns the ith Polygon.
func (g *MultiPolygon) Polygon(i int) *Polygon {
	if len(g.Endss[i]) == 0 {
		return NewPolygon(g.Layout)
	}
	// Find the offset from the previous non-empty polygon element.
	offset := 0
	lastNonEmptyIdx := i - 1
	for lastNonEmptyIdx >= 0 {
		ends := g.Endss[lastNonEmptyIdx]
		if len(ends) > 0 {
			offset = ends[len(ends)-1]
			break
		}
		lastNonEmptyIdx--
	}
	ends := make([]int, len(g.Endss[i]))
	if offset == 0 {
		copy(ends, g.Endss[i])
	} else {
		for j, end := range g.Endss[i] {
			ends[j] = end - offset
		}
	}
	return NewPolygonFlat(g.Layout, g.FlatCoords[offset:g.Endss[i][len(g.Endss[i])-1]], ends)
}

// Push appends a Polygon.
func (g *MultiPolygon) Push(p *Polygon) error {
	if p.Layout != g.Layout {
		return ErrLayoutMismatch{Got: p.Layout, Want: g.Layout}
	}
	offset := len(g.FlatCoords)
	var ends []int
	if len(p.Ends) > 0 {
		ends = make([]int, len(p.Ends))
		if offset == 0 {
			copy(ends, p.Ends)
		} else {
			for i, end := range p.Ends {
				ends[i] = end + offset
			}
		}
	}
	g.FlatCoords = append(g.FlatCoords, p.FlatCoords...)
	g.Endss = append(g.Endss, ends)
	return nil
}

// SetCoords sets the coordinates.
func (g *MultiPolygon) SetCoords(coords [][][]Coord) (*MultiPolygon, error) {
	if err := g.setCoords(coords); err != nil {
		return nil, err
	}
	return g, nil
}

// SetSRID sets the SRID of g.
func (g *MultiPolygon) SetSRID(srid int) *MultiPolygon {
	g.Srid = srid
	return g
}

// Swap swaps the values of g and g2.
func (g *MultiPolygon) Swap(g2 *MultiPolygon) {
	*g, *g2 = *g2, *g
}
