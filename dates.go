package nubarium

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type DateParser struct {
	expiryReferenceDate *time.Time
}

type Option func(*DateParser)

var (
	ErrDateEmpty = errors.New("date is empty")
)

func WithExpiryReferenceDate(date time.Time) Option {
	return func(p *DateParser) {
		p.expiryReferenceDate = &date
	}
}

func NewDateParser(options ...Option) (p *DateParser) {
	p = &DateParser{}
	for _, option := range options {
		option(p)
	}
	return p
}

func (p *DateParser) Parse(dateStr string) (t time.Time, err error) {
	if dateStr == "" || dateStr == "//" {
		return time.Time{}, ErrDateEmpty
	}

	parts := strings.Split(dateStr, "/")
	dayStr, monthStr, yearStr := parts[0], parts[1], parts[2]

	dayStr = removeNonDigits(dayStr)
	day, _ := strconv.Atoi(dayStr)

	// if month does not contain numbers, parse it as a month name
	monthNumStr := removeNonDigits(monthStr)
	month, err := strconv.Atoi(monthNumStr)
	if err != nil {
		m, err := time.Parse("Jan", monthStr)
		if err == nil {
			month = int(m.Month())
		} else {
			month = int(spanishMonths[strings.ToLower(strings.TrimSpace(monthStr))])
		}
	}

	yearStr = removeNonDigits(yearStr)
	year, _ := strconv.Atoi(yearStr)
	if year < 100 {
		year += 2000 // Handle 2-digit years as 2000s
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

var spanishMonths = map[string]time.Month{
	"ene":        time.January,
	"feb":        time.February,
	"mar":        time.March,
	"abr":        time.April,
	"may":        time.May,
	"jun":        time.June,
	"jul":        time.July,
	"ago":        time.August,
	"sep":        time.September,
	"oct":        time.October,
	"nov":        time.November,
	"dic":        time.December,
	"enero":      time.January,
	"febrero":    time.February,
	"marzo":      time.March,
	"abril":      time.April,
	"mayo":       time.May,
	"junio":      time.June,
	"julio":      time.July,
	"agosto":     time.August,
	"septiembre": time.September,
	"octubre":    time.October,
	"noviembre":  time.November,
	"diciembre":  time.December,
}

func removeNonDigits(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}
