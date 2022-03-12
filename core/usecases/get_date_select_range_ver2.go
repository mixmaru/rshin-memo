package usecases

import (
	"github.com/pkg/errors"
	"time"
)

type GetDateSelectRangeVer2UseCase struct {
	now time.Time
}

func NewGetDateSelectRangeVer2UseCase(now time.Time) *GetDateSelectRangeVer2UseCase {
	return &GetDateSelectRangeVer2UseCase{
		now: now,
	}
}

//type InsertMode int
//
//const (
//	INSERT_NEWER_MODE InsertMode = iota
//	INSERT_OLDER_MODE
//)
//const maxCount = 30

func (g *GetDateSelectRangeVer2UseCase) Handle(overCursorDate, currentCursorDate, underCursorDate time.Time, insertMode InsertMode) ([]time.Time, error) {
	//func (g *GetDateSelectRangeVer2UseCase) Handle(memoName string, date time.Time, insertMode InsertMode) ([]time.Time, error) {
	var from, to time.Time
	switch insertMode {
	case INSERT_NEWER_MODE:
		from = currentCursorDate
		to = overCursorDate
	case INSERT_OLDER_MODE:
		from = g.adjustFromDate(currentCursorDate, underCursorDate)
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

	if insertMode == INSERT_OLDER_MODE {
		for i := 0; i < len(retDates)/2; i++ {
			retDates[i], retDates[len(retDates)-i-1] = retDates[len(retDates)-i-1], retDates[i]
		}
	}
	return retDates, nil
}

func (g *GetDateSelectRangeVer2UseCase) adjustFromDate(currentCursorDate time.Time, underCursorDate time.Time) time.Time {
	duration := currentCursorDate.Sub(underCursorDate)
	if duration >= maxCount*24*60*60*1000*1000*1000 {
		return time.Date(currentCursorDate.Year(), currentCursorDate.Month(), currentCursorDate.Day()-(maxCount-1), 0, 0, 0, 0, time.Local)
	} else {
		return underCursorDate
	}
}
