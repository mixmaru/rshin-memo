package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
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
	defer fd.Close()
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
	insertIndex := -1
	for i, dailyData := range dailyDataList {
		date, err := time.Parse("2006-01-02", dailyData.Date)
		if err != nil {
			return errors.Wrapf(err, "日付パース失敗。%+v", dailyData)
		}
		if date.Equal(entity.Date()) {
			// 日付が同じだったら上書きする
			newDailyData := DailyData{
				Date:  entity.Date().Format("2006-01-02"),
				Notes: entity.NoteNames(),
			}
			dailyData = newDailyData
			break
		} else if date.After(entity.Date()) {

		}
		// else if 引数のdailyDataよりも未来の日付である
		//     continue
		// else
		//		挿入point := i
	}
	if insertIndex > 0 {
		// 指定位置へdataを挿入
	}

	// jsonファイル出力

	return nil
}

type DailyData struct {
	Date  string   `json:"date"`
	Notes []string `json:"notes"`
}
