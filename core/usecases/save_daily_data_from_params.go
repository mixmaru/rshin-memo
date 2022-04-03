package usecases

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/pkg/errors"
	"time"
)

type SaveDailyDataFromParamsUseCase struct {
	noteRepository      repositories.NoteRepositoryInterface
	dailyDataRepository repositories.DailyDataRepositoryInterface
}

func NewSaveDailyDataFromParamsUseCase(noteRepository repositories.NoteRepositoryInterface, dailyDataRepository repositories.DailyDataRepositoryInterface) *SaveDailyDataFromParamsUseCase {
	return &SaveDailyDataFromParamsUseCase{noteRepository: noteRepository, dailyDataRepository: dailyDataRepository}
}

func (c *SaveDailyDataFromParamsUseCase) Handle(
	baseMemoDate time.Time,
	baseMemoName string,
	newMemoDate time.Time,
	newMemoName string,
	newMemoContent string,
	mode InsertMode,
) error {
	// memoファイルを作成保存
	noteEntity := entities.NewNoteEntity(newMemoName, newMemoContent)
	err := c.noteRepository.Save(noteEntity)
	if err != nil {
		return err
	}

	// 対象のdailyDataEntityを作成
	dailyDataEntity, err := c.generateNewDailyDataEntity(newMemoDate, newMemoName, baseMemoDate, baseMemoName, mode)
	if err != nil {
		return err
	}

	// dailyDataEntityを保存
	err = c.dailyDataRepository.Save(dailyDataEntity)
	if err != nil {
		// todo: もし失敗したらファイルを削除site modosu
		return err
	}
	return nil

	//

	//println(allDailyData)
	//noteEntity := entities.NewNoteEntity("new_memo_name", "new_memo_内容")
	//saveingEntity := entities.NewDailyDataEntity(
	//	time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
	//	[]string{
	//		"newMemoName",
	//	},
	//)
	//c.dailyDataRepository.Save(saveingEntity)
	//c.noteRepository.Save(noteEntity)
	/*
		jsonファイルに入力


		create memo file
	*/

	//// dailyListの作成
	//date, err := time.ParseInLocation("2006-01-02", dailyData.Date, time.Local)
	//if err != nil {
	//	return errors.Wrapf(err, "日付Parseの失敗。%+v", dailyData)
	//}
	//newDailyDataEntity := entities.NewDailyDataEntity(
	//	date,
	//	dailyData.Notes,
	//)
	//err = c.dailyDataRepository.Save(newDailyDataEntity)
	//if err != nil {
	//	return err
	//}
	//
	//// Noteの作成（存在しないNoteであれば作成する）
	//for _, noteName := range dailyData.Notes {
	//	noteEntity, err := c.noteRepository.GetByNoteName(noteName)
	//	if err != nil {
	//		return err
	//	}
	//	if noteEntity == nil {
	//		newNote := entities.NewNoteEntity(noteName, "")
	//		err := c.noteRepository.Save(newNote)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}
	//return nil
}

// todo: リファクタ
func (c *SaveDailyDataFromParamsUseCase) generateNewDailyDataEntity(newMemoDate time.Time, newMemoname string, baseMemoDate time.Time, baseMemoName string, mode InsertMode) (*entities.DailyDataEntity, error) {
	// データ取得
	allDailyData, err := c.dailyDataRepository.Get()
	if err != nil {
		return nil, err
	}
	var pre *entities.DailyDataEntity
	for i, dailyData := range allDailyData {
		// 上からたどっていってベース日までたどる。
		if !dailyData.Date().Equal(baseMemoDate) {
			pre = dailyData
			continue
		}
		switch mode {
		case INSERT_NEWER_MODE:
			// 上挿入の場合
			for j, note := range dailyData.NoteNames() {
				// ベースメモ名までたどる。このとき番号を数えておく
				if note == baseMemoName {
					if j != 0 {
						// 番号が0(先頭)でなければ一つ前に挿入する。
						dailyData.InsertNoteName(newMemoname, j)
					} else {
						// 番号が0であれば挿入日付を確認し、同じでなければ一つ前のやつの日付を確認する
						if dailyData.Date() == newMemoDate {
							// 同じであればそこの先頭に追加する
							dailyData.InsertNoteName(newMemoname, 0)
						} else if pre != nil && pre.Date() == newMemoDate {
							// 同じでなければ一つ前のやつの日付を確認する。
							//同じであればそこの末に追加する
							dailyData.InsertNoteNameToLast(newMemoname)
						} else {
							// 異なれば新しいdailyDataを作成する
							dailyData = entities.NewDailyDataEntity(newMemoDate, []string{newMemoname})
						}
					}
					return dailyData, nil
				}
			}
			return nil, errors.Errorf("想定外. dailyData:%v, newMemoDate:%v, newMemoname:%v, baseMemoDate:%v, baseMemoName:%v", dailyData, newMemoDate, newMemoname, baseMemoDate, baseMemoName)
		case INSERT_OLDER_MODE:
			// 下挿入の場合
			for j, note := range dailyData.NoteNames() {
				// ベースメモ名までたどる。このとき番号を数えておく
				if note == baseMemoName {
					if j != len(dailyData.NoteNames())-1 {
						// 番号が末でなければ一つ後に挿入する。
						dailyData.InsertNoteName(newMemoname, j+1)
						return dailyData, nil
					} else {
						// 番号が末であれば挿入日付を確認し、
						if dailyData.Date() == newMemoDate {
							// 同じであればそこの末に追加する
							dailyData.InsertNoteName(newMemoname, j+1)
							return dailyData, nil
						} else if i+1 < len(allDailyData) && allDailyData[i+1].Date().Equal(newMemoDate) {
							// 同じでなければ次のやつの日付を確認する。
							// 同じであればそこの先頭に追加する
							allDailyData[j+1].InsertNoteName(newMemoname, 0)
							return allDailyData[i+1], nil
						} else {
							// 異なれば新しいdailyDataを作成する
							return entities.NewDailyDataEntity(newMemoDate, []string{newMemoname}), nil
						}
					}
				}
			}
			return nil, errors.Errorf("想定外. dailyData:%v, newMemoDate:%v, newMemoname:%v, baseMemoDate:%v, baseMemoName:%v", dailyData, newMemoDate, newMemoname, baseMemoDate, baseMemoName)
		default:
			return nil, errors.Errorf("想定外エラー mode: %v", mode)
		}
	}
	/*
	   まずベース日を確認する。
	   上からたどっていってベース日までたどる。
	   で、ベースメモ名までたどる。このとき番号を数えておく
	   で、上挿入か下挿入か見る
	   上挿入の場合
	       番号が0でなければ一つ前に挿入する。
	       番号が0であれば挿入日付を確認し、同じでなければ一つ前のやつの日付を確認する
	       一つ前のやつと一致していれば、
	           そいつの末に追加する
	       一致していなければ、
	           あたらしくデータを作ってそこに一つ追加する

	   下挿入の場合
	       番号が末でなければ一つあとに挿入する
	       番号が末であれば挿入日付を確認する
	           同じである
	               一番末に挿入する
	           違う
	               次のやつの日付を確認して追加日付を比較する
	                   同じである
	                       そこの先頭に追加する
	                   違う
	                       新しくデータを作ってそこに一つだけ追加する
	*/
	return nil, errors.Errorf("想定外. newMemoDate:%v, newMemoname:%v, baseMemoDate:%v, baseMemoName:%v", newMemoDate, newMemoname, baseMemoDate, baseMemoName)
}
