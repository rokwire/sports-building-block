package source

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// parseDate parse the string date format 9/21/2019 and returns month, day and year
func parseDate(date string) (int, int, int, error) {
	if len(date) <= 0 {
		return -1, -1, -1, errors.New("parseDate -> the xml date is empty")
	}
	dateArr := strings.Split(date, "/")
	if len(dateArr) != 3 {
		return -1, -1, -1, errors.New("parseDate -> the xml date is not split properly: " + date)
	}
	dateMonth, err := strconv.Atoi(dateArr[0])
	if err != nil {
		log.Println(err.Error())
		return -1, -1, -1, err
	}
	dateDay, err := strconv.Atoi(dateArr[1])
	if err != nil {
		log.Println(err.Error())
		return -1, -1, -1, err
	}
	dateYear, err := strconv.Atoi(dateArr[2])
	if err != nil {
		log.Println(err.Error())
		return -1, -1, -1, err
	}
	return dateMonth, dateDay, dateYear, nil
}

// isNumber checks if the given string is number
func isNumber(data string) bool {
	_, err := strconv.Atoi(data)
	if err == nil {
		return true
	}
	return false
}

// getOrdinal gives the ordinal of the input string
func getOrdinal(data string) string {
	input, err := strconv.Atoi(data)
	if err != nil {
		return ""
	}

	switch input {
	case 1:
		return "1st"
	case 2:
		return "2nd"
	case 3:
		return "3rd"
	default:
		return fmt.Sprintf("%dth", input)
	}
}
