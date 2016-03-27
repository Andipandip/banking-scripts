package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var categories = make(map[string][]string)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Must provide 3 arguments <ANZ|KB> <expenses.csv> <categories.csv>")
		return
	}

	mode := os.Args[1]
	file := os.Args[2]
	fmt.Println(file)
	f, err := os.Open(file)
	defer f.Close()
	if checkError(err) {
		return
	}

	lines, err := csv.NewReader(f).ReadAll()
	if checkError(err) {
		return
	}

	setupCategories()

	outputFile, err := os.OpenFile("/Users/dawnrichardson/Google Drive/results.csv", os.O_RDWR|os.O_APPEND, 0)
	if checkError(err) {
		return
	}

	length := len(lines) - 1
	for i := range lines[1:] { //ignore header
		l := lines[length-i]
		var a = ""
		description := ""
		date := ""
		if mode == "ANZ" {
			a = l[5]
			description = fmt.Sprintf("%s %s %s %s", l[1], l[2], l[3], l[4])
			date = l[6]
		} else { //Kiwibank
			a = l[3]
			description = l[1]
			date = l[0]
		}
		amount, err := strconv.ParseFloat(a, 64)
		if amount > 0 {
			continue
		}
		amount = amount * -1.0

		category := categorize(description)
		if category == "Ignore" {
			continue
		}
		formattedLine := fmt.Sprintf("%s, %.2f, %s, %s\n", date, amount, category, description)

		_, err = outputFile.WriteString(formattedLine)
		if checkError(err) {
			return
		}
	}
}

func setupCategories() {
	file := os.Args[3]
	f, err := os.Open(file)
	defer f.Close()
	if checkError(err) {
		return
	}

	lines, err := csv.NewReader(f).ReadAll()
	if checkError(err) {
		return
	}
	for i := range lines {
		l := lines[i]
		categories[l[0]] = l[1:]
	}
}

func categorize(desc string) string {
	for key, values := range categories {
		if check(values, desc) {
			return key
		}
	}
	return "Miscellaneous"
}

func check(values []string, desc string) bool {
	for _, value := range values {
		if value == "" {
			continue
		}
		d := strings.ToLower(desc)
		v := strings.ToLower(value)
		if strings.Contains(d, v) {
			return true
		}
	}
	return false
}

func checkError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
