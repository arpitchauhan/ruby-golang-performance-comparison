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

type csvRecord struct {
	id        string
	firstName string
	lastName  string
	age       string
}

func (r *csvRecord) isInvalid() bool {
	if r.id == "" || r.firstName == "" || r.lastName == "" || r.age == "" {
		return true
	}

	ageInteger, err := strconv.Atoi(r.age)
	if err != nil {
		return true
	}

	if ageInteger < 1 || ageInteger > 50 {
		return true
	}

	return false
}

func main() {
	reader, err := os.Open(os.Args[1])
	defer reader.Close()
	check(err)

	csvReader := csv.NewReader(reader)
	check(err)

	invalidRows := 0

	for {
		record, err := csvReader.Read()

		if err != nil {
			break
		}

		recordObject := csvRecord{id: record[0], firstName: record[1], lastName: record[2], age: record[3]}

		if recordObject.isInvalid() {
			invalidRows++
		}
	}

	fmt.Printf("Number of invalid rows = %v\n", invalidRows)
}
