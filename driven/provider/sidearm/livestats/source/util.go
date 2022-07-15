// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
