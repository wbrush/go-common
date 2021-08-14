package helpers

import (
	"reflect"
	"testing"
	"time"
)

func TestAddTime(t *testing.T) {
	type args struct {
		inTime   time.Time
		years    int
		months   int
		days     int
		hours    int
		timeZone string
	}
	tests := []struct {
		name        string
		args        args
		wantOutTime time.Time
		wantErr     bool
	}{
		{name: "test 0", args: args{
			inTime:   time.Unix(1613137108, 0),
			years:    0,
			months:   0,
			days:     0,
			hours:    24,
			timeZone: "",
		}, wantOutTime: time.Unix(1613223508, 0), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutTime, err := AddTime(tt.args.inTime, tt.args.years, tt.args.months, tt.args.days, tt.args.hours, tt.args.timeZone)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutTime, tt.wantOutTime) {
				t.Errorf("AddTime() gotOutTime = %v, want %v", gotOutTime, tt.wantOutTime)
			}
		})
	}
}
