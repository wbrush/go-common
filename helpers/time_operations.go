package helpers

import (
	"strings"
	"time"
)

func AddTime(inTime time.Time, years, months, days, hours int, timeZone string) (outTime time.Time, err error) {
	var loc *time.Location
	loc, err = time.LoadLocation(timeZone)
	//if err != nil { return outTime, err }     //  not sure that we should return here;

	if len(strings.TrimSpace(timeZone)) <= 0 || err != nil {
		outTime = inTime.AddDate(years, months, days)
		outTime = outTime.Add(time.Duration(hours) * time.Hour)
	} else {
		t2 := inTime.In(loc)
		t1 := t2.AddDate(years, months, days)
		t1 = t1.Add(time.Duration(hours) * time.Hour)

		outTime = t1.UTC()
	}

	return
}
