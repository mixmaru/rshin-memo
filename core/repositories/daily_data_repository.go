package repositories

import (
	bytes2 "bytes"
	"encoding/json"
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

// entityをjsonfileへ永続化する
// スレッドセーフではない
func (d *DailyDataRepository) Save(entity *entities.DailyDataEntity) error {
	// ファイルがなければ空データで作成する
	if _, err := os.Stat(d.filePath); err != nil {
		err := ioutil.WriteFile(d.filePath, []byte("[]"), 0644)
		if err != nil {
			return errors.Wrapf(err, "ファイル作成失敗。d.filePath: %v", d.filePath)
		}
	}

	// ファイル読み込み（一人でしかつかわないのでスレッドセーフは考慮してない）
	jsonTextBytes, err := ioutil.ReadFile(d.filePath)
	if err != nil {
		return errors.Wrapf(err, "ファイル読み込み失敗。d.filePath%v", d.filePath)
	}

	// jsonパース
	var dailyDataList []DailyData
	err = json.Unmarshal(jsonTextBytes, &dailyDataList)
	if err != nil {
		return errors.Wrapf(err, "jsonパース失敗。%v", string(jsonTextBytes))
	}

	updatedDailyDataList, err := generateNewDailyDataList(dailyDataList, entity)
	if err != nil {
		return err
	}

	// jsonに変換
	var buf bytes2.Buffer
	newText, err := json.Marshal(updatedDailyDataList)
	if err != nil {
		return errors.Wrapf(err, "json Marshal error. updatedDailyDataList: %+v", updatedDailyDataList)
	}
	// jsonを整形する
	err = json.Indent(&buf, newText, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "json indent error newText: %+v", newText)
	}
	// ファイルに書き出す。
	err = ioutil.WriteFile(d.filePath, buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrapf(err, "ファイルへの書き出し失敗。d.filePath: %v, buf.Bytes(): %v", d.filePath, buf.Bytes())
	}
	return nil
}

// currentDailyDataListの適切な位置にentityが表すDailyDataを上書き、もしくは挿入したものを返す
func generateNewDailyDataList(currentDailyDataList []DailyData, entity *entities.DailyDataEntity) ([]DailyData, error) {
	newDailyData := DailyData{
		Date:  entity.Date().Format("2006-01-02"),
		Notes: entity.NoteNames(),
	}

	inserted := false
	for i, dailyData := range currentDailyDataList {
		date, err := time.ParseInLocation("2006-01-02", dailyData.Date, time.Local)
		if err != nil {
			return nil, errors.Wrapf(err, "日付パース失敗。%+v", dailyData)
		}

		if date.After(entity.Date()) {
			// 同じ日付にヒットするまで過去にさかのぼって捜査する
			continue
		} else if date.Equal(entity.Date()) {
			// 日付が同じだったら上書きする
			currentDailyDataList[i] = newDailyData
			inserted = true
			break
		} else {
			// 日付が過去だったら、そこに挿入する
			insertedDailyList := []DailyData{}
			insertedDailyList = append(insertedDailyList, currentDailyDataList[:i]...)
			insertedDailyList = append(insertedDailyList, newDailyData)
			insertedDailyList = append(insertedDailyList, currentDailyDataList[i:]...)
			currentDailyDataList = insertedDailyList
			inserted = true
			break
		}
	}
	if !inserted {
		// 上記処理で挿入されなかったのなら末にappendされる
		currentDailyDataList = append(currentDailyDataList, newDailyData)
	}
	return currentDailyDataList, nil
}

type DailyData struct {
	Date  string   `json:"date"`
	Notes []string `json:"notes"`
}
