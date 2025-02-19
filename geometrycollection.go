package geom

// A GeometryCollection is a collection of arbitrary geometries with the same
// SRID.
type GeometryCollection struct {
	layout Layout
	geoms  []T
	srid   int
}

// NewGeometryCollection returns a new empty GeometryCollection.
func NewGeometryCollection() *GeometryCollection {
	return &GeometryCollection{}
}

// Geom returns the ith geometry in g.
func (g *GeometryCollection) Geom(i int) T {
	return g.geoms[i]
}

// Geoms returns the geometries in g.
func (g *GeometryCollection) Geoms() []T {
	return g.geoms
}

// GetLayout returns the smallest layout that covers all of the layouts in g's
// geometries.
func (g *GeometryCollection) GetLayout() Layout {
	if g.layout != NoLayout {
		return g.layout
	}
	maxLayout := NoLayout
	for _, g := range g.geoms {
		switch l := g.GetLayout(); l {
		case XYZ:
			if maxLayout == XYM {
				maxLayout = XYZM
			} else if l > maxLayout {
				maxLayout = l
			}
		case XYM:
			if maxLayout == XYZ {
				maxLayout = XYZM
			} else if l > maxLayout {
				maxLayout = l
			}
		default:
			if l > maxLayout {
				maxLayout = l
			}
		}
	}
	return maxLayout
}

// NumGeoms returns the number of geometries in g.
func (g *GeometryCollection) NumGeoms() int {
	return len(g.geoms)
}

// GetStride returns the stride of g's layout.
func (g *GeometryCollection) GetStride() int {
	return g.GetLayout().Stride()
}

// GetBounds returns the bounds of all the geometries in g.
func (g *GeometryCollection) GetBounds() *Bounds {
	// FIXME this needs work for mixing layouts, e.g. XYZ and XYM
	b := NewBounds(g.GetLayout())
	for _, g := range g.geoms {
		b = b.Extend(g)
	}
	return b
}

// IsEmpty returns true if the collection is empty.
// This can return true if the GeometryCollection contains multiple Geometry objects
// which are all empty.
func (g *GeometryCollection) IsEmpty() bool {
	for _, g := range g.geoms {
		if !g.IsEmpty() {
			return false
		}
	}
	return true
}

// GetFlatCoords panics.
func (g *GeometryCollection) GetFlatCoords() []float64 {
	panic("FlatCoords() called on a GeometryCollection")
}

// GetEnds panics.
func (g *GeometryCollection) GetEnds() []int {
	panic("Ends() called on a GeometryCollection")
}

// GetEndss panics.
func (g *GeometryCollection) GetEndss() [][]int {
	panic("Endss() called on a GeometryCollection")
}

// GetSRID returns g's GetSRID.
func (g *GeometryCollection) GetSRID() int {
	return g.srid
}

// MustPush pushes gs to g. It panics on any error.
func (g *GeometryCollection) MustPush(gs ...T) *GeometryCollection {
	if err := g.Push(gs...); err != nil {
		panic(err)
	}
	return g
}

// CheckLayout checks all geometries in the collection match the given
// layout.
func (g *GeometryCollection) CheckLayout(layout Layout) error {
	if layout != NoLayout {
		for _, geom := range g.geoms {
			if geomLayout := geom.GetLayout(); geomLayout != layout {
				return ErrLayoutMismatch{
					Got:  layout,
					Want: geomLayout,
				}
			}
		}
	}
	return nil
}

// MustSetLayout sets g's layout. It panics on any error.
func (g *GeometryCollection) MustSetLayout(layout Layout) *GeometryCollection {
	if err := g.SetLayout(layout); err != nil {
		panic(err)
	}
	return g
}

// Push appends geometries.
func (g *GeometryCollection) Push(gs ...T) error {
	if g.layout != NoLayout {
		for _, geom := range gs {
			if geomLayout := geom.GetLayout(); geomLayout != g.layout {
				return ErrLayoutMismatch{
					Got:  geomLayout,
					Want: g.layout,
				}
			}
		}
	}
	g.geoms = append(g.geoms, gs...)
	return nil
}

// SetLayout sets g's layout.
func (g *GeometryCollection) SetLayout(layout Layout) error {
	if err := g.CheckLayout(layout); err != nil {
		return err
	}
	g.layout = layout
	return nil
}

// SetSRID sets g's SRID and the SRID of all its elements.
func (g *GeometryCollection) SetSRID(srid int) *GeometryCollection {
	g.srid = srid
	return g
}
