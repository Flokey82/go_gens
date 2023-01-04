package main

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/Flokey82/go_gens/genworldvoronoi"
	"github.com/gorilla/mux"
)

var worldmap *genworldvoronoi.Map

var (
	seed      int64   = 1234
	numPlates int     = 25
	numPoints int     = 40000
	jitter    float64 = 0.0
)

func init() {
	flag.Int64Var(&seed, "seed", seed, "the world seed")
	flag.IntVar(&numPlates, "num_plates", numPlates, "number of plates")
	flag.IntVar(&numPoints, "num_points", numPoints, "number of points")
	flag.Float64Var(&jitter, "jitter", jitter, "jitter")
}

func main() {
	flag.Parse()

	// Initialize the planet.
	sp, err := genworldvoronoi.NewMap(seed, numPlates, numPoints, jitter)
	if err != nil {
		log.Fatal(err)
	}
	worldmap = sp

	// Start the server.
	router := mux.NewRouter()
	router.HandleFunc("/tiles/{z}/{x}/{y}", tileHandler)
	router.HandleFunc("/geojson_cities/{z}/{la1}/{lo1}/{la2}/{lo2}", geoJSONCitiesHandler)
	router.HandleFunc("/geojson_borders/{z}/{la1}/{lo1}/{la2}/{lo2}", geoJSONBorderHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe(":3333", router))
}

func parseBoundingBox(req *http.Request) (la1, lo1, la2, lo2 float64, z int, err error) {
	// Get the tile coordinates and zoom level.
	vars := mux.Vars(req)
	la1, err = strconv.ParseFloat(vars["la1"], 64)
	if err != nil {
		return
	}
	la2, err = strconv.ParseFloat(vars["la2"], 64)
	if err != nil {
		return
	}
	lo1, err = strconv.ParseFloat(vars["lo1"], 64)
	if err != nil {
		return
	}
	lo2, err = strconv.ParseFloat(vars["lo2"], 64)
	if err != nil {
		return
	}
	z, err = strconv.Atoi(vars["z"])
	if err != nil {
		return
	}
	return
}

func geoJSONCitiesHandler(res http.ResponseWriter, req *http.Request) {
	// Get the tile coordinates and zoom level.
	tileLa1, tileLo1, tileLa2, tileLo2, tileZ, err := parseBoundingBox(req)
	if err != nil {
		panic(err)
	}
	data, err := worldmap.GetGeoJSONCities(tileLa1, tileLo1, tileLa2, tileLo2, tileZ)
	if err != nil {
		panic(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Content-Length", strconv.Itoa(len(data)))
	res.Write(data)
}

func geoJSONBorderHandler(res http.ResponseWriter, req *http.Request) {
	// Get the tile coordinates and zoom level.
	tileLa1, tileLo1, tileLa2, tileLo2, tileZ, err := parseBoundingBox(req)
	if err != nil {
		panic(err)
	}
	data, err := worldmap.GetGeoJSONBorders(tileLa1, tileLo1, tileLa2, tileLo2, tileZ)
	if err != nil {
		panic(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Content-Length", strconv.Itoa(len(data)))
	res.Write(data)
}

func tileHandler(res http.ResponseWriter, req *http.Request) {
	// Get the url parameter 'd'.
	d := req.URL.Query().Get("d")
	if d == "" {
		d = "0"
	}
	displayMode, err := strconv.Atoi(d)
	if err != nil {
		panic(err)
	}

	// get the url parameter 'wind'.
	wind := req.URL.Query().Get("wind")
	if wind == "" {
		wind = "false"
	}

	// Get the tile coordinates and zoom level.
	vars := mux.Vars(req)
	tileX, err := strconv.Atoi(vars["x"])
	if err != nil {
		panic(err)
	}
	tileY, err := strconv.Atoi(vars["y"])
	if err != nil {
		panic(err)
	}
	tileZ, err := strconv.Atoi(vars["z"])
	if err != nil {
		panic(err)
	}

	// Get the tile image.
	img := worldmap.GetTile(tileX, tileY, tileZ, displayMode, wind == "true")
	writeImage(res, &img)
}

// writeImage writes the image to the response writer.
func writeImage(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, *img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
