package main

import (
    "flag"
    "fmt"
    "os"
    "../.."
)

var (
	mingeohashlength   = flag.Int("geohash-min-length", 5,
                                  "minimal length of geohashes")
	maxgeohashlength   = flag.Int("geohash-max-length", 9,
                                  "maximal length of geohashes")
	feature_prop_name  = flag.String("feature-prop-name", "id",
                                     "feature property name serving as its id")
	helpFlag    = flag.Bool("h", false, "display this help dialog")
)

var helpMsg = `make_geojson_index - build geohashtree index for GeoJSON input

Usage:

    make_geojson_index [options] <input_geojson_file> <output_csv_file>

Options:
`

func help() {
	fmt.Println(helpMsg)
    flag.PrintDefaults()
}

func run() error {
    flag.Parse()
	if *helpFlag {
		help()
		os.Exit(0)
	}

    geojsonfile := flag.Arg(0)
    outputfile := flag.Arg(1)
    return geohashtree.IndexFromGeoJSON(geojsonfile,
                                        outputfile,
                                        *mingeohashlength,
                                        *maxgeohashlength,
                                        *feature_prop_name)
}

func main() {
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}
