package main

import "time"

func DiffDate(a, b time.Time) (years, months, days, hours int) {
	yearA, monthA, dayA := a.Date()
	yearB, monthB, dayB := b.Date()
	hourA, _, _ := a.Clock()
	hourB, _, _ := b.Clock()

	years = int(yearB - yearA)
	months = int(monthB - monthA)
	days = int(dayB - dayA)
	hours = int(hourB - hourA)

	if hours < 0 {
		hours += 24
		days--
	}

	if days < 0 {
		daysInMonth := time.Date(yearA, monthA, 32, 0, 0, 0, 0, time.UTC)
		days += 32 - daysInMonth.Day()
		months--
	}

	if months < 0 {
		months += 12
		years--
	}

	return
}
