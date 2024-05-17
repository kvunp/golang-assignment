package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const numWorkers = 8
const chunkSize = 10 * 1024 * 1024

type PlaceData struct {
	sum   float64
	count int
	min   float64
	max   float64
}

func worker(id int, lineChan <-chan string, resultChan chan<- map[string]*PlaceData, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Started worker %d\n", id)
	partialResults := make(map[string]*PlaceData)
	for line := range lineChan {
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			continue
		}
		place := parts[0]
		temp, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		if _, exists := partialResults[place]; !exists {
			partialResults[place] = &PlaceData{min: temp, max: temp}
		}
		partialResults[place].sum += temp
		partialResults[place].count++
		if temp < partialResults[place].min {
			partialResults[place].min = temp
		}
		if temp > partialResults[place].max {
			partialResults[place].max = temp
		}
	}
	resultChan <- partialResults
}

func main() {
	var elapsedTime time.Duration
	start := time.Now()
	filename := "measurements.txt"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lineChan := make(chan string, 10000)
	resultChan := make(chan map[string]*PlaceData, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, lineChan, resultChan, &wg)
	}

	go func() {
		reader := bufio.NewReader(file)
		buffer := make([]byte, chunkSize)
		var remaining string

		for {
			n, err := reader.Read(buffer)
			if n > 0 {
				data := remaining + string(buffer[:n])
				lines := strings.Split(data, "\n")
				for i, line := range lines {
					if i == len(lines)-1 {
						remaining = line
					} else {
						lineChan <- line
					}
				}
			}

			if err != nil {
				break
			}
		}

		if remaining != "" {
			lineChan <- remaining
		}
		close(lineChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	finalResults := make(map[string]*PlaceData)
	for partialResults := range resultChan {
		for place, data := range partialResults {
			if _, exists := finalResults[place]; !exists {
				finalResults[place] = &PlaceData{min: data.min, max: data.max}
			}
			finalResults[place].sum += data.sum
			finalResults[place].count += data.count
			if data.min < finalResults[place].min {
				finalResults[place].min = data.min
			}
			if data.max > finalResults[place].max {
				finalResults[place].max = data.max
			}
		}
	}
	elapsedTime = time.Since(start)

	for place, data := range finalResults {
		average := data.sum / float64(data.count)
		fmt.Printf("Place: %s, Min Temperature: %.2f, Max Temperature: %.2f, Average Temperature: %.2f\n",
			place, data.min, data.max, average)
	}

	fmt.Println("Time taken: ", elapsedTime.Seconds(), " seconds")
}
