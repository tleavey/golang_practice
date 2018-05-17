// Homework Assignment 1: Read CSV File and Manipulate
// Tim Leavey

// Install Go
// To run from command line:
// go build hw1.go
// ./hw1.go filename

// Design note: I inserted all the rows (as array type) into one large array.
// So I work with and manipulate an array of (an array of strings).
// This program works for either file as input.

// Side note: This is my first Go program ever aside from 'Hello World'. Yay!
// That being said, I was learning all the syntax for the first time, and I
// realize some of this is inelegant code. Ideally, I'll come back and make this
// object oriented in the Go fashion.
// As I got a little more comfortable with Go, I started writing more
// functions later in the code.
// Did a huge refactor at the end to get more code into functions.

package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"sort"
	"strings"
	"strconv"
	"log"
)

// This custom sort uses the Go-required format of Len, Swap, and Less 
// to sort my rows based on the Product column in my csv file.
// I compare the length of the string and compare the string itself.
// I'm comparing "Product1" vs "Product2" for example.
// Then the full row gets swapped depending on the outcome.
type ByProduct [][]string

func (s ByProduct) Len() int {
    return len(s)
}
func (s ByProduct) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s ByProduct) Less(i, j int) bool {
    return s[i][1] < s[j][1]
}

// This function calculates the average from the transaction data.
// It takes in a slice of ints, finds the average, and returns it as an int.
func getAverageFromSlice(numbersToEvaluate []int) int {
	sum := 0
	for i := range numbersToEvaluate {
		sum += numbersToEvaluate[i]
	}
	average := sum / len(numbersToEvaluate)
	return average
}

// This function converts a slice (dynamic array) of strings into a slice of integers.
// It uses the go library "strconv" to change one value at a time in a loop and
// appends the values into a new slice of ints which it returns.
// It also ignores any value that is an empty string, meaning no price was set.
func sliceAtoi(transactionData []string) ([]int, error) {
    transactionDataInts := make([]int, 0, len(transactionData))
    for _, oneValue := range transactionData {
    	if oneValue != "" {
    		theAmount, err := strconv.Atoi(oneValue)
        	if err != nil {
            	return transactionDataInts, err
        	}
        	transactionDataInts = append(transactionDataInts, theAmount)
    	}
    }
    return transactionDataInts, nil
}

// This function prints my data set to the screen.
func printToScreen(allData [][]string) {
	for x := range allData {
		fmt.Println(allData[x])
	}
}

// This error message is from
// https://golangcode.com/write-data-to-a-csv-file/
// If err is not nil then an error message is logged.
func checkError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}

// This function counts the number of Amandas in either file.
func countAmandas(allRowsOfData [][]string) int {
	// The if statements check to see if it's file01 or file02 it's dealing with.
	// The for loops look through all the customers' names and puts them to lowercase.
	// It then stores those lowercased names in a slice (dynamic array).
	var namesOfCustomers []string
	if allRowsOfData[0][7] == "Name" {
		for z := range allRowsOfData {
			namesOfCustomers = append(namesOfCustomers, strings.ToLower(allRowsOfData[z][7]))
		}
	}
	if allRowsOfData[0][8] == "Name" {
		for z := range allRowsOfData {
			namesOfCustomers = append(namesOfCustomers, strings.ToLower(allRowsOfData[z][8]))
		}
	}
	// Using the slice of lowercase names, this for loop checks to see
	// which names contain "amanda" and adds it to the total number.
	numberOfAnyAmanda := 0
	for y := range namesOfCustomers {
		if strings.Contains(namesOfCustomers[y], "amanda") {
			numberOfAnyAmanda += 1
		}
	}
	return numberOfAnyAmanda
}

