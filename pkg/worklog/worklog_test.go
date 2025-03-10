package main

import (
	"fmt"
	"testing"
)

func makeIntPointer(i int) *int    { return &i }
func makeBoolPointer(b bool) *bool { return &b }

// testObject represents a set of inputs to the worklog program.
type testObject struct {
	cy *int
	cm *int
	cd *int

	ny *int
	nm *int
	nd *int

	py *int
	pm *int
	pd *int

	friday  *bool
	monday  *bool
	verbose *bool
}

func (o *testObject) next() (*int, *int, *int)  { return o.ny, o.nm, o.nd }
func (o *testObject) today() (*int, *int, *int) { return o.cy, o.cm, o.cd }
func (o *testObject) prev() (*int, *int, *int)  { return o.py, o.pm, o.pd }

func makeTestObj(cy, cm, cd, ny, nm, nd, py, pm, pd int, friday, monday, verbose bool) *testObject {
	return &testObject{
		cy: makeIntPointer(cy),
		cm: makeIntPointer(cm),
		cd: makeIntPointer(cd),

		ny: makeIntPointer(ny),
		nm: makeIntPointer(nm),
		nd: makeIntPointer(nd),

		py: makeIntPointer(py),
		pm: makeIntPointer(pm),
		pd: makeIntPointer(pd),

		friday:  makeBoolPointer(friday),
		monday:  makeBoolPointer(monday),
		verbose: makeBoolPointer(verbose),
	}
}

func make2025Feb28ToMondayObj() *testObject {
	return makeTestObj(2025, 2, 28, -1, -1, -1, -1, -1, -1, true, false, false)
}

func make2024Dec31() *testObject {
	return makeTestObj(2024, 12, 31, -1, -1, 2, -1, -1, -1, false, false, false)
}

// next day 2025-Jan-2, prev day 2024-Dec-30
func Test2024Dec31(t *testing.T) {
	to := make2024Dec31()
	cy, cm, cd := to.today()
	ny, nm, nd := to.next()
	py, pm, pd := to.prev()

	retObj := doItObj(cy, cm, cd, ny, nm, nd, py, pm, pd, to.friday, to.monday, to.verbose)
	shouldBe := "{{Worklog|currentyear=2024|currentmonth=Dec|currentday=31|nextyear=2025|nextmonth=Jan|nextday=2|prevyear=2024|prevmonth=Dec|prevday=30}}\n"
	if retObj != shouldBe {
		fmt.Printf("obj\nG: %s\nE: %s\n", retObj, shouldBe)
		t.Fail()
	}
}

func make2025Jan2Obj() *testObject {
	return makeTestObj(2025, 1, 2, -1, -1, -1, 2024, 12, 31, false, false, false)
}

func Test2025Jan2(t *testing.T) {
	to := make2025Jan2Obj()
	cy, cm, cd := to.today()
	ny, nm, nd := to.next()
	py, pm, pd := to.prev()

	retObj := doItObj(cy, cm, cd, ny, nm, nd, py, pm, pd, to.friday, to.monday, to.verbose)
	shouldBe := "{{Worklog|currentyear=2025|currentmonth=Jan|currentday=02|nextyear=2025|nextmonth=Jan|nextday=3|prevyear=2024|prevmonth=Dec|prevday=31}}\n"
	if retObj != shouldBe {
		fmt.Printf("obj\nG: %s\nE: %s\n", retObj, shouldBe)
		t.Fail()
	}
}
func Test2025Feb28ToMondayObj(t *testing.T) {
	to := make2025Feb28ToMondayObj()
	cy, cm, cd := to.today()
	ny, nm, nd := to.next()
	py, pm, pd := to.prev()

	ret := doItObj(cy, cm, cd, ny, nm, nd, py, pm, pd, to.friday, to.monday, to.verbose)
	shouldBe := "{{Worklog|currentyear=2025|currentmonth=Feb|currentday=28|nextyear=2025|nextmonth=Mar|nextday=3|prevyear=2025|prevmonth=Feb|prevday=27}}\n"
	if ret != shouldBe {
		fmt.Printf("G: %s\nE: %s\n", ret, shouldBe)
		t.Fail()
	}
}

