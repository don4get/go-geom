package xy_test

import (
	"math"
	"testing"

	"github.com/don4get/go-geom"
	"github.com/don4get/go-geom/geomtest"
	"github.com/don4get/go-geom/xy"
	"github.com/don4get/go-geom/xy/internal"
)

func TestLineCentroidCalculator_GetCentroid_NoGeomsAdded(t *testing.T) {
	calculator := xy.NewLineCentroidCalculator(geom.XY)
	centroid := calculator.GetCentroid()
	if !centroid.Equal(geom.XY, geom.Coord{math.NaN(), math.NaN()}) {
		t.Errorf("centroid with no coords added should return the [NaN NaN] coord but was: %v", centroid)
	}
}

type lineDataType struct {
	lines        []*geom.LineString
	lineCentroid geom.Coord
}

var lineTestData = []lineDataType{
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 0}),
		},
		lineCentroid: geom.Coord{5, 0},
	},
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 10}),
		},
		lineCentroid: geom.Coord{5, 5},
	},
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 0}),
			geom.NewLineStringFlat(geom.XY, []float64{0, 10, 10, 10}),
		},
		lineCentroid: geom.Coord{5, 5},
	},
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 0}),
			geom.NewLineStringFlat(geom.XY, []float64{0, 10, 5, 10}),
		},
		lineCentroid: geom.Coord{4.166666666666667, 3.3333333333333335},
	},
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 0, 10, 10, 0, 0}),
		},
		lineCentroid: geom.Coord{6.464466094067262, 3.5355339059327378},
	},
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(internal.TestRing.GetLayout(), internal.TestRing.GetFlatCoords()),
		},
		lineCentroid: geom.Coord{-44.10405031184597, 42.3149062174918},
	},
}

func TestLineGetCentroidLines(t *testing.T) {
	for i, tc := range lineTestData {
		verifyLineCentroid(t, i, tc)
		verifyMultiLineCentroid(t, i, tc)
		verifyLinearRingsCentroid(t, i, tc)
	}
}

func verifyLineCentroid(t *testing.T, i int, tc lineDataType) {
	t.Helper()

	if centroid := xy.LinesCentroid(tc.lines[0], tc.lines[1:]...); !geomtest.CoordsEqualRel(tc.lineCentroid, centroid, 1e-15) {
		t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.lineCentroid, centroid)
	}
}

func verifyMultiLineCentroid(t *testing.T, i int, tc lineDataType) {
	t.Helper()
	coords := []float64{}
	ends := []int{}
	for _, p := range tc.lines {
		coords = append(coords, p.GetFlatCoords()...)
		ends = append(ends, len(coords))
	}

	layout := tc.lines[0].GetLayout()
	multiPolygon := geom.NewMultiLineStringFlat(layout, coords, ends)
	centroid := xy.MultiLineCentroid(multiPolygon)

	if !geomtest.CoordsEqualRel(tc.lineCentroid, centroid, 1e-15) {
		t.Errorf("Test '%v' failed: expected centroid for multipolygon to be\n%v but was \n%v", i+1, tc.lineCentroid, centroid)
	}
}

func verifyLinearRingsCentroid(t *testing.T, i int, tc lineDataType) {
	t.Helper()
	rings := make([]*geom.LinearRing, len(tc.lines))
	for i, p := range tc.lines {
		coords := append([]float64{}, p.GetFlatCoords()...)
		if coords[0] != coords[len(coords)-2] || coords[1] != coords[len(coords)-1] {
			coords = append(coords, coords[0], coords[1])
		}
		rings[i] = geom.NewLinearRingFlat(p.GetLayout(), coords)
	}

	centroid := xy.LinearRingsCentroid(rings[0], rings[1:]...)

	if !geomtest.CoordsEqualRel(tc.lineCentroid, centroid, 1e-15) {
		t.Errorf("Test '%v' failed: expected centroid for linear rings to be\n%v but was \n%v", i+1, tc.lineCentroid, centroid)
	}
}

func TestLineGetCentroidPolygons(t *testing.T) {
	for i, tc := range polygonTestData {
		calc := xy.NewLineCentroidCalculator(tc.polygons[0].GetLayout())
		for _, p := range tc.polygons {
			calc.AddPolygon(p)
		}
		centroid := calc.GetCentroid()

		if !geomtest.CoordsEqualRel(tc.lineCentroid, centroid, 1e-15) {
			t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.lineCentroid, centroid)
		}
	}
}
