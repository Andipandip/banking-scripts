package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Must provide 4 arguments <ANZ|KB> <expenses.csv> <categories.csv> <output_file.csv>")
		return
	}

	mode := os.Args[1]
	lines := readExpenses(os.Args[2])
	categories := setupCategories(os.Args[3])

	output, err := os.OpenFile(os.Args[4], os.O_RDWR|os.O_APPEND, 0)
	checkError(err)

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

		category := categorize(description, categories)
		if category == "Ignore" {
			continue
		}
		formattedLine := fmt.Sprintf("%s, %.2f, %s, %s\n", date, amount, category, description)

		_, err = output.WriteString(formattedLine)
		checkError(err)
	}
}

func readExpenses(fileName string) [][]string {
	f, err := os.Open(fileName)
	defer f.Close()
	checkError(err)

	lines, err := csv.NewReader(f).ReadAll()
	checkError(err)
	return lines
}

func setupCategories(categoryFileName string) map[string][]string {
	f, err := os.Open(categoryFileName)
	defer f.Close()
	checkError(err)

	lines, err := csv.NewReader(f).ReadAll()
	checkError(err)

	categories := make(map[string][]string)
	for i := range lines {
		l := lines[i]
		categories[l[0]] = l[1:]
	}
	return categories
}

func categorize(desc string, categories map[string][]string) string {
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

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
