package dto

import (
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
)

type InsertData struct {
	insertNum       int    // dailyDataの上から何番目に挿入するか
	noteName        string // note名
	targetDailyData usecases.DailyData
}

func (i *InsertData) TargetDailyData() usecases.DailyData {
	return i.targetDailyData
}

func (i *InsertData) SetTargetDailyData(targetDailyData usecases.DailyData) {
	i.targetDailyData = targetDailyData
}

func (i *InsertData) NoteName() string {
	return i.noteName
}

func (i *InsertData) SetNoteName(noteName string) {
	i.noteName = noteName
}

func (i *InsertData) InsertNum() int {
	return i.insertNum
}

func (i *InsertData) SetInsertNum(insertNum int) error {
	if insertNum > len(i.targetDailyData.Notes) {
		return errors.New("不整合な値がsetされました")
	}
	i.insertNum = insertNum
	return nil
}

func (i *InsertData) GenerateNewDailyData() usecases.DailyData {
	var retDailyData usecases.DailyData
	retDailyData.Date = i.targetDailyData.Date
	for index, noteName := range i.targetDailyData.Notes {
		if index == i.InsertNum() {
			retDailyData.Notes = append(retDailyData.Notes, i.NoteName())
		}
		retDailyData.Notes = append(retDailyData.Notes, noteName)
	}
	if i.insertNum == len(i.targetDailyData.Notes) {
		// 末への追加
		retDailyData.Notes = append(retDailyData.Notes, i.NoteName())
	}
	return retDailyData
}
