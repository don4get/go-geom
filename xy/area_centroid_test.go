package xy_test

import (
	"math"
	"testing"

	"github.com/don4get/go-geom"
	"github.com/don4get/go-geom/geomtest"
	"github.com/don4get/go-geom/xy"
	"github.com/don4get/go-geom/xy/internal"
)

func TestAreaCentroidCalculator_GetCentroid_NoGeomsAdded(t *testing.T) {
	calculator := xy.NewAreaCentroidCalculator(geom.XY)
	centroid := calculator.GetCentroid()
	if !centroid.Equal(geom.XY, geom.Coord{math.NaN(), math.NaN()}) {
		t.Errorf("centroid with no coords added should return the [NaN NaN] coord but was: %v", centroid)
	}
}

var polygonTestData = []struct {
	polygons                   []*geom.Polygon
	areaCentroid, lineCentroid geom.Coord
}{
	{
		polygons: []*geom.Polygon{
			geom.NewPolygonFlat(geom.XY, []float64{0, 0, 2, 0, 2, 2, 0, 2, 0, 0}, []int{10}),
		},
		areaCentroid: geom.Coord{1, 1},
		lineCentroid: geom.Coord{1, 1},
	},
	{
		polygons: []*geom.Polygon{
			geom.NewPolygonFlat(geom.XY, []float64{
				0, 0, 2, 0, 2, 2, 0, 2, 0, 0,
				0.5, 0.5, 0.75, 0.5, 0.75, 0.75, 0.5, 0.75, 0.5, 0.5,
				1.25, 1.25, 1.5, 1.25, 1.5, 1.5, 1.25, 1.5, 1.25, 1.25,
			}, []int{10, 20, 30}),
		},
		areaCentroid: geom.Coord{1, 1},
		lineCentroid: geom.Coord{1, 1},
	},
	{
		polygons: []*geom.Polygon{
			geom.NewPolygonFlat(geom.XY, []float64{-100, 100, 100, 100, 10, -100, -10, -100, -100, 100}, []int{10}),
		},
		areaCentroid: geom.Coord{0.0, 27.272727272727273},
		lineCentroid: geom.Coord{0.0, 27.329280498653272},
	},
	{
		polygons: []*geom.Polygon{
			geom.NewPolygonFlat(geom.XY, []float64{-100, 100, 100, 100, 10, -100, -10, -100, -100, 100}, []int{10}),
			geom.NewPolygonFlat(geom.XY, []float64{-100, -100, 100, -100, 10, 100, -10, 100, -100, -100}, []int{10}),
		},
		areaCentroid: geom.Coord{0.0, 0.0},
		lineCentroid: geom.Coord{0, 0},
	},
	{
		polygons: []*geom.Polygon{
			geom.NewPolygonFlat(geom.XY, internal.TestRing.GetFlatCoords(), []int{internal.TestRing.NumCoords() * 2}),
		},
		areaCentroid: geom.Coord{-53.10266611446687, 42.314777901050384},
		lineCentroid: geom.Coord{-44.10405031184597, 42.3149062174918},
	},
}

func TestAreaGetCentroid(t *testing.T) {
	for i, tc := range polygonTestData {
		centroid := xy.PolygonsCentroid(tc.polygons[0], tc.polygons[1:]...)

		if !geomtest.CoordsEqualRel(tc.areaCentroid, centroid, 1e-15) {
			t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.areaCentroid, centroid)
		}

		coords := []float64{}
		endss := [][]int{}
		lastEnd := 0
		for _, p := range tc.polygons {
			coords = append(coords, p.GetFlatCoords()...)
			ends := append([]int{}, p.GetEnds()...)

			for i := range p.GetEnds() {
				ends[i] += lastEnd
			}
			endss = append(endss, ends)
			lastEnd = len(coords)
		}

		layout := tc.polygons[0].GetLayout()
		multiPolygon := geom.NewMultiPolygonFlat(layout, coords, endss)
		centroid = xy.MultiPolygonCentroid(multiPolygon)

		if !geomtest.CoordsEqualRel(tc.areaCentroid, centroid, 1e-15) {
			t.Errorf("Test '%v' failed: expected centroid for multipolygon to be\n%v but was \n%v", i+1, tc.areaCentroid, centroid)
		}
	}
}
