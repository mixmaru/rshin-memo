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
