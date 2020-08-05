package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	files := os.Args[1:]
	wg.Add(len(files))

	for _, filename := range files {
		go SortFileAsync(filename)
	}

	wg.Wait()
	fmt.Println("All done!")
}

// SortFileAsync reads the supplied filename, sorts the values, and
// writes the contents to a sorted file. Sends the name of the file
// or the error to the supplied channel.
func SortFileAsync(filename string) {
	defer wg.Done()
	fmt.Println("started " + filename)

	values, err := Read(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Ints(values) // Sort the array of integers

	writefile := strings.Replace(filename, "r", "s", 1) // Create a new version of the file name

	if err = Write(writefile, values); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("completed " + filename)
}

// SortFile reads the supplied filename, sorts the values, and
// writes the contents to a sorted file.
func SortFile(filename string) error {
	fmt.Println("Started " + filename)
	values, err := Read(filename)
	if err != nil {
		return err
	}

	sort.Ints(values) // Sort the array of integers

	writefile := strings.Replace(filename, "r", "s", 1) // Create a new version of the file name

	return Write(writefile, values)
}

// Read a file of integers and returns the values in an array
func Read(filename string) ([]int, error) {
	// Open the supplied file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close() // Defer closing the file until the end of the function

	values := []int{}              // Create an empty array to fill with the values
	scanner := bufio.NewScanner(f) // Create a scanner that will read the file line by line

	// Iterate through each line
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text()) // Convert the text from the file to integers
		if err != nil {
			return nil, err
		}

		values = append(values, v) // Append the value to the array of integers
	}

	return values, scanner.Err() // Return the values and the error if there was one
}

// Write values of an integer array to a file, separated by new lines
func Write(filename string, values []int) error {
	// Create a file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close() // Defer closing the file until the end of the function

	// Iterate through the array values and write them to the
	for _, v := range values {
		_, err := f.WriteString(strconv.Itoa(v) + "\n")
		if err != nil {
			return err
		}
	}

	// No errors? Return nil
	return nil
}
