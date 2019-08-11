package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
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

func main() {
	reader, err := os.Open(os.Args[1])
	defer reader.Close()
	check(err)

	csvReader := csv.NewReader(reader)
	check(err)

	var invalidRowsCount int32
	var wg sync.WaitGroup

	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}

		wg.Add(1)

		go func(record []string) {
			if isRecordInvalid(record) {
				atomic.AddInt32(&invalidRowsCount, 1)
			}
			wg.Done()
		}(record)
	}

	wg.Wait()

	fmt.Printf("Number of invalid rows = %v\n", atomic.LoadInt32(&invalidRowsCount))
}
