package finder

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"
)

func TestGetTopMin(t *testing.T) {

	var tests = []struct {
		name                            string
		limit                           int
		worker, workerCapacity, howMany int
	}{
		{
			name:           "Test 1",
			limit:          10,
			howMany:        5,
			worker:         2,
			workerCapacity: 4,
		},
		{
			name:           "Test 2",
			limit:          1000,
			howMany:        5,
			worker:         2,
			workerCapacity: 4,
		},
		{
			name:           "Test 3",
			limit:          100000,
			howMany:        5,
			worker:         2,
			workerCapacity: 4,
		},
	}
	for _, test := range tests {

		readChan := make(chan Coordinate, test.limit)
		var testMin MinHeapCoordinate = make([]Coordinate, 0, test.limit)
		var testMax = MaxHeapCoordinate{make([]Coordinate, 0, test.limit)}

		rand.Seed(time.Now().Unix())
		mapOfNumbers := map[int]struct{}{}
		for i := 0; i < test.limit; i++ {
			gen := func() int {
				for {
					n := rand.Intn(test.limit)
					if _, ok := mapOfNumbers[n]; ok {
						continue
					}
					mapOfNumbers[n] = struct{}{}
					return n
				}
			}()
			c := Coordinate{Id: strconv.Itoa(i), Distance: float64(gen)}
			readChan <- c
			testMin = append(testMin, c)
			testMax.MinHeapCoordinate = append(testMax.MinHeapCoordinate, c)
		}
		close(readChan)

		min, max := GetTopMinMax(readChan, test.worker, test.workerCapacity, test.howMany)
		sort.Sort(testMax)
		sort.Sort(testMin)
		if !compareArray(min, testMin[:test.howMany]) {
			t.Fatal(fmt.Sprint(min, "\n", testMin[:test.howMany]))
		}
		if !compareArray(max, testMax.MinHeapCoordinate[:test.howMany]) {
			t.Fatal(fmt.Sprint(max, "\n", testMax.MinHeapCoordinate[:test.howMany]))
		}
	}
}

func Test_getMin(t *testing.T) {

	var tests = []struct {
		name     string
		input1   *MinHeapCoordinate
		input2   int
		expArray MinHeapCoordinate
	}{
		{
			name:     "Test 1, happy path",
			input1:   &MinHeapCoordinate{{Id: "4", Distance: 4}, {Id: "3", Distance: 3}, {Id: "5", Distance: 5}, {Id: "2", Distance: 2}, {Id: "7", Distance: 7}, {Id: "4", Distance: 4}},
			input2:   2,
			expArray: []Coordinate{{Id: "2", Distance: 2}, {Id: "3", Distance: 3}},
		},
		{
			name:     "Test 2, where howMany is less than array size",
			input1:   &MinHeapCoordinate{{Id: "4", Distance: 4}, {Id: "3", Distance: 3}, {Id: "5", Distance: 5}},
			input2:   4,
			expArray: []Coordinate{{Id: "3", Distance: 3}, {Id: "4", Distance: 4}, {Id: "5", Distance: 5}},
		},
	}
	for _, test := range tests {
		expArray := *test.input1
		output := getOrderedNumbers(test.input1, test.input2)
		if len(output) != len(test.expArray) {
			t.Fatalf("Test %s, failed as the expected length of output was %d, but got %d", test.name, len(test.expArray), len(output))
		}
		sort.Sort(test.input1)
		sort.Sort(expArray)
		expArray = expArray[:len(test.expArray)]
		if !compareArray(output, expArray) {
			t.Fatalf("Test %s, failed as the expected output was %v, but got %v", test.name, expArray, output)
		}
	}
}
func Test_getMax(t *testing.T) {
	var tests = []struct {
		name     string
		input1   *MaxHeapCoordinate
		input2   int
		expArray MaxHeapCoordinate
	}{
		{
			name:     "Test 1, happy path",
			input1:   &MaxHeapCoordinate{MinHeapCoordinate{{Id: "4", Distance: 4}, {Id: "3", Distance: 3}, {Id: "5", Distance: 5}, {Id: "2", Distance: 2}, {Id: "7", Distance: 7}, {Id: "4", Distance: 4}}},
			input2:   2,
			expArray: MaxHeapCoordinate{[]Coordinate{{Id: "7", Distance: 7}, {Id: "5", Distance: 5}}},
		},
		{
			name:     "Test 2, where howMany is less than array size",
			input1:   &MaxHeapCoordinate{MinHeapCoordinate{{Id: "4", Distance: 4}, {Id: "3", Distance: 3}, {Id: "5", Distance: 5}}},
			input2:   4,
			expArray: MaxHeapCoordinate{[]Coordinate{{Id: "5", Distance: 5}, {Id: "4", Distance: 4}, {Id: "3", Distance: 3}}},
		},
	}
	for _, test := range tests {
		expArray := *test.input1
		output := getOrderedNumbers(test.input1, test.input2)
		if len(output) != len(test.expArray.MinHeapCoordinate) {
			t.Fatalf("Test %s, failed as the expected length of output was %d, but got %d", test.name, len(test.expArray.MinHeapCoordinate), len(output))
		}
		sort.Sort(test.input1)
		sort.Sort(expArray)
		expArray.MinHeapCoordinate = expArray.MinHeapCoordinate[:len(test.expArray.MinHeapCoordinate)]
		if !compareArray(output, expArray.MinHeapCoordinate) {
			t.Fatalf("Test %s, failed as the expected output was %v, but got %v", test.name, expArray, output)
		}
	}
}
func Test_getMinMax(t *testing.T) {
	var tests = []struct {
		name        string
		input1      []Coordinate
		input2      int
		expMinArray []Coordinate
		expMaxArray []Coordinate
	}{
		{
			name:        "Test 1, happy path",
			input1:      []Coordinate{{Id: "4", Distance: 4}, {Id: "3", Distance: 3}, {Id: "5", Distance: 5}, {Id: "2", Distance: 2}, {Id: "7", Distance: 7}, {Id: "4", Distance: 4}},
			input2:      2,
			expMinArray: []Coordinate{{Id: "2", Distance: 2}, {Id: "3", Distance: 3}},
			expMaxArray: []Coordinate{{Id: "7", Distance: 7}, {Id: "5", Distance: 5}},
		},
	}
	for i := 0; i < 50; i++ {
		for _, test := range tests {
			min, max := getMinMax(test.input1, test.input2)
			if len(min) != test.input2 {
				t.Fatalf("Test %s, failed as the expected length of min array was %d, but got %d", test.name, test.input2, len(min))
			}
			if len(max) != test.input2 {
				t.Fatalf("Test %s, failed as the expected length of max array was %d, but got %d", test.name, test.input2, len(max))
			}
			if !compareArray(min, test.expMinArray) {
				t.Fatalf("Test %s, failed as the expected min output was %v, but got %v", test.name, test.expMinArray, min)
			}
			if !compareArray(max, test.expMaxArray) {
				t.Fatalf("Test %s, failed as the expected max output was %v, but got %v", test.name, test.expMaxArray, max)
			}
		}
	}
}
func compareArray(a, b []Coordinate) bool {
	if len(a) != len(b) {
		return false
	}
	for i, aValue := range a {
		if !reflect.DeepEqual(aValue, b[i]) {
			return false
		}
	}
	return true
}
