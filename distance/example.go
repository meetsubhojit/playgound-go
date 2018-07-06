package main

import (
	"flag"
	"fmt"
	"github.com/TheGUNNER13/playgound-go/distance/finder"
	"os"
	"time"
)

var (
	filePath        string
	readWorkerCount int
	findWorkerCount int
	findBathSize    int
	officeLat       float64
	officeLng       float64
	howMany         int
)

func init() {
	gopath := os.Getenv("GOPATH")
	flag.StringVar(&filePath, "filePath", fmt.Sprintf("%s/%s", gopath, "src/github.com/TheGUNNER13/playgound-go/distance/resources/geoData.csv"), "The path to the csv file which contains all the coordinates")
	flag.IntVar(&readWorkerCount, "readWorkerCount", 20, "The number of threads to be spawned to concurrently find the distances by the coordinates")
	flag.IntVar(&findWorkerCount, "findWorkerCount", 100, "The number of threads to be spawned to find min/max distance")
	flag.IntVar(&findBathSize, "findBathSize", 100, "The batch size for the heap to be built")
	flag.Float64Var(&officeLat, "officeLat", 51.925146, "The office latitude")
	flag.Float64Var(&officeLng, "officeLng", 4.478617, "The office longitude")
	flag.IntVar(&howMany, "howMany", 5, "How many closest/farthest point to show")
}

func main() {
	flag.Parse()
	now := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error occured while reading file, ", err.Error())
		os.Exit(1)
	}

	readChan := make(chan finder.Coordinate, findBathSize*findWorkerCount)

	go finder.Reader(file, officeLat, officeLng, readChan, readWorkerCount)
	min, max := finder.GetTopMinMax(readChan, findWorkerCount, findBathSize, howMany)

	fmt.Println("The closest points (in KMs) are")
	printPoints(min)
	fmt.Println("\n\nThe farthest points (in KMs) are")
	printPoints(max)

	fmt.Println("\nTime taken ", time.Now().Sub(now))
}
func printPoints(ar []finder.Coordinate) {
	for _, a := range ar {
		fmt.Println("ID ", a.Id, " Distance from office ", a.Distance)
	}
}
