package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	numberSlice []int
	mu          sync.RWMutex
)

func loadData(filePath string) error {

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// read number from from file
	var numbers []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		number, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("invalid number in file: %v", err)
		}
		numbers = append(numbers, number)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	numberSlice = numbers

	log.Printf("%d numbers loaded from file", len(numbers))
	return nil
}

// findValue returns the index of the target element.
// If the element does not exist, it find closest to the target within 10%
// conformity of the target. Otherwise, it returns -1.
func findValue(arr []int, target int) int {
	n := len(arr)

	tolerance := float64(target) * 0.1

	// edge cases
	if target <= arr[0] {
		if math.Abs(float64(arr[0]-target)) <= tolerance {
			return 0
		}
		return -1
	}
	if target >= arr[n-1] {
		if math.Abs(float64(arr[n-1]-target)) <= tolerance {
			return n - 1
		}
		return -1
	}

	// use binary search for search as data could be huge
	i, j := 0, n
	var mid int
	for i < j {
		mid = (i + j) / 2

		if arr[mid] == target {
			return mid
		}

		// If target is less than the middle element, search in the left half
		if target < arr[mid] {
			if mid > 0 && target > arr[mid-1] {
				return getClosestIndexByTolerance(mid-1, mid, arr, target, tolerance)
			}
			j = mid
		} else {
			// If target is greater than the middle element, search in the right half
			if mid < n-1 && target < arr[mid+1] {
				return getClosestIndexByTolerance(mid, mid+1, arr, target, tolerance)
			}
			i = mid + 1
		}
	}

	// Only a single element remains after search
	if math.Abs(float64(arr[mid]-target)) <= tolerance {
		return mid
	}

	// Return -1 if no element is within the 10% tolerance
	return -1
}

func getClosestIndexByTolerance(index1, index2 int, arr []int, target int, tolerance float64) int {
	diff1 := math.Abs(float64(arr[index1] - target))
	diff2 := math.Abs(float64(arr[index2] - target))

	// Both indices must satisfy the tolerance constraint
	if diff1 <= tolerance && diff2 <= tolerance {
		if diff1 <= diff2 {
			return index1
		}
		return index2
	}

	// Check if only one of them satisfies the tolerance
	if diff1 <= tolerance {
		return index1
	}
	if diff2 <= tolerance {
		return index2
	}

	return -1
}

func searchValueHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var number = r.PathValue("num")
	target, err := strconv.Atoi(number)
	if err != nil {
		http.Error(w, "Invalid number in path", http.StatusBadRequest)
		return
	}

	// Search for the number in the slice
	mu.RLock()
	defer mu.RUnlock()
	index := findValue(numberSlice, target)

	if index == -1 {
		http.Error(w, "Value not found", http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(strconv.Itoa(index)))
	}
}

func main() {
	// Load numbers from file into a slice on-startup
	if err := loadData("data/input.txt"); err != nil {
		log.Fatal("Error while reading file: ", err)
		os.Exit(1)
	}
	http.HandleFunc("/search/{num}", searchValueHandler)
	log.Printf("Server listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
