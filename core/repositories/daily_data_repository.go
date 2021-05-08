package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type DailyDataRepositoryInterface interface {
	Save(entity *entities.DailyDataEntity) error
}

// DailyDataをjsonファイルで永続化するリポジトリ
type DailyDataRepository struct {
	filePath string //jsonファイルの絶対パス
}

func NewDailyDataRepository(filePath string) *DailyDataRepository {
	return &DailyDataRepository{filePath: filePath}
}

func (d *DailyDataRepository) Save(entity *entities.DailyDataEntity) error {
	// ファイルがなければ空データで作成する
	if _, err := os.Stat(d.filePath); err != nil {
		fd, err := os.Create(d.filePath)
		if err != nil {
			return errors.Wrapf(err, "file作成失敗。%v", d.filePath)
		}
		text, err := json.Marshal([]DailyData{})
		if err != nil {
			return errors.Wrap(err, "空json作成失敗。")
		}
		_, err = fmt.Fprint(fd, string(text))
		if err != nil {
			return errors.Wrapf(err, "書き込み失敗。%v", string(text))
		}
		err = fd.Close()
		if err != nil {
			return errors.Wrap(err, "ファイルClose失敗。")
		}
	}

	fd, err := os.OpenFile(d.filePath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return errors.Wrapf(err, "fileopen失敗。%v", d.filePath)
	}
	// json読み込み
	bytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return errors.Wrap(err, "jsonファイル読み込み失敗")
	}
	var dailyDataList []DailyData
	err = json.Unmarshal(bytes, &dailyDataList)
	if err != nil {
		return errors.Wrapf(err, "jsonパース失敗。%v", string(bytes))
	}

	// データ書き換え
	if len(dailyDataList) == 0 {
		dailyData := DailyData{
			Date:  entity.Date().Format("2006-01-02"),
			Notes: entity.NoteNames(),
		}
		dailyDataList = append(dailyDataList, dailyData)
	} else {
		inserted := false
		for i, dailyData := range dailyDataList {
			date, err := time.ParseInLocation("2006-01-02", dailyData.Date, time.Local)
			if err != nil {
				return errors.Wrapf(err, "日付パース失敗。%+v", dailyData)
			}
			if date.After(entity.Date()) {
				// 同じ日付にヒットするまで過去にさかのぼって捜査する
				continue
			} else if date.Equal(entity.Date()) {
				// 日付が同じだったら上書きする
				newDailyData := DailyData{
					Date:  entity.Date().Format("2006-01-02"),
					Notes: entity.NoteNames(),
				}
				dailyDataList[i] = newDailyData
				inserted = true
				break
			} else {
				// 日付が過去だったら、そこに挿入する
				newDailyData := DailyData{
					Date:  entity.Date().Format("2006-01-02"),
					Notes: entity.NoteNames(),
				}
				insertedDailyList := []DailyData{}
				insertedDailyList = append(insertedDailyList, dailyDataList[:i]...)
				insertedDailyList = append(insertedDailyList, newDailyData)
				insertedDailyList = append(insertedDailyList, dailyDataList[i:]...)
				dailyDataList = insertedDailyList
				inserted = true
				break
			}
		}
		if !inserted {
			// 上記処理で挿入されなかったのなら末にappendされる
			newDailyData := DailyData{
				Date:  entity.Date().Format("2006-01-02"),
				Notes: entity.NoteNames(),
			}
			dailyDataList = append(dailyDataList, newDailyData)
		}
	}

	// jsonファイル出力
	newText, err := json.Marshal(dailyDataList)
	if err != nil {
		return errors.Errorf("json Mar")
	}
	// tmpファイルに書き出す。
	tmpFilePath := filepath.Join(filepath.Dir(d.filePath), "tmp_daily_data.json")
	err = ioutil.WriteFile(tmpFilePath, newText, 0644)
	if err != nil {
		return errors.Wrapf(err, "tmpファイルの作成失敗。%v", tmpFilePath)
	}

	// 読み込みファイルを削除する
	fd.Close()
	err = os.Remove(d.filePath)
	if err != nil {
		return errors.Wrapf(err, "ファイルの削除失敗。%v", d.filePath)
	}

	// tmpファイルの名前を変更する
	err = os.Rename(tmpFilePath, d.filePath)
	if err != nil {
		return errors.Wrapf(err, "ファイルのリネーム失敗。%v => %v", tmpFilePath, d.filePath)
	}

	return nil
}

type DailyData struct {
	Date  string   `json:"date"`
	Notes []string `json:"notes"`
}
