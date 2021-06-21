package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

const dayTemplate string = `= %s =

== Handovers ==

=== %s AM Handover ===

=== %s PM Handover ===

== %s %s ==

`

var (
	shiftMap = map[string]string{
		"primary":   "Pages",
		"weekend":   "Pages",
		"secondary": "Tickets",
	}
)

func monthToShort(m time.Month) string {
	switch m {
	case time.January:
		return "Jan"
	case time.February:
		return "Feb"
	case time.March:
		return "Mar"
	case time.April:
		return "Apr"
	case time.May:
		return "May"
	case time.June:
		return "Jun"
	case time.July:
		return "Jul"
	case time.August:
		return "Aug"
	case time.September:
		return "Sep"
	case time.October:
		return "Oct"
	case time.November:
		return "Nov"
	case time.December:
		return "Dec"
	}
	return "Unknown"
}

func printCategories(day int, month time.Month, year int, kind string) {

	fmt.Printf("[[Category:On-Call]]\n[[Category:On-Call/%d]]\n[[Category:On-Call/%d/%s]]\n",
		year, year, monthToShort(month))
	end := time.Date(year, month, day+7, 8, 0, 0, 0, time.Now().Location())
	if end.Month() > month {
		fmt.Printf("[[Category:On-Call/%d/%s]]\n", end.Year(), monthToShort(end.Month()))
	}
	fmt.Printf("[[Category:On-Call/%s]]\n", kind)
}

func printDay(day, worktype string) {
	fmt.Printf(dayTemplate, day, day, day, day, worktype)
}

func main() {
	currentYear := flag.Int("cy", time.Now().Year(), "Current year")
	currentMonth := flag.Int("cm", int(time.Now().Month()), "Current oncall's month (eg 1-12)")
	currentDay := flag.Int("cd", time.Now().Day(), "Current workday's day (eg 1-31)")
	onCallType := flag.String("type", "primary", "On-Call Type (Primary/Secondary/Weekend)")
	flag.Parse()

	today := time.Date(*currentYear, time.Month(*currentMonth), *currentDay, 8, 0, 0, 0, time.Now().Location())

	// On-call day order
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	printCategories(today.Day(), today.Month(), today.Year(), *onCallType)
	fmt.Println()

	for _, day := range days {
		printDay(day, shiftMap[strings.ToLower(*onCallType)])
	}
}
