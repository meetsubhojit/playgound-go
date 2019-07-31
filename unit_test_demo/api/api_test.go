package api

import (
	"errors"
	"github.com/TheGUNNER13/playgound-go/unit_test_demo/external"
	"strings"
	"testing"
)

func TestDoAddAndSave(t *testing.T) {

	tests := []struct {
		name string
		//input
		in string
		//output
		err    string
		output int
		//mock
		mock external.Cache
	}{
		{
			name:   "error 1",
			in:     "4+5",
			err:    "some forced error",
			output: 9,
			mock: &mock{
				myOwnFunc: func(input interface{}) error {
					return errors.New("some forced error")
				},
			},
		},
		//{
		//	name:   "sucess 2",
		//	in:     "5+5",
		//	err:    "",
		//	output: 10,
		//},
		//{
		//	name: "error 1",
		//	in:   "5",
		//	err:  "need to pass two numbers atleast",
		//},
		//{
		//	name: "error 1",
		//	in:   "+5",
		//	err:  "need to pass two numbers atleast",
		//},
		//{
		//	name: "error 1",
		//	in:   "5,5",
		//	err:  "need to pass two numbers atleast",
		//},
	}

	for _, test := range tests {

		cache = test.mock
		//NewAPI(cache)

		actualOutput, actualErr := DoAddAndSave(test.in)
		//fmt.Println(actualOutput, actualErr)
		if actualErr != nil {
			if test.err == "" {
				t.Fatal(test.name, "expected ", test.err, " got ", actualErr.Error())
			}
			if !strings.Contains(actualErr.Error(), test.err) {
				t.Fatal(test.name)
			}
		} else {
			if test.err != "" {
				t.Fatal(test.name)
			}
		}
		if actualOutput != test.output {
			t.Fatal(test.name)
		}
	}
}

type mock struct {
	myOwnFunc func(input interface{}) error
}

func (m *mock) Save(input interface{}) error {
	return m.myOwnFunc(input)
}
