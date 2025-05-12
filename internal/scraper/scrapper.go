package scraper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	db "tanggalan-api/internal/database"
)

type Holiday struct {
	Date              string `json:"date"`
	Title             string `json:"title"`
	IsNationalHoliday bool   `json:"is_national_holiday"`
}

func generateDateList(startDayStr, endDayStr, monthStr, yearStr string) []string {
	var result []string

	// Format layout input dan output
	inputLayout := "2 January 2006"
	outputLayout := "2006-01-02"

	// Gabungkan untuk parsing
	startDateStr := fmt.Sprintf("%s %s %s", startDayStr, monthStr, yearStr)
	endDateStr := fmt.Sprintf("%s %s %s", endDayStr, monthStr, yearStr)

	startDate, err1 := time.Parse(inputLayout, startDateStr)
	endDate, err2 := time.Parse(inputLayout, endDateStr)
	if err1 != nil || err2 != nil {
		return result
	}

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		result = append(result, d.Format(outputLayout))
	}

	return result
}

func expandDateRange(text, currentMonth string, currentYear int) []string {
	text = strings.ToLower(strings.TrimSpace(text))
	dateStr := TranslateIndoMonth(text)

	// "2 - 4 april 2025"
	if strings.Contains(dateStr, "-") {
		parts := strings.Split(dateStr, "-")
		if len(parts) == 2 {
			start := strings.TrimSpace(parts[0])
			rest := strings.TrimSpace(parts[1]) // example: "4 april 2025"
			endParts := strings.Fields(rest)

			if len(endParts) >= 3 {
				endDay, endMonth, endYear := endParts[0], endParts[1], endParts[2]
				return generateDateList(start, endDay, endMonth, endYear)
			} else if len(endParts) == 2 {
				endDay, endMonth := endParts[0], endParts[1]
				return generateDateList(start, endDay, endMonth, fmt.Sprint(currentYear))
			} else if len(endParts) == 1 {
				endDay := endParts[0]
				return generateDateList(start, endDay, currentMonth, fmt.Sprint(currentYear))
			}
		}
	}

	// fallback: try full single date first (e.g., "1 april 2025")
	if t, err := time.Parse("2 January 2006", dateStr); err == nil {
		return []string{t.Format("2006-01-02")}
	}

	// fallback: single day format like "6 april"
	if t, err := time.Parse("2 January", dateStr); err == nil {
		t = t.AddDate(currentYear-t.Year(), 0, 0)
		return []string{t.Format("2006-01-02")}
	}

	return nil
}

func ScrapEventByMonthAndYear(month string, year int) ([]Holiday, error) {
	month = strings.ToLower(month)
	url := fmt.Sprintf("https://tanggalan.com/%s-%d", month, year)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var events []Holiday

	doc.Find("#events > div > .event").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".namaevent").Text())
		if title == "" {
			return
		}

		dateText := strings.TrimSpace(s.Find(".tanggal").Text())
		if dateText == "" {
			s.Find("div").Each(func(i int, div *goquery.Selection) {
				if i > 0 && dateText == "" {
					dateText = strings.TrimSpace(div.Text())
				}
			})
		}

		dateList := expandDateRange(dateText, month, year)
		isNationalHoliday := s.Find(".libur").Length() > 0

		for _, date := range dateList {
			events = append(events, Holiday{
				Title:             title,
				Date:              date,
				IsNationalHoliday: isNationalHoliday,
			})
		}
	})

	return events, nil
}

func SyncEventsFromTanggalan(ctx context.Context, q *db.Queries, month string, year int) (string, error) {
	events, err := ScrapEventByMonthAndYear(month, year)
	if err != nil {
		return "failed", err
	}

	for _, e := range events {
		t, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			continue
		}
		insertErr := q.InsertEvent(ctx, db.InsertEventParams{
			Title:             e.Title,
			Date:              t.Format("2006-01-02"),
			IsNationalHoliday: e.IsNationalHoliday,
		})
		if insertErr != nil {
			log.Printf("Failed to insert event %s on %s: %v", e.Title, e.Date, err)
		}
	}

	return "success", nil
}
