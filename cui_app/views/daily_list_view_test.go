package views

import (
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDateRange_IsIn(t *testing.T) {
	t.Run("範囲内ならtrue、範囲外ならfalse", func(t *testing.T) {
		dateRange := &DateRange{
			From: time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		}
		inDate := time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local)
		assert.True(t, dateRange.IsIn(inDate))

		beforeDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)
		assert.False(t, dateRange.IsIn(beforeDate))

		afterDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
		assert.False(t, dateRange.IsIn(afterDate))
	})

	t.Run("From, Toと同日ならTrue", func(t *testing.T) {
		dateRange := &DateRange{
			From: time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		}
		equalBeforeDate := time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local)
		assert.True(t, dateRange.IsIn(equalBeforeDate))

		equalAfterDate := time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local)
		assert.True(t, dateRange.IsIn(equalAfterDate))
	})

	t.Run("From, To範囲が同日の場合", func(t *testing.T) {
		dateRange := &DateRange{
			From: time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		}

		t.Run("From, Toと同日ならTrue", func(t *testing.T) {
			equalDate := time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local)
			assert.True(t, dateRange.IsIn(equalDate))
		})

		t.Run("範囲外ならFalse", func(t *testing.T) {
			beforeDate := time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local)
			assert.False(t, dateRange.IsIn(beforeDate))
			afterDate := time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local)
			assert.False(t, dateRange.IsIn(afterDate))
		})
	})

	t.Run("Fromがゼロ値の場合", func(t *testing.T) {
		dateRange := &DateRange{
			To: time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		}
		beforeDate := time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local)
		assert.True(t, dateRange.IsIn(beforeDate))
		equalDate := time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local)
		assert.True(t, dateRange.IsIn(equalDate))
		afterDate := time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local)
		assert.False(t, dateRange.IsIn(afterDate))
	})

	t.Run("Toがゼロ値の場合", func(t *testing.T) {
		dateRange := &DateRange{
			From: time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		}
		beforeDate := time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local)
		assert.False(t, dateRange.IsIn(beforeDate))
		equalDate := time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local)
		assert.True(t, dateRange.IsIn(equalDate))
		afterDate := time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local)
		assert.True(t, dateRange.IsIn(afterDate))
	})
}

func TestIsEndOfDateList(t *testing.T) {
	dailyList := []usecases.DailyData{
		{
			Date: "2021-03-30",
			Notes: []string{
				"a",
				"b",
				"c",
			},
		},
		{
			Date: "2021-03-29",
			Notes: []string{
				"a",
				"b",
			},
		},
	}

	t.Run("日の末noteを指す番号だったらTrue", func(t *testing.T) {
		assert.True(t, IsEndOfDateList(2, dailyList))
		assert.True(t, IsEndOfDateList(4, dailyList))
	})

	t.Run("日の末noteを指す番号でなければfalse", func(t *testing.T) {
		assert.False(t, IsEndOfDateList(0, dailyList))
		assert.False(t, IsEndOfDateList(3, dailyList))
	})

	t.Run("範囲以外だったらfalse", func(t *testing.T) {
		assert.False(t, IsEndOfDateList(-1, dailyList))
		assert.False(t, IsEndOfDateList(10, dailyList))
	})
}

func TestIsFirstOfDateList(t *testing.T) {
	dailyList := []usecases.DailyData{
		{
			Date: "2021-03-30",
			Notes: []string{
				"a",
				"b",
				"c",
			},
		},
		{
			Date: "2021-03-29",
			Notes: []string{
				"a",
				"b",
			},
		},
	}

	t.Run("日の先頭noteを指す番号だったらTrue", func(t *testing.T) {
		assert.True(t, IsFirstOfDateList(0, dailyList))
		assert.True(t, IsFirstOfDateList(3, dailyList))
	})

	t.Run("日の末noteを指す番号でなければfalse", func(t *testing.T) {
		assert.False(t, IsFirstOfDateList(1, dailyList))
		assert.False(t, IsFirstOfDateList(2, dailyList))
		assert.False(t, IsFirstOfDateList(4, dailyList))
	})

	t.Run("範囲以外だったらfalse", func(t *testing.T) {
		assert.False(t, IsEndOfDateList(-1, dailyList))
		assert.False(t, IsEndOfDateList(10, dailyList))
	})
}

