package main

import (
    "fmt"
    "os"
    "strconv"
    "../.."
)

func run() error {
    geojsonfile := os.Args[1]
    outputfile := os.Args[2]
    mingeohashlength, err := strconv.Atoi(os.Args[3])
    if err != nil {
        return err
    }
    maxgeohashlength, err := strconv.Atoi(os.Args[4])
    if err != nil {
        return err
    }
    propname := os.Args[5]
    return geohashtree.IndexFromGeoJSON(geojsonfile,
                                        outputfile,
                                        mingeohashlength,
                                        maxgeohashlength,
                                        propname)
}

func main() {
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}