// No idea why the only the current day is zero-padded
func TestCurrentDayShouldBeZeroPadded(t *testing.T) {
	to := makeTestObj(2025, 3, 4, // 2025-Mar-4
		-1, -1, -1, // 2025-Mar-5
		-1, -1, -1, //2025-Mar-3
		false, false, false)
	cy, cm, cd := to.today()
	ny, nm, nd := to.next()
	py, pm, pd := to.prev()
	ret := doItObj(cy, cm, cd, ny, nm, nd, py, pm, pd, to.friday, to.monday, to.verbose)
	shouldBe := "{{Worklog|currentyear=2025|currentmonth=Mar|currentday=04|nextyear=2025|nextmonth=Mar|nextday=5|prevyear=2025|prevmonth=Mar|prevday=3}}\n"
	if ret != shouldBe {
		fmt.Printf("G: %s\nE: %s\n", ret, shouldBe)
		t.Fail()
	}
}

func TestLeapYear(t *testing.T) {
	to := makeTestObj(2024, 2, 29, // 2024-Feb-29
		-1, -1, -1, // 2024-Mar-1
		-1, -1, -1, //2024-Feb-28
		false, false, false)
	cy, cm, cd := to.today()
	ny, nm, nd := to.next()
	py, pm, pd := to.prev()
	ret := doItObj(cy, cm, cd, ny, nm, nd, py, pm, pd, to.friday, to.monday, to.verbose)
	shouldBe := "{{Worklog|currentyear=2024|currentmonth=Feb|currentday=29|nextyear=2024|nextmonth=Mar|nextday=1|prevyear=2024|prevmonth=Feb|prevday=28}}\n"
	if ret != shouldBe {
		fmt.Printf("G: %s\nE: %s\n", ret, shouldBe)
		t.Fail()
	}
}

func TestLeapYearFriday(t *testing.T) {
	to := makeTestObj(2008, 2, 29, // 2008-Feb-29
		-1, -1, -1, // 2008-Mar-1
		-1, -1, -1, //2008-Feb-28
		true, false, false)
	cy, cm, cd := to.today()
	ny, nm, nd := to.next()
	py, pm, pd := to.prev()
	ret := doItObj(cy, cm, cd, ny, nm, nd, py, pm, pd, to.friday, to.monday, to.verbose)
	shouldBe := "{{Worklog|currentyear=2008|currentmonth=Feb|currentday=29|nextyear=2008|nextmonth=Mar|nextday=3|prevyear=2008|prevmonth=Feb|prevday=28}}\n"
	if ret != shouldBe {
		fmt.Printf("G: %s\nE: %s\n", ret, shouldBe)
		t.Fail()
	}
}
func TestLeapYearMonday(t *testing.T) {
	to := makeTestObj(2016, 2, 29, // 2016-Feb-29
		-1, -1, -1, // 2016-Mar-1
		-1, -1, -1, //2016-Feb-26
		false, true, false)
	cy, cm, cd := to.today()
	ny, nm, nd := to.next()
	py, pm, pd := to.prev()
	ret := doItObj(cy, cm, cd, ny, nm, nd, py, pm, pd, to.friday, to.monday, to.verbose)
	shouldBe := "{{Worklog|currentyear=2016|currentmonth=Feb|currentday=29|nextyear=2016|nextmonth=Mar|nextday=1|prevyear=2016|prevmonth=Feb|prevday=26}}\n"
	if ret != shouldBe {
		fmt.Printf("G: %s\nE: %s\n", ret, shouldBe)
		t.Fail()
	}
}

func TestShortMonthNames(t *testing.T) {
	type monthShortTests struct {
		testN    int
		shouldBe string
	}
	monthtests := []monthShortTests{
		{
			testN:    1,
			shouldBe: "Jan",
		},
		{
			testN:    2,
			shouldBe: "Feb",
		},
		{
			testN:    3,
			shouldBe: "Mar",
		},
		{
			testN:    4,
			shouldBe: "Apr",
		},
		{
			testN:    5,
			shouldBe: "May",
		},
		{
			testN:    6,
			shouldBe: "Jun",
		},
		{
			testN:    7,
			shouldBe: "Jul",
		},
		{
			testN:    8,
			shouldBe: "Aug",
		},
		{
			testN:    9,
			shouldBe: "Sep",
		},
		{
			testN:    10,
			shouldBe: "Oct",
		},
		{
			testN:    11,
			shouldBe: "Nov",
		},
		{
			testN:    12,
			shouldBe: "Dec",
		},
		{
			testN:    13,
			shouldBe: "Unknown",
		},
	}
	for _, testCase := range monthtests {
		testWorklog := NewWorklog(2025, testCase.testN, 10)
		if monthToShort(testWorklog.Month) != testCase.shouldBe {
			t.Fatalf("Failed to convert month %d to %s. Test day: %s", testCase.testN, testCase.shouldBe, testWorklog)
		}
	}
}
