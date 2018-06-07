package math

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var firstDayMonday bool
var timeFormats = []string{"1/2/2006", "1/2/2006 15:4:5", "2006-1-2 15:4:5", "2006-1-2 15:4", "2006-1-2", "1-2", "15:4:5", "15:4", "15", "15:4:5 Jan 2, 2006 MST", "2006-01-02 15:04:05.999999999 -0700 MST"}

// lay gio hien tai he thong
func GetTimeNowVietNam() time.Time {
	return time.Now().In(loc)
}

type Now struct {
	time.Time
}

func New(t time.Time) *Now {
	return &Now{t}
}

func BeginningOfMinute() time.Time {
	return New(time.Now()).BeginningOfMinute()
}

func BeginningOfHour() time.Time {
	return New(time.Now()).BeginningOfHour()
}

func BeginningOfDay() time.Time {
	return New(time.Now()).BeginningOfDay()
}

func BeginningOfWeek() time.Time {
	return New(time.Now()).BeginningOfWeek()
}

func BeginningOfMonth() time.Time {
	return New(time.Now()).BeginningOfMonth()
}

func BeginningOfQuarter() time.Time {
	return New(time.Now()).BeginningOfQuarter()
}

func BeginningOfYear() time.Time {
	return New(time.Now()).BeginningOfYear()
}

func EndOfMinute() time.Time {
	return New(time.Now()).EndOfMinute()
}

func EndOfHour() time.Time {
	return New(time.Now()).EndOfHour()
}

func EndOfDay() time.Time {
	return New(time.Now()).EndOfDay()
}

func EndOfWeek() time.Time {
	return New(time.Now()).EndOfWeek()
}

func EndOfMonth() time.Time {
	return New(time.Now()).EndOfMonth()
}

func EndOfQuarter() time.Time {
	return New(time.Now()).EndOfQuarter()
}

func EndOfYear() time.Time {
	return New(time.Now()).EndOfYear()
}

func Monday() time.Time {
	return New(time.Now()).Monday()
}

func Sunday() time.Time {
	return New(time.Now()).Sunday()
}

func EndOfSunday() time.Time {
	return New(time.Now()).EndOfSunday()
}

func Parse(strs ...string) (time.Time, error) {
	return New(time.Now()).Parse(strs...)
}

func MustParse(strs ...string) time.Time {
	return New(time.Now()).MustParse(strs...)
}

func Between(time1, time2 string) bool {
	return New(time.Now()).Between(time1, time2)
}

func (now *Now) BeginningOfMinute() time.Time {
	return now.Truncate(time.Minute)
}

func (now *Now) BeginningOfHour() time.Time {
	return now.Truncate(time.Hour)
}

func (now *Now) BeginningOfDay() time.Time {
	d := time.Duration(-now.Hour()) * time.Hour
	return now.BeginningOfHour().Add(d)
}

func BeginningOfDayInt64(date int64) time.Time {
	var now = time.Unix(date, 0)
	return New(now).BeginningOfDay()
}

func (now *Now) BeginningOfWeek() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if firstDayMonday {
		if weekday == 0 {
			weekday = 7
		}
		weekday = weekday - 1
	}

	d := time.Duration(-weekday) * 24 * time.Hour
	return t.Add(d)
}

func (now *Now) BeginningOfMonth() time.Time {
	t := now.BeginningOfDay()
	d := time.Duration(-int(t.Day())+1) * 24 * time.Hour
	return t.Add(d)
}

