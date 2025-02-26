// parse-igc parsers all IGC files passed on the command line and prints any
// errors encountered.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/don4get/go-geom/encoding/igc"
)

func parseIGC(filename string) (*igc.T, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return igc.Read(f)
}

//nolint:unparam
func run() error {
	flag.Parse()
	for _, arg := range flag.Args() {
		if _, err := parseIGC(arg); err != nil {
			fmt.Printf("%s: %v\n", arg, err)
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
