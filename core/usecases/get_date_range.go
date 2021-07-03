package usecases

import (
	"time"
)

type GetDateRangeUseCase struct {
}

func NewGetDateRangeUseCase() *GetDateRangeUseCase {
	return &GetDateRangeUseCase{}
}

func (g *GetDateRangeUseCase) Handle(baseTime time.Time, beforeRange, afterRange int) (from, to time.Time) {
	baseTime = baseTime.In(time.Local)
	from = time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day()-beforeRange, 0, 0, 0, 0, time.Local)
	to = time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day()+afterRange, 0, 0, 0, 0, time.Local)
	return from, to
}

//type DateRange struct {
//	From time.Time
//	To   time.Time
//}
//
//// 日付の範囲から最大指定個数の日をスライスで返す
//func (d *DateRange) GetSomeDateInRange(num int) ([]time.Time, error) {
//	retDates := []time.Time{}
//	var startData time.Time
//	if !d.From.IsZero() {
//		startData = d.From
//	} else if !d.To.IsZero() {
//		startData = d.To.AddDate(0, 0, -(num - 1))
//	} else {
//		return nil, errors.Errorf("このDateRangeはFromもToも設定されていませんので実行できません。%+v", d)
//	}
//
//	for i := 0; i < num; i++ {
//		date := startData.AddDate(0, 0, i)
//		if !d.To.IsZero() && date.After(d.To) {
//			break
//		}
//		retDates = append(retDates, date)
//	}
//	return retDates, nil
//}
//
//func (d *DateRange) IsIn(targetDate time.Time) bool {
//	if !d.From.IsZero() {
//		if targetDate.Before(d.From) {
//			return false
//		}
//	}
//	if !d.To.IsZero() {
//		if targetDate.After(d.To) {
//			return false
//		}
//	}
//	return true
//}
//
//func (d *DateRange) SetFromByString(dateStr string) error {
//	var err error
//	d.From, err = time.Parse("2006-01-02", dateStr)
//	if err != nil {
//		return errors.WithStack(err)
//	}
//	return err
//}
//
//func (d *DateRange) SetToByString(dateStr string) error {
//	var err error
//	d.To, err = time.Parse("2006-01-02", dateStr)
//	if err != nil {
//		return errors.WithStack(err)
//	}
//	return err
//}
