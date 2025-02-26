package bigxy_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/don4get/go-geom"
	"github.com/don4get/go-geom/bigxy"
	"github.com/don4get/go-geom/xy/orientation"
)

func TestOrientationIndex(t *testing.T) {
	for i, testData := range []struct {
		vectorOrigin, vectorEnd, point geom.Coord
		result                         orientation.Type
	}{
		{
			vectorOrigin: geom.Coord{-1.0, -1.0},
			vectorEnd:    geom.Coord{1.0, 1.0},
			point:        geom.Coord{0, 0},
			result:       orientation.Collinear,
		},
		{
			vectorOrigin: geom.Coord{1.0, 1.0},
			vectorEnd:    geom.Coord{-1.0, -1.0},
			point:        geom.Coord{0, 0},
			result:       orientation.Collinear,
		},
		{
			vectorOrigin: geom.Coord{10.0, 10.0},
			vectorEnd:    geom.Coord{20.0, 20.0},
			point:        geom.Coord{10.0, 20.0},
			result:       orientation.CounterClockwise,
		},
		{
			vectorOrigin: geom.Coord{10.0, 10.0},
			vectorEnd:    geom.Coord{20.0, 20.0},
			point:        geom.Coord{20.0, 10.0},
			result:       orientation.Clockwise,
		},
		{
			vectorOrigin: geom.Coord{10.0, 20.0},
			vectorEnd:    geom.Coord{20.0, 10.0},
			point:        geom.Coord{10.0, 10.0},
			result:       orientation.Clockwise,
		},
		{
			vectorOrigin: geom.Coord{10.0, 20.0},
			vectorEnd:    geom.Coord{20.0, 10.0},
			point:        geom.Coord{20.0, 20.00},
			result:       orientation.CounterClockwise,
		},
		{
			vectorOrigin: geom.Coord{-71.104126, 42.314675},
			vectorEnd:    geom.Coord{-17.104138, 42.314732},
			point:        geom.Coord{-17.1041375307579, 42.3147318674446},
			result:       orientation.Clockwise,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			orientationIndex := bigxy.OrientationIndex(testData.vectorOrigin, testData.vectorEnd, testData.point)
			assert.Equal(t, testData.result, orientationIndex)
		})
	}
}

func TestIntersection(t *testing.T) {
	for _, tc := range []struct {
		desc                                       string
		line1Start, line1End, line2Start, line2End geom.Coord
		result                                     geom.Coord
	}{
		{
			desc:       "Plus",
			line1Start: geom.Coord{0, 1},
			line1End:   geom.Coord{0, -1},
			line2Start: geom.Coord{-1, 0},
			line2End:   geom.Coord{1, 0},
			result:     geom.Coord{0, 0},
		},
		{
			desc:       "X",
			line1Start: geom.Coord{0, 1},
			line1End:   geom.Coord{1, 0},
			line2Start: geom.Coord{1, 1},
			line2End:   geom.Coord{0, 0},
			result:     geom.Coord{0.5, 0.5},
		},
		{
			desc:       "Ends Touch",
			line1Start: geom.Coord{0, 1},
			line1End:   geom.Coord{1, 1},
			line2Start: geom.Coord{1, 1},
			line2End:   geom.Coord{1, 0},
			result:     geom.Coord{1, 1},
		},
		{
			desc:       "Close Parallel",
			line1Start: geom.Coord{0, 0},
			line1End:   geom.Coord{1, 0},
			line2Start: geom.Coord{1e-66, 1e-66},
			line2End:   geom.Coord{1 + 1e-66, 1e-66},
			result:     geom.Coord{math.Inf(1), math.Inf(1)}, // response when not possible to calculate
		},
		{
			desc:       "No Intersection",
			line1Start: geom.Coord{0, 0},
			line1End:   geom.Coord{1, 0},
			line2Start: geom.Coord{2, 1},
			line2End:   geom.Coord{2, 2},
			result:     geom.Coord{2, 0},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			calculatedIntersection := bigxy.Intersection(tc.line1Start, tc.line1End, tc.line2Start, tc.line2End)
			assert.Equal(t, tc.result, calculatedIntersection)
		})
	}
}
