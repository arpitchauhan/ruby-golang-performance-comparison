package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	reader, err := os.Open(os.Args[1])
	defer reader.Close()
	check(err)

	csvReader := csv.NewReader(reader)
	check(err)

	invalidRows := 0

	records, err := csvReader.ReadAll()
	check(err)

	for _, record := range records {
		if record[0] == "" ||
			record[1] == "" ||
			record[2] == "" ||
			record[3] == "" {
			invalidRows++
		} else {
			i, err := strconv.Atoi(record[3])
			if err != nil || i < 1 || i > 50 {
				invalidRows++
			}
		}
	}

	fmt.Printf("Number of invalid rows = %v\n", invalidRows)
}
