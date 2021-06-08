package views

import (
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

func TestDateRange_GetSomeDateInRange(t *testing.T) {
	t.Run("From, Toが設定されているRange：", func(t *testing.T) {
		d := &DateRange{
			From: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		}

		t.Run("指定numがRangeの範囲内ならnum個のDateが返ってくる", func(t *testing.T) {
			result, err := d.GetSomeDateInRange(5, time.Now())
			assert.NoError(t, err)
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
			result, err := d.GetSomeDateInRange(7, time.Now())
			assert.NoError(t, err)
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
			result, err := d.GetSomeDateInRange(10, time.Now())
			assert.NoError(t, err)
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
		d := &DateRange{
			From: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
		}

		t.Run("Fromからの指定num個のDateが返ってくる", func(t *testing.T) {
			result, err := d.GetSomeDateInRange(3, time.Now())
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("Toだけが設定されているRange：", func(t *testing.T) {
		d := &DateRange{
			To: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
		}
		t.Run("Toまでの指定num個のDateが返ってくる", func(t *testing.T) {
			result, err := d.GetSomeDateInRange(3, time.Now())
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2020, 12, 30, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, result)
		})

	})
	t.Run("ToもFromも設定されていないRangeの場合errorが返る：", func(t *testing.T) {
		d := &DateRange{}
		result, err := d.GetSomeDateInRange(4, time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local))
		assert.NoError(t, err)
		expected := []time.Time{
			time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 17, 0, 0, 0, 0, time.Local),
		}
		assert.Equal(t, expected, result)
	})
}
