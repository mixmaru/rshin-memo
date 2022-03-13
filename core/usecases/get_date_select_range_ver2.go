package usecases

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/pkg/errors"
	"time"
)

type GetDateSelectRangeVer2UseCase struct {
	now                 time.Time
	dailyDataRepository repositories.DailyDataRepositoryInterface
}

func NewGetDateSelectRangeVer2UseCase(now time.Time, repositoryInterface repositories.DailyDataRepositoryInterface) *GetDateSelectRangeVer2UseCase {
	return &GetDateSelectRangeVer2UseCase{
		now:                 now,
		dailyDataRepository: repositoryInterface,
	}
}

//type InsertMode int
//
//const (
//	INSERT_NEWER_MODE InsertMode = iota
//	INSERT_OLDER_MODE
//)
//const maxCount = 30

//func (g *GetDateSelectRangeVer2UseCase) Handle(overCursorDate, currentCursorDate, underCursorDate time.Time, insertMode InsertMode) ([]time.Time, error) {
func (g *GetDateSelectRangeVer2UseCase) Handle(memoName string, date time.Time, insertMode InsertMode) ([]time.Time, error) {
	// dateのmemoNameのmemo一覧を取得
	// dateの前後1日のmemo一覧も取得
	// NewerModeの場合
	// 		指定のmemoの上に同日でmemoがある場合
	//			指定dateのみを返す
	//		そうでない場合
	//			次のdateまでの範囲を返す
	// OlderModeの場合
	//		指定のmemoの↓に同日でmemoがある場合
	//			指定dateのみを返す
	// 		そうでない場合
	//			次のdateまでの範囲を返す
	dailyDataList, err := g.dailyDataRepository.Get()
	if err != nil {
		return nil, err
	}
	switch insertMode {
	case INSERT_NEWER_MODE:
		isExist, err := existUpperMemo(dailyDataList, date, memoName)
		if err != nil {
			return nil, err
		}
		if isExist {
			// 指定dateのみを返す
			retDates := []time.Time{
				date,
			}
			return retDates, nil
		} else {
			// 次のdateまでの範囲を返す
			// 一つ前の日付を取得
			fromDate := date
			toDate, err := getToDate(dailyDataList, date)
			if err != nil {
				return nil, err
			}
			// 1日addながら1つまえの日付までloop
			retDates := []time.Time{}
			for date := fromDate; date.Equal(toDate) || date.Before(toDate); date = date.AddDate(0, 0, 1) {
				retDates = append(retDates, date)
			}
			return retDates, nil
		}
	case INSERT_OLDER_MODE:
		isExist, err := existUnderMemo(dailyDataList, date, memoName)
		if err != nil {
			return nil, err
		}
		if isExist {
			// 指定dateのみを返す
			retDates := []time.Time{
				date,
			}
			return retDates, nil
		} else {
			// 次のdateまでの範囲を返す
			return nil, errors.Errorf("想定外エラー insertMode: %v", insertMode)
		}
	default:
		return nil, errors.Errorf("想定外エラー insertMode: %v", insertMode)
	}

	//retDates := []time.Time{
	//	time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
	//	time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
	//	time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
	//	time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
	//	time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
	//	time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
	//}
	//return retDates, nil
	//var from, to time.Time
	//switch insertMode {
	//case INSERT_NEWER_MODE:
	//	from = currentCursorDate
	//	to = overCursorDate
	//case INSERT_OLDER_MODE:
	//	from = g.adjustFromDate(currentCursorDate, underCursorDate)
	//	to = currentCursorDate
	//default:
	//	return nil, errors.Errorf("想定外値 overCursorDate: %+v, currentCursorDate: %+v, underCursorDate: %+v, insertMode: %+v", overCursorDate, currentCursorDate, underCursorDate, insertMode)
	//}
	//
	//if from.IsZero() && to.IsZero() {
	//	// 両方zero値のときはnowを真ん中に前後でmaxCount分
	//	from = g.now.AddDate(0, 0, -(maxCount / 2))
	//	to = time.Unix(1<<63-62135596801, 999999999) // compareを使える範囲での最大値。https://stackoverflow.com/questions/25065055/what-is-the-maximum-time-time-in-go/32620397#32620397
	//} else {
	//	// fromがゼロ値(無指定)の場合の補正
	//	if from.IsZero() {
	//		from = to.AddDate(0, 0, -(maxCount - 1))
	//	}
	//	// toがゼロ値(無指定)の場合の補正
	//	if to.IsZero() {
	//		to = time.Unix(1<<63-62135596801, 999999999) // compareを使える範囲での最大値。https://stackoverflow.com/questions/25065055/what-is-the-maximum-time-time-in-go/32620397#32620397
	//	}
	//
	//}
	//
	//retDates := []time.Time{}
	//appendDate := from
	//counter := 0
	//for counter < maxCount && (appendDate.Before(to) || appendDate.Equal(to)) {
	//	retDates = append(retDates, appendDate)
	//	appendDate = appendDate.AddDate(0, 0, 1)
	//	counter++
	//}
	//
	//if insertMode == INSERT_OLDER_MODE {
	//	for i := 0; i < len(retDates)/2; i++ {
	//		retDates[i], retDates[len(retDates)-i-1] = retDates[len(retDates)-i-1], retDates[i]
	//	}
	//}
	//return retDates, nil
}

func getToDate(dailyList []*entities.DailyDataEntity, date time.Time) (time.Time, error) {
	for i := len(dailyList); i >= 0; i-- {
		if dailyList[i-1].Date().Equal(date) {
			return dailyList[i-2].Date(), nil
		}
		continue
	}
	return time.Time{}, errors.Errorf("想定外エラー")
}

func existUpperMemo(dailyDataList []*entities.DailyDataEntity, date time.Time, memoName string) (bool, error) {
	for _, dailyData := range dailyDataList {
		if dailyData.Date() == date {
			for index, memo := range dailyData.NoteNames() {
				if memo == memoName {
					if index != 0 {
						return true, nil
					} else {
						return false, nil
					}
				}
				continue
			}
		}
		continue
	}
	return false, errors.Errorf("想定外エラー dailyDataList: %v, date: %v, memoName: %v", dailyDataList, date, memoName)
}

func existUnderMemo(dailyDataList []*entities.DailyDataEntity, date time.Time, memoName string) (bool, error) {
	for _, dailyData := range dailyDataList {
		if dailyData.Date() == date {
			for i := len(dailyData.NoteNames()); i >= 0; i-- {
				if dailyData.NoteNames()[i-1] == memoName {
					if i != len(dailyData.NoteNames()) {
						return true, nil
					} else {
						return false, nil
					}
				}
				continue
			}
		}
		continue
	}
	return false, errors.Errorf("想定外エラー dailyDataList: %v, date: %v, memoName: %v", dailyDataList, date, memoName)
}

func (g *GetDateSelectRangeVer2UseCase) adjustFromDate(currentCursorDate time.Time, underCursorDate time.Time) time.Time {
	duration := currentCursorDate.Sub(underCursorDate)
	if duration >= maxCount*24*60*60*1000*1000*1000 {
		return time.Date(currentCursorDate.Year(), currentCursorDate.Month(), currentCursorDate.Day()-(maxCount-1), 0, 0, 0, 0, time.Local)
	} else {
		return underCursorDate
	}
}
