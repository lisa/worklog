package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	// One day is 24 hours
	day = 24 * time.Hour
)

// Any possible day we can generate
type Worklog struct {
	Year  int
	Month time.Month
	Day   int
}

func NewWorklogFromTimeObj(t time.Time) *Worklog {
	return &Worklog{
		Year:  t.Year(),
		Month: t.Month(),
		Day:   t.Day(),
	}
}

func (w *Worklog) String() string {
	return fmt.Sprintf("%d-%s-%d", w.Year, monthToShort(w.Month), w.Day)
}

func (w *Worklog) toTime() time.Time {
	return time.Date(w.Year, w.Month, w.Day, 8, 0, 0, 0, time.Now().Location())
}

// Add a duration to the Worklog to make a new Worklog to represent a future
// date
func (w *Worklog) Add(d time.Duration) *Worklog {
	return NewWorklogFromTimeObj(w.toTime().Add(d))
}

// Subtract a duration from the worklog to get a new Worklog to represent a past
// date
func (w *Worklog) Subtract(d time.Duration) *Worklog {
	return w.Add(-1 * d)
}

// Create a Worklog object for the next workday in the most
func (w *Worklog) NextDay(todayIsFriday *bool) *Worklog {
	if *todayIsFriday {
		// Add 3 days to today to get Monday
		return w.Add(3 * day)
	}
	// Otherwise, just add one day
	return w.Add(1 * day)
}

func (w *Worklog) PrevDay(todayIsMonday *bool) *Worklog {
	if *todayIsMonday {
		// Subtract 3 days from today to get Friday
		return w.Subtract(3 * day)
	}
	// Otherwise, just subtract a day
	return w.Subtract(1 * day)
}

// NewWorklog will make a new Worklog object for the given year, month and day.
func NewWorklog(year, month, day int) *Worklog {
	return &Worklog{
		Year:  year,
		Month: time.Month(month),
		Day:   day,
	}
}

func makeWikiString(prevDate, today, nextDate *Worklog) string {
	ret := fmt.Sprintf("{{Worklog|currentyear=%d|currentmonth=%s|currentday=%02d", today.Year, monthToShort(today.Month), today.Day)
	ret = ret + fmt.Sprintf("|nextyear=%d|nextmonth=%s|nextday=%d", nextDate.Year, monthToShort(nextDate.Month), nextDate.Day)
	ret = ret + fmt.Sprintf("|prevyear=%d|prevmonth=%s|prevday=%d}}\n", prevDate.Year, monthToShort(prevDate.Month), prevDate.Day)
	return ret
}

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
	log.Printf("Got an unexpected month, which is weird. m=%#v", m)
	return "Unknown"
}

func doItObj(currentYear, currentMonth, currentDay,
	nextYear, nextMonth, nextDay,
	prevYear, prevMonth, prevDay *int,
	friday, monday, verbose *bool) string {
	today := NewWorklog(*currentYear, *currentMonth, *currentDay)
	var prev, next *Worklog
	if (*nextYear == -1) && (*nextMonth == -1) && (*nextDay == -1) {
		// Easy case - we're not given any explicit part of the next date's components
		if *verbose {
			log.Printf("Easy case for next day's calculation")
		}
		next = today.NextDay(friday)
	} else {
		// Hard case - we were given a part of the next date's component
		tomorrow := today.Add(1 * day).toTime()
		if *verbose {
			log.Printf("Hard case for next day's (default next=%s) calculation: ny: %d; nm: %d; nd: %d", NewWorklogFromTimeObj(tomorrow), *nextYear, *nextMonth, *nextDay)
		}
		if *nextYear == -1 {
			y := tomorrow.Year()
			nextYear = &y
		}
		if *nextMonth == -1 {
			m := int(tomorrow.Month())
			nextMonth = &m
		}
		if *nextDay == -1 {
			d := tomorrow.Day()
			nextDay = &d
		}
		next = NewWorklog(*nextYear, *nextMonth, *nextDay)
	}
	if (*prevYear == -1) && (*prevMonth == -1) && (*prevDay == -1) {
		if *verbose {
			log.Printf("Easy case for previous day's calculation")
		}
		prev = today.PrevDay(monday)
	} else {
		yesterday := today.Subtract(1 * day).toTime()
		if *verbose {
			log.Printf("Hard case for previous day's (default prev=%s) calculation: py: %d; pm: %d; pd: %d", NewWorklogFromTimeObj(yesterday), *prevYear, *prevMonth, *prevDay)
		}
		if *prevYear == -1 {
			y := yesterday.Year()
			prevYear = &y
		}
		if *prevMonth == -1 {
			m := int(yesterday.Month())
			prevMonth = &m
		}
		if *prevDay == -1 {
			d := yesterday.Day()
			prevDay = &d
		}
		prev = NewWorklog(*prevYear, *prevMonth, *prevDay)
	}
	if *verbose {
		log.Printf("Yesterday is %s", prev)
		log.Printf("Today is %s", today)
		log.Printf("Tomorrow is %s", next)
	}

	return makeWikiString(prev, today, next)
}

func main() {
	verbose := flag.Bool("verbose", false, "Be verbose?")

	currentYear := flag.Int("cy", time.Now().Year(), "Current year")
	currentMonth := flag.Int("cm", int(time.Now().Month()), "Current workday's month (eg 1-12)")
	currentDay := flag.Int("cd", time.Now().Day(), "Current workday's day (eg 1-31)")

	nextYear := flag.Int("ny", -1, "Next workday's year")
	nextMonth := flag.Int("nm", -1, "Next workday's month (1-12)")
	nextDay := flag.Int("nd", -1, "Next workday's day (eg 1-31)")

	prevYear := flag.Int("py", -1, "Previous workday's year")
	prevMonth := flag.Int("pm", -1, "Previous workday's month (eg 1-12)")
	prevDay := flag.Int("pd", -1, "Previous workday's day (eg 1-31)")

	friday := flag.Bool("friday", false, "Today is Friday, so act as if -nd currentDay+3 was passed to save mental math")
	monday := flag.Bool("monday", false, "Today is Monday, so act as if -pd currentDay-3 was passed to save mental math")

	flag.Parse()
	log.SetOutput(os.Stderr)
	if *monday && *friday {
		fmt.Println("Providing -monday and -friday at the same time is weird. It can't be both.")
		os.Exit(1)
	}
	if *friday && (*nextDay != -1) {
		log.Printf("Providing -nd with -friday is weird")
	}
	if *monday && (*prevDay != -1) {
		log.Printf("Providing -pd with -monday is weird")
	}
	o := doItObj(currentYear, currentMonth, currentDay,
		nextYear, nextMonth, nextDay,
		prevYear, prevMonth, prevDay,
		friday, monday, verbose)
	fmt.Printf("%s\n", o)
}
