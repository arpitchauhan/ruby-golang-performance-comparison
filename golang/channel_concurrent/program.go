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

func isRecordInvalid(record []string) bool {
	id, firstName, lastName, age := record[0], record[1], record[2], record[3]

	if id == "" || firstName == "" || lastName == "" || age == "" {
		return true
	}

	ageInteger, err := strconv.Atoi(age)
	if err != nil {
		return true
	}

	if ageInteger < 1 || ageInteger > 50 {
		return true
	}

	return false
}

func executeRecordValidation(record []string, invalidRowsChan chan<- bool) {
	invalidRowsChan <- isRecordInvalid(record)
}

func main() {
	reader, err := os.Open(os.Args[1])
	defer reader.Close()
	check(err)

	csvReader := csv.NewReader(reader)
	check(err)

	records, err := csvReader.ReadAll()
	check(err)

	invalidRowsChan := make(chan bool, len(records))

	invalidRowsCount := 0

	for _, record := range records {
		go executeRecordValidation(record, invalidRowsChan)
	}

	for i := 1; i <= len(records); i++ {
		result := <-invalidRowsChan
		if result {
			invalidRowsCount++
		}
	}

	fmt.Printf("Number of invalid rows = %v\n", invalidRowsCount)
}
