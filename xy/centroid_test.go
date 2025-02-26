package xy_test

import (
	"testing"

	"github.com/don4get/go-geom"
	"github.com/don4get/go-geom/geomtest"
	"github.com/don4get/go-geom/xy"
	"github.com/don4get/go-geom/xy/internal"
)

func TestCentroid(t *testing.T) {
	for _, tc := range []struct {
		id       int
		geometry geom.T
		centroid geom.Coord
	}{
		{
			id:       1,
			geometry: geom.NewPolygonFlat(geom.XY, []float64{-100, 100, 100, 100, 10, -100, -10, -100, -100, 100}, []int{10}),
			centroid: geom.Coord{0.0, 27.272727272727273},
		},
		{
			id:       2,
			geometry: geom.NewMultiPolygonFlat(geom.XY, []float64{-100, 100, 100, 100, 10, -100, -10, -100, -100, 100}, [][]int{{10}}),
			centroid: geom.Coord{0.0, 27.272727272727273},
		},
		{
			id:       3,
			geometry: geom.NewLineStringFlat(internal.TestRing.GetLayout(), internal.TestRing.GetFlatCoords()),
			centroid: geom.Coord{-44.10405031184597, 42.3149062174918},
		},
		{
			id:       4,
			geometry: geom.NewMultiLineStringFlat(internal.TestRing.GetLayout(), internal.TestRing.GetFlatCoords(), []int{len(internal.TestRing.GetFlatCoords())}),
			centroid: geom.Coord{-44.10405031184597, 42.3149062174918},
		},
		{
			id:       5,
			geometry: internal.TestRing,
			centroid: geom.Coord{-44.10405031184597, 42.3149062174918},
		},
		{
			id:       6,
			geometry: geom.NewPointFlat(geom.XY, []float64{2, 2}),
			centroid: geom.Coord{2, 2},
		},
		{
			id:       7,
			geometry: geom.NewMultiPointFlat(geom.XY, []float64{0, 0, 2, 2}),
			centroid: geom.Coord{1, 1},
		},
	} {
		calculated, err := xy.Centroid(tc.geometry)

		if !geomtest.CoordsEqualRel(calculated, tc.centroid, 1e-15) {
			t.Errorf("Test %v failed.  Expected \n\t%v but got \n\t%v", tc.id, tc.centroid, calculated)
		}
		if err != nil {
			t.Errorf("Test %v failed.  Expected \n\t%v but got \n\t%v", tc.id, tc.centroid, calculated)
		}
	}
}