func replaceUnitedStatesWithUSA(unsortedData [][]string) [][]string {
	if unsortedData[0][6] == "Country" {
    	// The for loop replaces each "United States" with "USA"
    	for key := range unsortedData {
    		if unsortedData[key][6] == "United States" {
    			unsortedData[key][6] = strings.Replace(unsortedData[key][6], unsortedData[key][6],"USA", -1)
    		}
    	}
    }
    if unsortedData[0][7] == "Country" {
    	// The for loop replaces each "United States" with "USA"
    	for key := range unsortedData {
    		if unsortedData[key][7] == "United States" {
    			unsortedData[key][7] = strings.Replace(unsortedData[key][7], unsortedData[key][7],"USA", -1)
    		}
    	}
    }
    return unsortedData
}

// Main Program
func main() {
	// I tried making a function to pull in the csv file but kept getting blank arrays.
	// So unfortunately it is hardcoded here for now in main.
	// Get the name of the csv file from the params and store it in a variable
	fileName := os.Args[1]
	file, _ := os.Open(fileName)
	defer file.Close()

	// Create a csv file so that Go recognizes it as such using the csv library.
	// The bufio library provides buffering for I/O.
	theCSVfile := csv.NewReader(bufio.NewReader(file))

	// Create an array that will hold an array of strings (the data from our csv file)
	allRowsOfData := make([]([]string), 41)
	unsortedData := make([]([]string), 41)

	// The for loop iteratively appends each row of the csv file into an array that holds all of them.
	for i := range allRowsOfData {
		oneRowOfData, err := theCSVfile.Read()
		// This if statement ends the loop if the arrays of strings end.
		if err == io.EOF {
			break
		}
		allRowsOfData[i] = append(oneRowOfData)
		unsortedData[i] = append(oneRowOfData)
	}

	// I sort all the rows of data using my custom sort type.
	// In this case, it's sorted by 'Product'.
	sort.Sort(ByProduct(allRowsOfData))
	
	// Print all the data to the screen.
	printToScreen(allRowsOfData)

	// Count how many Amandas in either CSV file and store it in a variable.
	numberOfAnyAmanda := countAmandas(allRowsOfData)

	// Print out the number of Amandas from either csv file.
	fmt.Println("The number of customers with the name Amanda:", numberOfAnyAmanda)

	// The next section works toward displaying the average transaction amount in the db.
	// This for loop creates a slice (dynamic array) containing all the 'Price' amounts.
	// However, it contains a non-price in the first element, the label 'Price'.
	var transactionAmounts []string
	for rangeOfAllRows := range allRowsOfData {
		transactionAmounts = append(transactionAmounts, allRowsOfData[rangeOfAllRows][2])
	}
	// I remove the first element which is the label, leaving only the prices.
	transactionAmounts = append(transactionAmounts[:0], transactionAmounts[1:]...)
	// Because transactionAmounts is a slice of strings, I have to make a new slice
	// And fill it with ints instead. So I convert the strings to ints.
	var usableTransactionAmounts = []int{}
	usableTransactionAmounts, err := sliceAtoi(transactionAmounts)
    if err != nil {
        fmt.Println(err)
        return
    }
    // I find the average of all the amounts using my getAverageFromSlice() function.
    newTransactionAverage := getAverageFromSlice(usableTransactionAmounts)
    fmt.Println("The average transaction amount:", newTransactionAverage)

    // This next line replaces "United States" with "USA" on the unsorted data set.
    unsortedData = replaceUnitedStatesWithUSA(unsortedData)

    // This next section creates a new csv file and writes the results to it.
    // It is based off the code from
    // https://golangcode.com/write-data-to-a-csv-file/
    newCSVfile, err := os.Create("results.csv")
    checkError("Cannot create file", err)
    defer newCSVfile.Close()

    // After the new csv file is created, Go creates an object
    // of type Writer, which writes records to a CSV encoded file.
    // Flush() writes any buffered data to the underlying io.
    writer := csv.NewWriter(newCSVfile)
    defer writer.Flush()
    // This for loop writes each row of data to results.csv.
    // It also checks for errors.
    for _, value := range unsortedData {
        err := writer.Write(value)
        checkError("Cannot write to file", err)
    }
}