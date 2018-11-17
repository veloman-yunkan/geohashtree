package main

import (
    "flag"
    "fmt"
    "os"
    "net/http"
    "strconv"
//    "time"
    "../.."
)

var (
	port   = flag.Int("port", 8888, "port to serve on")
	helpFlag    = flag.Bool("h", false, "display this help dialog")
)

var helpMsg = `http_server - serve geohashtree query requests over HTTP

Usage:

    http_server [options] <geohashtree_index_csv_file>

Options:
`

func help() {
	fmt.Println(helpMsg)
    flag.PrintDefaults()
}

var tree *geohashtree.GeohashTree = nil;

func handleRequest(w http.ResponseWriter, r *http.Request) {
    lon_header := r.Header["Lon"]
    if len(lon_header) == 0 {
        http.Error(w, "Longitude value not provided", 400)
        return
    }
    lat_header := r.Header["Lat"]
    if len(lat_header) == 0 {
        http.Error(w, "Lattitude value not provided", 400)
        return
    }
    lonstr := lon_header[0]
    latstr := lat_header[0]
    lon, err := strconv.ParseFloat(lonstr, 64)
    if err != nil {
        http.Error(w, "Invalid longitude: " + lonstr, 400)
        return
    }
    lat, err := strconv.ParseFloat(latstr, 64)
    if err != nil {
        http.Error(w, "Invalid lattitude: " + latstr, 400)
        return
    }
    geohash, ok := tree.Query([]float64{lon, lat})
    if ok {
        fmt.Fprintf(w, `[ "%s" ]`, geohash)
    } else {
        fmt.Fprintf(w, "[]", )
    }
}

func run() error {
    flag.Parse()
	if *helpFlag {
		help()
		os.Exit(0)
	}

    var err error
    indexfile := flag.Arg(0)
    tree, err = geohashtree.OpenGeohashTreeCSV(indexfile)
    if err != nil {
        return err
    }

//    s := &http.Server{
//        Addr:           ":" + strconv.Itoa(*port),
//        Handler:        handleRequest,
//        ReadTimeout:    1 * time.Second,
//        WriteTimeout:   1 * time.Second,
//        MaxHeaderBytes: 1 << 10,
//    }
    http.HandleFunc("/", handleRequest)
    return http.ListenAndServe(":" + strconv.Itoa(*port), nil);
}

func main() {
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}