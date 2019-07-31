package api

import (
	"strconv"
	"strings"
)

import (
	"errors"
	"github.com/TheGUNNER13/playgound-go/unit_test_demo/external"
	//"strconv"
	//"strings"
)

//var cache = external.Save

var cache = external.NewCacheClient()

// "4+5"
func DoAddAndSave(input string) (int, error) {

	numbers := strings.Split(input, "+")
	if len(numbers) < 2 {
		return 0, errors.New("need to pass two numbers atleast")
	}
	a, err := strconv.Atoi(numbers[0])
	if err != nil {
		return 0, errors.New("first number has wrong format " + err.Error())
	}
	b, err := strconv.Atoi(numbers[1])
	if err != nil {
		return 0, errors.New("second number has wrong format " + err.Error())
	}

	sum := a + b
	return a + b, cache.Save(sum)
}
