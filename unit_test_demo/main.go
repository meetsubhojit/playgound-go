package main

import (
	"github.com/TheGUNNER13/playgound-go/unit_test_demo/api"
	"log"
)

func main() {

	input := "1+2"
	sum, err := api.DoAddAndSave(input)

	if err != nil {
		log.Fatalf("error occurred %s, ", err.Error())
	}
	log.Printf("Sum saved as : %d", sum)
}
