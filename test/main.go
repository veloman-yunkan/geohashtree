package main

import (
    "fmt"
    "os"
    ".."
)

func run() error {
    return geohashtree.IndexFromGeoJSON("input.geojson", "output.csv", 5, 9, "id")
}

func main() {
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}
