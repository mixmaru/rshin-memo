package usecases

import (
	"github.com/pkg/errors"
	"time"
)

type GetDateSelectRangeUseCase struct {
}

func NewGetDateSelectRangeUseCase() *GetDateSelectRangeUseCase {
	return &GetDateSelectRangeUseCase{}
}

type InsertMode int

const (
	INSERT_OVER_MODE InsertMode = iota
	INSERT_UNDER_MODE
)
const maxCount = 30

func (g *GetDateSelectRangeUseCase) Handle(overCursorDate, currentCursorDate, underCursorDate time.Time, insertMode InsertMode) ([]time.Time, error) {
	var from, to time.Time
	switch insertMode {
	case INSERT_OVER_MODE:
		from = currentCursorDate
		to = overCursorDate
	case INSERT_UNDER_MODE:
		from = underCursorDate
		to = currentCursorDate
	default:
		return nil, errors.Errorf("想定外値 overCursorDate: %+v, currentCursorDate: %+v, underCursorDate: %+v, insertMode: %+v", overCursorDate, currentCursorDate, underCursorDate, insertMode)
	}

	// fromがゼロ値(無指定)の場合の補正
	if from.IsZero() {
		from = to.AddDate(0, 0, -maxCount)
	}

	retDates := []time.Time{}
	appendDate := from
	for appendDate.Before(to) || appendDate.Equal(to) {
		retDates = append(retDates, appendDate)
		appendDate = appendDate.AddDate(0, 0, 1)
	}
	return retDates, nil
}
