package xy_test

import (
	"fmt"

	"github.com/don4get/go-geom"
	"github.com/don4get/go-geom/xy"
)

func ExampleConvexHull() {
	polygon := geom.NewLineStringFlat(geom.XY, []float64{1, 1, 3, 3, 4, 4, 2, 5})

	convexHull := xy.ConvexHull(polygon)

	fmt.Println(convexHull.GetFlatCoords())
	// Output: [1 1 2 5 4 4 1 1]
}

func ExampleConvexHullFlat() {
	polygon := geom.NewLineStringFlat(geom.XY, []float64{1, 1, 3, 3, 4, 4, 2, 5})

	convexHull := xy.ConvexHullFlat(polygon.GetLayout(), polygon.GetFlatCoords())

	fmt.Println(convexHull.GetFlatCoords())
	// Output: [1 1 2 5 4 4 1 1]
}
