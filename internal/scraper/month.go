package scraper

import (
	"fmt"
	"strings"
)

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

var indoToEnglishMonths = map[string]string{
	"januari":   "January",
	"februari":  "February",
	"maret":     "March",
	"april":     "April",
	"mei":       "May",
	"juni":      "June",
	"juli":      "July",
	"agustus":   "August",
	"september": "September",
	"oktober":   "October",
	"november":  "November",
	"desember":  "December",
}

func TranslateIndoMonth(dateStr string) string {
	for indo, eng := range indoToEnglishMonths {
		if strings.Contains(strings.ToLower(dateStr), indo) {
			return strings.Replace(strings.ToLower(dateStr), indo, eng, 1)
		}
	}
	return dateStr
}
