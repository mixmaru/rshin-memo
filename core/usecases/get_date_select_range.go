package usecases

import (
	"github.com/pkg/errors"
	"time"
)

type GetDateSelectRangeUseCase struct {
	now time.Time
}

func NewGetDateSelectRangeUseCase(now time.Time) *GetDateSelectRangeUseCase {
	return &GetDateSelectRangeUseCase{
		now: now,
	}
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
	case INSERT_UNDER_MODE:
		from = currentCursorDate
		to = overCursorDate
	case INSERT_OVER_MODE:
		from = underCursorDate
		to = currentCursorDate
	default:
		return nil, errors.Errorf("想定外値 overCursorDate: %+v, currentCursorDate: %+v, underCursorDate: %+v, insertMode: %+v", overCursorDate, currentCursorDate, underCursorDate, insertMode)
	}

	if from.IsZero() && to.IsZero() {
		// 両方zero値のときはnowを真ん中に前後でmaxCount分
		from = g.now.AddDate(0, 0, -(maxCount / 2))
		to = time.Unix(1<<63-62135596801, 999999999) // compareを使える範囲での最大値。https://stackoverflow.com/questions/25065055/what-is-the-maximum-time-time-in-go/32620397#32620397
	} else {
		// fromがゼロ値(無指定)の場合の補正
		if from.IsZero() {
			from = to.AddDate(0, 0, -(maxCount - 1))
		}
		// toがゼロ値(無指定)の場合の補正
		if to.IsZero() {
			to = time.Unix(1<<63-62135596801, 999999999) // compareを使える範囲での最大値。https://stackoverflow.com/questions/25065055/what-is-the-maximum-time-time-in-go/32620397#32620397
		}

	}

	retDates := []time.Time{}
	appendDate := from
	counter := 0
	for counter < maxCount && (appendDate.Before(to) || appendDate.Equal(to)) {
		retDates = append(retDates, appendDate)
		appendDate = appendDate.AddDate(0, 0, 1)
		counter++
	}
	return retDates, nil
}
