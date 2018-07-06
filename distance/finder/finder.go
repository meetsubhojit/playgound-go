package finder

import (
	"bufio"
	"container/heap"
	"github.com/golang/geo/s2"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
)

const earthRadius = 6371

type Coordinate struct {
	Id       string
	Distance float64
}

func Reader(r io.Reader, aLat, aLng float64, readChan chan Coordinate, workers int) {
	aLatLng := s2.LatLngFromDegrees(aLat, aLng)

	workerChan := make(chan struct{}, workers)
	wg := &sync.WaitGroup{}

	br := bufio.NewReader(r)

	for {
		b, _, err := br.ReadLine()
		if err != nil {
			wg.Wait()
			close(readChan)
			return
		}
		bString := string(b)
		workerChan <- struct{}{}
		wg.Add(1)

		go func(readChan chan Coordinate, doneChan chan struct{}, bStringCopy string, wgCopy *sync.WaitGroup) {
			defer func() {
				wgCopy.Done()
				<-doneChan
			}()

			bArray := strings.Split(bStringCopy, ",")
			if len(bArray) != 3 {
				return
			}
			bLat, err := strconv.ParseFloat(bArray[1], 64)
			if err != nil {
				return
			}
			bLng, err := strconv.ParseFloat(bArray[2], 64)
			if err != nil {
				return
			}
			bLatLng := s2.LatLngFromDegrees(bLat, bLng)
			readChan <- Coordinate{Id: bArray[0], Distance: math.Abs(float64(aLatLng.Distance(bLatLng))) * earthRadius}

		}(readChan, workerChan, bString, wg)
	}
	wg.Wait()
}
func write(writeChan chan Coordinate, ar []Coordinate) {
	for _, i := range ar {
		writeChan <- i
	}
}
func getOrderedNumbers(ar heap.Interface, howMany int) []Coordinate {
	if ar.Len() == 0 {
		return nil
	}
	heap.Init(ar)
	var out = make([]Coordinate, 0, howMany)
	c := 0
	for ar.Len() > 0 {
		c++
		out = append(out, heap.Pop(ar).(Coordinate))
		if c == howMany {
			break
		}
	}
	return out
}
func getMinMax(ar []Coordinate, howMany int) ([]Coordinate, []Coordinate) {
	var minOut, maxOut []Coordinate
	var arC = make([]Coordinate, len(ar))
	copy(arC, ar)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(arCopy []Coordinate) {
		defer wg.Done()
		var minAr MinHeapCoordinate = arCopy
		minOut = getOrderedNumbers(&minAr, howMany)
	}(ar)
	go func(arCopy []Coordinate) {
		defer wg.Done()
		var maxAr = MaxHeapCoordinate{arCopy}
		maxOut = getOrderedNumbers(&maxAr, howMany)
	}(arC)
	wg.Wait()

	return minOut, maxOut
}

func GetTopMinMax(readChan chan Coordinate, workerNum, workerCapacity, howMany int) ([]Coordinate, []Coordinate) {

	minChan := make(chan Coordinate, workerCapacity)
	maxChan := make(chan Coordinate, workerCapacity)
	minMaxWg := &sync.WaitGroup{}
	minMaxWg.Add(2)

	var minOut, maxOut []Coordinate

	go func() {
		defer minMaxWg.Done()
		var minArray = make(MinHeapCoordinate, 0, workerCapacity)
		count := 0
		for i := range minChan {
			minArray = append(minArray, i)
			count++
			if count >= workerCapacity {
				minArray = getOrderedNumbers(&minArray, howMany)
				count = 0
			}
		}
		minOut = getOrderedNumbers(&minArray, howMany)
	}()
	go func() {
		defer minMaxWg.Done()
		var maxArray = MaxHeapCoordinate{make([]Coordinate, 0, workerCapacity)}
		count := 0
		for i := range maxChan {
			maxArray.MinHeapCoordinate = append(maxArray.MinHeapCoordinate, i)
			count++
			if count >= workerCapacity {
				maxArray.MinHeapCoordinate = getOrderedNumbers(&maxArray, howMany)
				count = 0
			}
		}
		maxOut = getOrderedNumbers(&maxArray, howMany)
	}()

	ar := make([]Coordinate, 0, workerCapacity)
	c := 0
	wC := make(chan bool, workerNum)
	wg := &sync.WaitGroup{}

	do := func(ar []Coordinate, done chan bool, wg *sync.WaitGroup) {
		defer wg.Done()
		min, max := getMinMax(ar, howMany)
		write(minChan, min)
		write(maxChan, max)
		<-done
	}
	for i := range readChan {
		ar = append(ar, i)
		c++
		if c >= workerCapacity {
			wC <- true
			wg.Add(1)
			go do(ar, wC, wg)
			c = 0
			ar = make([]Coordinate, 0, workerCapacity)
		}
	}
	wC <- true
	wg.Add(1)
	go do(ar, wC, wg)

	wg.Wait()
	close(minChan)
	close(maxChan)
	minMaxWg.Wait()

	return minOut, maxOut
}
