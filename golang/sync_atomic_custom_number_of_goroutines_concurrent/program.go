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

func allocateWorkAmongGoroutines(goroutines int, workItems int) []int {
	workItemsPerGoroutine := workItems / goroutines

	leftoverItems := workItems % goroutines

	allocatedWork := make([]int, 0, goroutines)

	for i := 0; i < goroutines; i++ {
		var work int
		if i < leftoverItems {
			work = workItemsPerGoroutine + 1
		} else {
			work = workItemsPerGoroutine
		}

		allocatedWork = append(allocatedWork, work)
	}

	return allocatedWork
}

func allocateRowsAmongGoroutines(goroutines int, rows int) [][]int {
	allocatedWork := allocateWorkAmongGoroutines(goroutines, rows)

	on := 1

	allocatedRows := make([][]int, goroutines)

	for i := range allocatedRows {
		allocatedRows[i] = make([]int, 2)
	}

	for i, workForGoroutine := range allocatedWork {
		startRow := on
		endRow := -1
		if on != -1 {
			endRow = on + workForGoroutine - 1
		}

		allocatedRows[i] = []int{startRow, endRow}

		if endRow == rows || endRow == -1 {
			on = -1
		} else {
			on = endRow + 1
		}
	}

	return allocatedRows
}

func main() {
	reader, err := os.Open(os.Args[1])
	defer reader.Close()
	check(err)

	csvReader := csv.NewReader(reader)
	check(err)

	records, err := csvReader.ReadAll()
	check(err)

	numberOfGoroutines, err := strconv.Atoi(os.Args[2])
	check(err)

	allocatedRows := allocateRowsAmongGoroutines(numberOfGoroutines, len(records))

	var invalidRowsCount int32
	var wg sync.WaitGroup

	numberOfGoroutinesNeeded := 0

	for _, rowRange := range allocatedRows {
		startRow, endRow := rowRange[0], rowRange[1]

		if startRow == -1 && endRow == -1 {
			continue
		}

		numberOfGoroutinesNeeded++
		wg.Add(1)

		go func(records [][]string, startRow int, endRow int) {
			var invalidRowsCountInternal int32

			for i := startRow - 1; i <= endRow-1; i++ {
				if isRecordInvalid(records[i]) {
					invalidRowsCountInternal++
				}
			}
			atomic.AddInt32(&invalidRowsCount, invalidRowsCountInternal)
			wg.Done()
		}(records, startRow, endRow)
	}

	wg.Wait()
	fmt.Printf("Number of invalid rows = %v\n", invalidRowsCount)
}
