package counting

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateNumbers - random number generation
func GenerateNumbers(max int) []int {
	rand.Seed(time.Now().UnixNano())
	numbers := make([]int, max)
	for i := 0; i < max; i++ {
		numbers[i] = rand.Intn(10)
	}
	return numbers
}

// Add - sequential code to add numbers
func Add(numbers []int) int64 {
	var sum int64
	for _, n := range numbers {
		sum += int64(n)
	}
	return sum
}

//TODO: complete the concurrent version of add function.

// AddConcurrent - concurrent code to add numbers
func AddConcurrent(numbers []int) int64 {
	var sum int64
	// Utilize all cores on machine
	cores := runtime.NumCPU()
	// Divide the input into parts
	var parts [][]int
	slice := int(math.Ceil(float64(len(numbers)) / float64(cores)))
	for i := 0; i < len(numbers); i += slice {
		end := i + slice
		if end > len(numbers) {
			end = len(numbers)
		}
		parts = append(parts, numbers[i:end])
	}
	// Run computation for each part in seperate goroutine.
	var wg sync.WaitGroup
	wg.Add(len(parts))
	sums := make([]int64, len(parts))
	for i, part := range parts {
		go func(part []int, i int) {
			defer wg.Done()
			sums[i] = Add(part)
		}(part, i)
	}
	wg.Wait()
	// Add part sum to cummulative sum
	for _, partsum := range sums {
		sum += partsum
	}
	return sum
}
