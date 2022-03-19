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

func (g *GetDateSelectRangeVer2UseCase) Handle(memoName string, memoDate time.Time, insertMode InsertMode) ([]time.Time, error) {
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

	// データ取得
	dailyDataList, err := g.dailyDataRepository.Get()
	if err != nil {
		return nil, err
	}
	switch insertMode {
	case INSERT_NEWER_MODE:
		// dailyDataListから日付範囲を取得する
		return g.getDateList(dailyDataList, memoDate, memoName)
	case INSERT_OLDER_MODE:
		isExist, err := existUnderMemo(dailyDataList, memoDate, memoName)
		if err != nil {
			return nil, err
		}
		if isExist {
			// 指定dateのみを返す
			retDates := []time.Time{
				memoDate,
			}
			return retDates, nil
		} else {
			// 次のdateまでの範囲を返す
			fromDate := memoDate
			toDate, err := getToDateForOlderMode(dailyDataList, memoDate)
			if err != nil {
				return nil, err
			}
			// 1日減らしながら1つまえの日付までloop
			retDates := []time.Time{}
			for date := fromDate; date.Equal(toDate) || date.After(toDate); date = date.AddDate(0, 0, -1) {
				retDates = append(retDates, date)
				if len(retDates) >= maxCount {
					break
				}
			}
			return retDates, nil
		}
	default:
		return nil, errors.Errorf("想定外エラー insertMode: %v", insertMode)
	}
}

func (g *GetDateSelectRangeVer2UseCase) getDateList(dailyDataList []*entities.DailyDataEntity, memoDate time.Time, memoName string) ([]time.Time, error) {
	isExist, err := existUpperMemo(dailyDataList, memoDate, memoName)
	if err != nil {
		return nil, err
	}
	if isExist {
		// 指定dateのみを返す
		retDates := []time.Time{
			memoDate,
		}
		return retDates, nil
	} else {
		// 次のdateまでの範囲を返す
		// 一つ前の日付を取得
		fromDate := memoDate
		toDate, err := getToDate(dailyDataList, memoDate, g.now)
		if err != nil {
			return nil, err
		}
		// 1日addながら1つまえの日付までloop
		retDates := []time.Time{}
		for date := fromDate; date.Equal(toDate) || date.Before(toDate); date = date.AddDate(0, 0, 1) {
			retDates = append(retDates, date)
			if len(retDates) >= maxCount {
				break
			}
		}
		return retDates, nil
	}
}

func getToDate(dailyList []*entities.DailyDataEntity, fromDate, limitDate time.Time) (time.Time, error) {
	for i := len(dailyList); i >= 0; i-- {
		if dailyList[i-1].Date().Equal(fromDate) {
			if i-1 == 0 {
				// 先頭だった場合はmax値を返す
				maxDate := fromDate.AddDate(0, 0, maxCount-1)
				if maxDate.Before(limitDate) {
					return maxDate, nil
				} else {
					return limitDate, nil
				}
			} else {
				// 一つ前の日付を返す
				return dailyList[i-2].Date(), nil
			}
		}
		continue
	}
	return time.Time{}, errors.Errorf("想定外エラー dailyList: %v, fromDate: %v", dailyList, fromDate)
}

func getToDateForOlderMode(dailyList []*entities.DailyDataEntity, date time.Time) (time.Time, error) {
	for index := range dailyList {
		if dailyList[index].Date().Equal(date) {
			if index == len(dailyList)-1 {
				// max
				maxDate := date.AddDate(0, 0, -maxCount-1)
				return maxDate, nil
			} else {
				return dailyList[index+1].Date(), nil
			}
		}
		continue
	}
	return time.Time{}, errors.Errorf("想定外エラー dailyList: %v, date: %v", dailyList, date)
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
