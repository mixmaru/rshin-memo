package dto

import (
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
)

type InsertData struct {
	InsertNum       int    // dailyDataの上から何番目に挿入するか
	NoteName        string // note名
	DateStr         string
	TargetDailyData []usecases.DailyData
}

func (i *InsertData) GenerateNewDailyData() (usecases.DailyData, error) {
	//var retDailyData usecases.DailyData
	//retDailyData.Date = i.targetDailyData.Date
	//for index, noteName := range i.targetDailyData.Notes {
	//	if index == i.InsertNum() {
	//		retDailyData.Notes = append(retDailyData.Notes, i.NoteName())
	//	}
	//	retDailyData.Notes = append(retDailyData.Notes, noteName)
	//}
	//if i.insertNum == len(i.targetDailyData.Notes) {
	//	// 末への追加
	//	retDailyData.Notes = append(retDailyData.Notes, i.NoteName())
	//}
	//return retDailyData
	//return usecases.DailyData{}
	retDailyData, err := generateNewDailyData(i.TargetDailyData, i.NoteName, i.DateStr, i.InsertNum)
	if err != nil {
		return usecases.DailyData{}, err
	}
	return retDailyData, nil
}

/*
dailyListとdate文字列と新Note名と挿入希望位置を渡すとdailyDataを生成して返す
*/
func generateNewDailyData(dailyList []usecases.DailyData, newNoteName string, date string, insertLineNum int) (usecases.DailyData, error) {
	if len(dailyList) == 0 {
		// 新しいdailyDataを作成する
		retDailyData := usecases.DailyData{
			Date: date,
			Notes: []string{
				newNoteName,
			},
		}
		return retDailyData, nil
	}

	if insertLineNum == 0 {
		//新しく先頭にdailyDateをつくるのか、最初のdailyDateの先頭に挿入するのか判断しないといけない
		if dailyList[0].Date != date {
			// 新しいdailyDataを作成する
			retDailyData := usecases.DailyData{
				Date: date,
				Notes: []string{
					newNoteName,
				},
			}
			return retDailyData, nil
		}
	}

	insertLineNum++ // lenと比較しやすくするため+1する
	newNotes := []string{}
	for index, dailyData := range dailyList {
		if insertLineNum > len(dailyData.Notes) {
			// noteListの数よりinsert位置があとなので次のdailyDateを確認する
			insertLineNum -= len(dailyData.Notes)
			continue
		} else {
			if insertLineNum == 1 {
				// このdailyDataの先頭noteに追加すべきなのか、一つ前のdailyDataの末尾noteに追加すべきなのか判断しないといけない
				if index > 0 && dailyList[index-1].Date == date {
					// 一つ前の末尾に追加
					dailyList[index-1].Notes = append(dailyList[index-1].Notes, newNoteName)
					return dailyList[index-1], nil
				} else if dailyData.Date == date {
					// このdailyDataの先頭についか
					newNotes = append(newNotes, newNoteName)
					newNotes = append(newNotes, dailyData.Notes...)
					dailyData.Notes = newNotes
					return dailyData, nil
				} else {
					// 新しくdailyDataを作成する
					newDailyData := usecases.DailyData{
						Date: date,
						Notes: []string{
							newNoteName,
						},
					}
					return newDailyData, nil
				}
			} else if insertLineNum <= len(dailyData.Notes) {
				// このdailyDataのnoteの途中に挿入
				if dailyData.Date != date {
					return usecases.DailyData{}, errors.Errorf("dateと挿入位置に矛盾があります。date: %v, dailyData: %v", date, dailyData)
				}
				newNotes = append(newNotes, dailyData.Notes[:insertLineNum-1]...)
				newNotes = append(newNotes, newNoteName)
				newNotes = append(newNotes, dailyData.Notes[insertLineNum-1:]...)
				dailyData.Notes = newNotes
				return dailyData, nil
			} else {
				return usecases.DailyData{}, errors.New("想定外")
			}
		}
	}
	// どこにも挿入されずにここまで来たということは末尾に新dailyDateを追加するということなので新dailyDataを作って返す
	if dailyList[len(dailyList)-1].Date == date {
		// 一つ前の末尾に追加
		dailyList[len(dailyList)-1].Notes = append(dailyList[len(dailyList)-1].Notes, newNoteName)
		return dailyList[len(dailyList)-1], nil
	} else {
		newDailyData := usecases.DailyData{
			Date: date,
			Notes: []string{
				newNoteName,
			},
		}
		return newDailyData, nil
	}
}
