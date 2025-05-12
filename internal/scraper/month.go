package scraper

import "fmt"

var indonesianMonths = map[int]string{
	1:  "januari",
	2:  "februari",
	3:  "maret",
	4:  "april",
	5:  "mei",
	6:  "juni",
	7:  "juli",
	8:  "agustus",
	9:  "september",
	10: "oktober",
	11: "november",
	12: "desember",
}

func GetMonthName(month int) (string, error) {
	name, ok := indonesianMonths[month]
	if !ok {
		return "", fmt.Errorf("invalid month: %d", month)
	}
	return name, nil
}
