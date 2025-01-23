package calendarManager

// This file holds the structs of the calendar manager

// This file holds the structures of the calendar years (I.e 2021, 2022, etc)
type YearTree struct {
	Years map[int]WeekTree
}

// This holds the week numbers of the year (I.e 01, 02, etc)
type WeekTree struct {
	Weeks map[int]MonthDayTree
}

// This holds the month and days of the week (I.e 0121, 0122, etc)
type MonthDayTree struct {
	MonthDays []string
}