func (now *Now) BeginningOfQuarter() time.Time {
	month := now.BeginningOfMonth()
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

func (now *Now) BeginningOfYear() time.Time {
	t := now.BeginningOfDay()
	d := time.Duration(-int(t.YearDay())+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func (now *Now) EndOfMinute() time.Time {
	return now.BeginningOfMinute().Add(time.Minute - time.Nanosecond)
}

func (now *Now) EndOfHour() time.Time {
	return now.BeginningOfHour().Add(time.Hour - time.Nanosecond)
}

func (now *Now) EndOfDay() time.Time {
	return now.BeginningOfDay().Add(24*time.Hour - time.Nanosecond)
}

func (now *Now) EndOfWeek() time.Time {
	return now.BeginningOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

func (now *Now) EndOfMonth() time.Time {
	return now.BeginningOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func (now *Now) EndOfQuarter() time.Time {
	return now.BeginningOfQuarter().AddDate(0, 3, 0).Add(-time.Nanosecond)
}

func (now *Now) EndOfYear() time.Time {
	return now.BeginningOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

func CompareDayWeek(dateTime1 time.Time, date2 int64) int {
	var dateTime2 = time.Unix(date2, 0)
	if dateTime1.Weekday() == dateTime2.Weekday() {
		return 0
	}
	return -1
}

func CompareDayTime(dateTime1 time.Time, date2 int64) int {
	var dateTime2 = time.Unix(date2, 0)
	var year1, month1, day1 = dateTime1.Date()
	var year2, month2, day2 = dateTime2.Date()
	if year1 == year2 && month1 == month2 && day1 == day2 {
		return 0
	}
	return -1
}

func CompareWeekDay(date1 int64, date2 int64) (isCheck bool) {
	var dateTime1 = time.Unix(date1, 0)
	var dateTime2 = time.Unix(date2, 0)
	if dateTime1.Weekday() == dateTime2.Weekday() {
		isCheck = true
	}
	return
}

func (now *Now) Monday() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := time.Duration(-weekday+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func (now *Now) Sunday() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if weekday == 0 {
		return t
	}
	d := time.Duration(7-weekday) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)

}

func (now *Now) EndOfSunday() time.Time {
	return now.Sunday().Add(24*time.Hour - time.Nanosecond)
}

func parseWithFormat(str string) (t time.Time, err error) {
	for _, format := range timeFormats {
		t, err = time.Parse(format, str)
		if err == nil {
			return
		}
	}
	err = errors.New("Can't parse string as time: " + str)
	return
}

func (now *Now) Parse(strs ...string) (t time.Time, err error) {
	var setCurrentTime bool
	parseTime := []int{}
	currentTime := []int{now.Second(), now.Minute(), now.Hour(), now.Day(), int(now.Month()), now.Year()}
	currentLocation := now.Location()

	for _, str := range strs {
		onlyTime := regexp.MustCompile(`^\s*\d+(:\d+)*\s*$`).MatchString(str) // match 15:04:05, 15

		t, err = parseWithFormat(str)
		location := t.Location()
		if location.String() == "UTC" {
			location = currentLocation
		}

		if err == nil {
			parseTime = []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
			onlyTime = onlyTime && (parseTime[3] == 1) && (parseTime[4] == 1)

			for i, v := range parseTime {
				// Don't reset hour, minute, second if it is a time only string
				if onlyTime && i <= 2 {
					continue
				}

				// Fill up missed information with current time
				if v == 0 {
					if setCurrentTime {
						parseTime[i] = currentTime[i]
					}
				} else {
					setCurrentTime = true
				}

				// Default day and month is 1, fill up it if missing it
				if onlyTime {
					if i == 3 || i == 4 {
						parseTime[i] = currentTime[i]
						continue
					}
				}
			}
		}

		if len(parseTime) > 0 {
			t = time.Date(parseTime[5], time.Month(parseTime[4]), parseTime[3], parseTime[2], parseTime[1], parseTime[0], 0, location)
			currentTime = []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
		}
	}
	return
}

func (now *Now) MustParse(strs ...string) (t time.Time) {
	t, err := now.Parse(strs...)
	if err != nil {
		panic(err)
	}
	return t
}

func (now *Now) Between(time1, time2 string) bool {
	restime := now.MustParse(time1)
	restime2 := now.MustParse(time2)
	return now.After(restime) && now.Before(restime2)
}

func DateDiff(s, e int64) []string {
	var res = make([]string, 0)
	var start = time.Unix(s, 0)
	var end = time.Unix(e, 0)
	if start.YearDay() == time.Now().YearDay() || end.YearDay() == start.YearDay() {
		return []string{strings.Split(start.String(), " ")[0]}
	}
	for {
		if end.YearDay() == start.YearDay() {
			break
		}
		s += 24 * 3600
		start = time.Unix(s, 0)
		res = append(res, strings.Split(start.String(), " ")[0])
	}
	return res
}

func ConvertTimeToString(timex int64) string {
	var timeConvert = time.Unix(timex, 0)
	return timeConvert.Format("02-01-2006")
}

func HourMinute() float32 {
	var timeNow = time.Now().In(loc)
	return float32(timeNow.Hour()) + float32(timeNow.Minute())/60
}

var loc, _ = time.LoadLocation("Asia/Ho_Chi_Minh")

func BeginningOfDayVN() time.Time {
	return BeginningOfDay().In(loc)
}
func NewTimeVN() time.Time {
	return new(time.Time).In(loc)
}
func HourMinuteEpoch(timeValue int64) float32 {
	var timeNow = time.Unix(timeValue, 0).In(loc)
	return float32(timeNow.Hour()) + float32(timeNow.Minute())/60
}

func ConvertTimeEpochToWeek(timeValue int64) time.Weekday {
	return time.Unix(timeValue, 0).Weekday()
}

func GetDayCount(dateStart int64, dateEnd int64) (dayCount int, dateTime1 time.Time, dateTime2 time.Time) {
	dateTime1 = time.Unix(dateStart, 0)
	fmt.Println(dateTime1)
	dateTime2 = time.Unix(dateEnd, 0)
	fmt.Println(dateTime2)
	dayCount = int(dateTime2.Sub(dateTime1).Hours() / 24)
	return
}

func BeginAndEndDay(val int64) (start int64, end int64) {
	var timeNow = time.Unix(val, 0).In(loc)
	start = New(timeNow).BeginningOfDay().Unix()
	end = New(timeNow).EndOfDay().Unix()
	return
}

func BeginAndEndDayNow() (start int64, end int64) {
	start = BeginningOfDay().Unix()
	end = EndOfDay().Unix()
	return
}

func TimeToString(val int64) string {
	var timeNow = time.Unix(val, 0).In(loc)
	var minute = timeNow.Minute()
	var minuteStr string
	if minute == 0 {
		minuteStr = "00"
	} else {
		minuteStr = strconv.Itoa(minute)
	}
	return strconv.Itoa(timeNow.Hour()) + ":" + minuteStr
}