func Test_generateNewDailyData(t *testing.T) {
	dailyList := []usecases.DailyData{
		{
			Date: "2021-03-30",
			Notes: []string{
				"a",
				"b",
			},
		},
		{
			Date: "2021-03-29",
			Notes: []string{
				"a",
				"b",
			},
		},
		{
			Date: "2021-03-27",
			Notes: []string{
				"a",
				"b",
			},
		},
	}

	t.Run("先頭に新規追加", func(t *testing.T) {
		result, err := generateNewDailyData(dailyList, "newNote", "2021-04-01", 0)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-04-01",
			Notes: []string{
				"newNote",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("日付の先頭位置への挿入", func(t *testing.T) {
		result, err := generateNewDailyData(dailyList, "newNote", "2021-03-30", 0)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-03-30",
			Notes: []string{
				"newNote",
				"a",
				"b",
			},
		}
		assert.Equal(t, expected, result)

		result, err = generateNewDailyData(dailyList, "newNote", "2021-03-28", 2)
		assert.NoError(t, err)
		expected = usecases.DailyData{
			Date: "2021-03-29",
			Notes: []string{
				"newNote",
				"a",
				"b",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("日付の中間位置への挿入", func(t *testing.T) {
		result, err := generateNewDailyData(dailyList, "newNote", "2021-03-30", 1)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-03-30",
			Notes: []string{
				"a",
				"newNote",
				"b",
			},
		}
		assert.Equal(t, expected, result)

		result, err = generateNewDailyData(dailyList, "newNote", "2021-03-28", 3)
		assert.NoError(t, err)
		expected = usecases.DailyData{
			Date: "2021-03-29",
			Notes: []string{
				"a",
				"newNote",
				"b",
			},
		}
		assert.Equal(t, expected, result)
	})
	//
	//t.Run("日付の末尾位置への挿入", func(t *testing.T) {
	//	result, err := generateNewDailyData(dailyList, "newNote", "2021-03-30", 0)
	//	assert.NoError(t, err)
	//	expected := []usecases.DailyData{
	//		{
	//			Date: "2021-03-30",
	//			Notes: []string{
	//				"a",
	//				"b",
	//				"c",
	//				"newNote",
	//			},
	//		},
	//		{
	//			Date: "2021-03-29",
	//			Notes: []string{
	//				"a",
	//				"b",
	//			},
	//		},
	//	}
	//	assert.Equal(t, expected, result)
	//
	//	result, err = generateNewDailyData(dailyList, "newNote", "2021-03-29", 0)
	//	assert.NoError(t, err)
	//	expected = []usecases.DailyData{
	//		{
	//			Date: "2021-03-30",
	//			Notes: []string{
	//				"a",
	//				"b",
	//				"c",
	//			},
	//		},
	//		{
	//			Date: "2021-03-29",
	//			Notes: []string{
	//				"a",
	//				"b",
	//				"newNote",
	//			},
	//		},
	//	}
	//	assert.Equal(t, expected, result)
	//})
	//
	//t.Run("日付の先頭位置に別日の挿入", func(t *testing.T) {
	//	dailyList = []usecases.DailyData{
	//		{
	//			Date: "2021-03-30",
	//			Notes: []string{
	//				"a",
	//				"b",
	//				"c",
	//			},
	//		},
	//		{
	//			Date: "2021-03-28",
	//			Notes: []string{
	//				"a",
	//				"b",
	//			},
	//		},
	//	}
	//
	//	result, err := generateNewDailyData(dailyList, "newNote", "2021-04-01", 0)
	//	assert.NoError(t, err)
	//	expected := []usecases.DailyData{
	//		{
	//			Date: "2021-04-01",
	//			Notes: []string{
	//				"newNote",
	//			},
	//		},
	//		{
	//			Date: "2021-03-30",
	//			Notes: []string{
	//				"a",
	//				"b",
	//				"c",
	//			},
	//		},
	//		{
	//			Date: "2021-03-28",
	//			Notes: []string{
	//				"a",
	//				"b",
	//			},
	//		},
	//	}
	//	assert.Equal(t, expected, result)
	//
	//	result, err = generateNewDailyData(dailyList, "newNote", "2021-03-29", 3)
	//	assert.NoError(t, err)
	//	expected = []usecases.DailyData{
	//		{
	//			Date: "2021-03-30",
	//			Notes: []string{
	//				"a",
	//				"b",
	//				"c",
	//			},
	//		},
	//		{
	//			Date: "2021-03-29",
	//			Notes: []string{
	//				"newNote",
	//			},
	//		},
	//		{
	//			Date: "2021-03-28",
	//			Notes: []string{
	//				"a",
	//				"b",
	//			},
	//		},
	//	}
	//	assert.Equal(t, expected, result)
	//})
	//
	//t.Run("日付の末尾位置に別日の挿入", func(t *testing.T) {
	//	dailyList = []usecases.DailyData{
	//		{
	//			Date: "2021-03-30",
	//			Notes: []string{
	//				"a",
	//				"b",
	//				"c",
	//			},
	//		},
	//		{
	//			Date: "2021-03-28",
	//			Notes: []string{
	//				"a",
	//				"b",
	//			},
	//		},
	//	}
	//
	//	result, err := generateNewDailyData(dailyList, "newNote", "2021-03-27", 5)
	//	assert.NoError(t, err)
	//	expected := []usecases.DailyData{
	//		{
	//			Date: "2021-03-30",
	//			Notes: []string{
	//				"a",
	//				"b",
	//				"c",
	//			},
	//		},
	//		{
	//			Date: "2021-03-28",
	//			Notes: []string{
	//				"a",
	//				"b",
	//			},
	//		},
	//		{
	//			Date: "2021-03-27",
	//			Notes: []string{
	//				"newNote",
	//			},
	//		},
	//	}
	//	assert.Equal(t, expected, result)
	//})
	//
	//t.Run("異なる日付の位置に挿入しようとするとエラー", func(t *testing.T) {
	//	_, err := generateNewDailyData(dailyList, "newNote", "2021-01-27", 2)
	//	assert.Error(t, err)
	//})
}
