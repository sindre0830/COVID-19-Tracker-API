package fun

import (
	"errors"
	"strconv"
	"strings"
)

// ValidateCountry checks if country is empty.
func ValidateCountry(country string) error {
	if country == "" {
		err := errors.New("country validation: doesn't exist in database")
		return err
	}
	return nil
}

// ValidateDates checks if there are any formatting mistakes in scope.
func ValidateDates(dates string) error {
	//split date by '-' for format checking
	arrDate := strings.Split(dates, "-")
	//branch if date doesn't have correct amount of elements
	if (len(arrDate) != 6) || (len(dates) != 21) {
		return errors.New("date validation: not enough elements")
	}
	//branch if start date isn't using correct format (YYYY-MM-DD)
	if (len(arrDate[0]) != 4) || (len(arrDate[1]) != 2) || (len(arrDate[2]) != 2) {
		return errors.New("date validation: start date doesn't follow required date format")
	}
	//branch if end date isn't using correct format (YYYY-MM-DD)
	if (len(arrDate[3]) != 4) || (len(arrDate[4]) != 2) || (len(arrDate[5]) != 2) {
		return errors.New("date validation: end date doesn't follow required date format")
	}
	//branch if date elements aren't integers or larger than 0. 'hehe-01-00' == false
	for _, elemDate := range arrDate {
		elemDateNum, err := strconv.Atoi(elemDate)
		if err != nil || elemDateNum < 1 {
			return errors.New("date validation: wrong type, should be numbers that are larger than 0")
		}
	}
	//branch if end date isn't larger or equal to start date
	if dates[:10] > dates[11:] {
		return errors.New("date validation: start date is larger than end date")
	}
	return nil
}
