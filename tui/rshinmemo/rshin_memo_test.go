package main

import (
	"testing"
	"time"
)

func TestIsDateInRange(t *testing.T) {
	type args struct {
		date time.Time
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"dateがfromよりも前",
			args{
				date: time.Date(2021, 4, 15, 0, 0, 0, 0, time.Local),
				from: time.Date(2021, 4, 16, 0, 0, 0, 0, time.Local),
				to:   time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"dateがfromと同じ",
			args{
				date: time.Date(2021, 4, 16, 0, 0, 0, 0, time.Local),
				from: time.Date(2021, 4, 16, 0, 0, 0, 0, time.Local),
				to:   time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
			},
			true,
		},
		{
			"dateがfromより後、toより前",
			args{
				date: time.Date(2021, 4, 17, 0, 0, 0, 0, time.Local),
				from: time.Date(2021, 4, 16, 0, 0, 0, 0, time.Local),
				to:   time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
			},
			true,
		},
		{
			"dateがtoと同じ",
			args{
				date: time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
				from: time.Date(2021, 4, 16, 0, 0, 0, 0, time.Local),
				to:   time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
			},
			true,
		},
		{
			"dateがtoより後",
			args{
				date: time.Date(2021, 4, 19, 0, 0, 0, 0, time.Local),
				from: time.Date(2021, 4, 16, 0, 0, 0, 0, time.Local),
				to:   time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"fromがゼロ値、dateがtoより前",
			args{
				date: time.Date(2021, 4, 17, 0, 0, 0, 0, time.Local),
				from: time.Time{},
				to:   time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
			},
			true,
		},
		{
			"fromがゼロ値、dateがtoと同じ",
			args{
				date: time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
				from: time.Time{},
				to:   time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
			},
			true,
		},
		{
			"fromがゼロ値、dateがtoより後",
			args{
				date: time.Date(2021, 4, 19, 0, 0, 0, 0, time.Local),
				from: time.Time{},
				to:   time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"toがゼロ値、dateがfromより前",
			args{
				date: time.Date(2021, 4, 16, 0, 0, 0, 0, time.Local),
				from: time.Date(2021, 4, 17, 0, 0, 0, 0, time.Local),
				to:   time.Time{},
			},
			false,
		},
		{
			"toがゼロ値、dateがfromと同じ",
			args{
				date: time.Date(2021, 4, 17, 0, 0, 0, 0, time.Local),
				from: time.Date(2021, 4, 17, 0, 0, 0, 0, time.Local),
				to:   time.Time{},
			},
			true,
		},
		{
			"toがゼロ値、dateがfromより後",
			args{
				date: time.Date(2021, 4, 18, 0, 0, 0, 0, time.Local),
				from: time.Date(2021, 4, 17, 0, 0, 0, 0, time.Local),
				to:   time.Time{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDateInRange(tt.args.date, tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("IsDateInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
