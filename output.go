package geohashtree

/*
This files handles abstracting the core function geohashtree.go uses.
*/

import (
    "fmt"
    "os"
    "strings"
    "encoding/json"
    "github.com/bcicen/jstream"
    "github.com/paulmach/go.geojson"
)

// creates a string that can be appended to a csv file
func CleanOutput(outputgeohashs []string, idstring string, minval int) string {
	// creating stringlist
	mymap := map[string]string{}
	newlist := []string{}
	for i, val := range outputgeohashs {
		currentval := val
		for len(currentval) != minval {
			_, boolval := mymap[currentval[:len(currentval)-1]]
			if !boolval {
				mymap[currentval[:len(currentval)-1]] = ""
				newlist = append(newlist, fmt.Sprintf("%s,%s", currentval[:len(currentval)-1], "-1"))
			}
			currentval = currentval[:len(currentval)-1]
		}
		outputgeohashs[i] = fmt.Sprintf("%s,%s", val, idstring)
	}
	return strings.Join(append(newlist, outputgeohashs...), "\n") + "\n"
}

type IndexOutput struct {
	Min, Max      int
	File          *os.File
	FileName      string
	TotalPolygons int
}

// creates a csv to start appending to
func CreateCSV(filename string, minp, maxp int) (*IndexOutput, error) {
	os.Create(filename)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	file.WriteString(fmt.Sprintf("GEOHASH,ID\nmin,%d\nmax,%d\ndummy,-1\n", minp, maxp))
	return &IndexOutput{
		Min:           minp,
		Max:           maxp,
		File:          file,
		FileName:      filename,
		TotalPolygons: 0,
	}, err
}

func (output *IndexOutput) AddFeature(feature *geojson.Feature, field string) string {
	val, boolval := feature.Properties[field]
	var valstr string
	if !boolval {
		return ""
	} else {
		valstr, boolval = val.(string)
		if !boolval {
			return ""
		}
	}
	output.TotalPolygons++
	if feature.Geometry.Type == "Polygon" {
		return CleanOutput(
			MakePolygonIndex(feature.Geometry.Polygon, output.Min, output.Max),
			valstr,
			output.Min,
		)
	} else if feature.Geometry.Type == "MultiPolygon" {
		totaloutput := []string{}
		for _, polygon := range feature.Geometry.MultiPolygon {
			totaloutput = append(totaloutput, MakePolygonIndex(polygon, output.Min, output.Max)...)
		}
		return CleanOutput(totaloutput, valstr, output.Min)
	}
	return ""
}

// creates an index from geojson and dumps it into a csv
func IndexFromGeoJSON(filename string, outfilename string, minp, maxp int, geojsonfield string) error {
    infile, err := os.Open(filename)
	if err != nil {
		return err
	}

    defer infile.Close()
	decoder := jstream.NewDecoder(infile, 2)
	output, err := CreateCSV(outfilename, minp, maxp)
	if err != nil {
		return err
	}

    i := 0
	for mv := range decoder.Stream() {
        bs, err := json.Marshal(mv.Value)
        if err != nil {
            return err
        }
        feature, err := geojson.UnmarshalFeature(bs)
        if err != nil {
            return err
        }

        val := output.AddFeature(feature, geojsonfield)
        output.File.WriteString(val)
        i++
        fmt.Printf("%d features written to output csv.\n", i)
	}

	if err := decoder.Err(); err != nil {
		return err
	}

	fmt.Printf("\nFinished making output csv: %s\n", outfilename)
	return err
}
