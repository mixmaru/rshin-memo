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

func createTestDailyList() []usecases.DailyData {
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
	return dailyList
}

func Test_generateNewDailyData(t *testing.T) {
	t.Run("先頭に新規追加", func(t *testing.T) {
		dailyList := createTestDailyList()
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
		dailyList := createTestDailyList()
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

		dailyList = createTestDailyList()
		result, err = generateNewDailyData(dailyList, "newNote", "2021-03-29", 2)
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
		dailyList := createTestDailyList()
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

		dailyList = createTestDailyList()
		result, err = generateNewDailyData(dailyList, "newNote", "2021-03-29", 3)
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

	t.Run("日付の末尾位置への挿入", func(t *testing.T) {
		dailyList := createTestDailyList()
		result, err := generateNewDailyData(dailyList, "newNote", "2021-03-30", 2)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-03-30",
			Notes: []string{
				"a",
				"b",
				"newNote",
			},
		}
		assert.Equal(t, expected, result)

		dailyList = createTestDailyList()
		result, err = generateNewDailyData(dailyList, "newNote", "2021-03-29", 4)
		assert.NoError(t, err)
		expected = usecases.DailyData{
			Date: "2021-03-29",
			Notes: []string{
				"a",
				"b",
				"newNote",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("日付の先頭位置に別日の挿入", func(t *testing.T) {
		dailyList := createTestDailyList()
		result, err := generateNewDailyData(dailyList, "newNote", "2021-04-01", 0)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-04-01",
			Notes: []string{
				"newNote",
			},
		}
		assert.Equal(t, expected, result)

		dailyList = createTestDailyList()
		result, err = generateNewDailyData(dailyList, "newNote", "2021-03-28", 4)
		assert.NoError(t, err)
		expected = usecases.DailyData{
			Date: "2021-03-28",
			Notes: []string{
				"newNote",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("日付の末尾位置に別日の挿入", func(t *testing.T) {
		dailyList := createTestDailyList()
		result, err := generateNewDailyData(dailyList, "newNote", "2021-03-26", 6)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-03-26",
			Notes: []string{
				"newNote",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("異なる日付の位置に挿入しようとするとエラー", func(t *testing.T) {
		dailyList := createTestDailyList()
		_, err := generateNewDailyData(dailyList, "newNote", "2021-01-27", 1)
		assert.Error(t, err)
	})
}

func TestDateRange_GetSomeDateInRange(t *testing.T) {
	t.Run("From, Toが設定されているRange：", func(t *testing.T) {
		d := &DateRange{
			From: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		}

		t.Run("指定numがRangeの範囲内ならnum個のDateが返ってくる", func(t *testing.T) {
			result := d.GetSomeDateInRange(5)
			expected := []time.Time{
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, result)
		})

		t.Run("指定numがRangeの範囲と同じならnum個のDateが返ってくる", func(t *testing.T) {
			result := d.GetSomeDateInRange(7)
			expected := []time.Time{
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, result)
		})

		t.Run("指定numがRangeの範囲外ならRangeの最大個のDateが返ってくる", func(t *testing.T) {
			result := d.GetSomeDateInRange(10)
			expected := []time.Time{
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("Fromだけが設定されているRange：", func(t *testing.T) {

	})

	t.Run("Toだけが設定されているRange：", func(t *testing.T) {

	})
}
