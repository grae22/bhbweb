package domain

import (
	"testing"
	"time"
)

func Test_PeriodIdForDate(t *testing.T) {
	// Arrange.
	date, _ := time.Parse("2006/1/2", "2026/8/16")

	// Act.
	periodId := PeriodIdForDate(date)

	// Assert.
	if periodId != "2608" {
		t.Errorf("incorrect id \"%s\"", periodId)
	}
}

func Test_DateForPeriodId(t *testing.T) {
	// Arrange.
	periodId := "2408"

	// Act.
	date := DateForPeriodId(periodId)

	// Assert.
	if date.Year() != 2024 ||
		date.Month() != 8 ||
		date.Day() != 1 {

		t.Errorf("incorrect date \"%v\".", date)
	}
}
