package usecases

import (
	"time"
)

type GetDateSelectRangeUseCase struct {
}

func NewGetDateSelectRangeUseCase() *GetDateSelectRangeUseCase {
	return &GetDateSelectRangeUseCase{}
}

func (g *GetDateSelectRangeUseCase) Handle(baseTime time.Time, beforeRangeDate, afterRangeDate int) []time.Time {
	baseTime = baseTime.In(time.Local)
	from := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day()-beforeRangeDate, 0, 0, 0, 0, time.Local)
	to := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day()+afterRangeDate, 0, 0, 0, 0, time.Local)

	retDates := []time.Time{}
	appendDate := from
	for appendDate.Before(to) || appendDate.Equal(to) {
		retDates = append(retDates, appendDate)
		appendDate = appendDate.AddDate(0, 0, 1)
	}
	return retDates
}
