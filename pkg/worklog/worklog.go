package main

import (
	"flag"
	"fmt"
	"time"
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

func main() {
	currentYear := flag.Int("cy", time.Now().Year(), "Current year")
	currentMonth := flag.Int("cm", int(time.Now().Month()), "Current workday's month (eg 1-12)")
	currentDay := flag.Int("cd", time.Now().Day(), "Current workday's day (eg 1-31)")

	nextYear := flag.Int("ny", -1, "Next workday's year")
	nextMonth := flag.Int("nm", -1, "Next workday's month (1-12)")
	nextDay := flag.Int("nd", -1, "Next workday's day (eg 1-31)")

	prevYear := flag.Int("py", -1, "Previous workday's year")
	prevMonth := flag.Int("pm", -1, "Previous workday's month (eg 1-12)")
	prevDay := flag.Int("pd", -1, "Previous workday's day (eg 1-31)")

	flag.Parse()

	// Today, as far as the wiki is concerned. All of these are from POV of the Wiki talking about work days
	today := time.Date(*currentYear, time.Month(*currentMonth), *currentDay, 8, 0, 0, 0, time.Now().Location())
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	var py, pd, ny, nd int
	var pm, nm time.Month

	if *nextYear == -1 {
		ny = tomorrow.Year()
	} else {
		ny = *nextYear
	}
	if *nextMonth == -1 {
		nm = tomorrow.Month()
	} else {
		nm = time.Month(*nextMonth)
	}
	if *nextDay == -1 {
		nd = tomorrow.Day()
	} else {
		nd = *nextDay
	}
	if *prevYear == -1 {
		py = yesterday.Year()
	} else {
		py = *prevYear
	}
	if *prevMonth == -1 {
		pm = yesterday.Month()
	} else {
		pm = time.Month(*prevMonth)
	}
	if *prevDay == -1 {
		pd = yesterday.Day()
	} else {
		pd = *prevDay
	}

	nextDate := time.Date(ny, nm, nd, 8, 0, 0, 0, tomorrow.Location())
	prevDate := time.Date(py, pm, pd, 8, 0, 0, 0, yesterday.Location())

	fmt.Printf("{{Worklog|currentyear=%d|currentmonth=%s|currentday=%02d", today.Year(), monthToShort(today.Month()), today.Day())
	fmt.Printf("|nextyear=%d|nextmonth=%s|nextday=%d", nextDate.Year(), monthToShort(nextDate.Month()), nextDate.Day())
	fmt.Printf("|prevyear=%d|prevmonth=%s|prevday=%d}}\n", prevDate.Year(), monthToShort(prevDate.Month()), prevDate.Day())

}
