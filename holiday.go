// (c) 2014 Rick Arnold. Licensed under the BSD license (see LICENSE).

package cal

import (
	"time"
)

// ObservedRule represents a rule for observing a holiday that falls
// on a weekend.
type ObservedRule int

//ObservedRule are the specific ObservedRules
const (
	ObservedNearest ObservedRule = iota // nearest weekday (Friday or Monday)
	ObservedExact                       // the exact day only
	ObservedMonday                      // Monday always
)

var (
	// United States holidays
	US_NewYear      = NewHoliday(time.January, 1)
	US_MLK          = NewHolidayFloat(time.January, time.Monday, 3)
	US_Presidents   = NewHolidayFloat(time.February, time.Monday, 3)
	US_Memorial     = NewHolidayFloat(time.May, time.Monday, -1)
	US_Independence = NewHoliday(time.July, 4)
	US_Labor        = NewHolidayFloat(time.September, time.Monday, 1)
	US_Columbus     = NewHolidayFloat(time.October, time.Monday, 2)
	US_Veterans     = NewHoliday(time.November, 11)
	US_Thanksgiving = NewHolidayFloat(time.November, time.Thursday, 4)
	US_Christmas    = NewHoliday(time.December, 25)

	// Target2 holidays
	ECB_GoodFriday       = NewHolidayFunc(calculateGoodFriday)
	ECB_EasterMonday     = NewHolidayFunc(calculateEasterMonday)
	ECB_NewYearsDay      = NewHoliday(time.January, 1)
	ECB_LabourDay        = NewHoliday(time.May, 1)
	ECB_ChristmasDay     = NewHoliday(time.December, 25)
	ECB_ChristmasHoliday = NewHoliday(time.December, 26)

	// Holidays in Germany
	DE_Neujahr                = US_NewYear
	DE_KarFreitag             = NewHolidayFunc(calculateGoodFriday)
	DE_Ostermontag            = NewHolidayFunc(calculateEasterMonday)
	DE_TagderArbeit           = NewHoliday(time.May, 1)
	DE_Himmelfahrt            = NewHolidayFunc(calculateHimmelfahrt)
	DE_Pfingstmontag          = NewHolidayFunc(calculatePfingstMontag)
	DE_TagderDeutschenEinheit = NewHoliday(time.October, 3)
	DE_ErsterWeihnachtstag    = ECB_ChristmasDay
	DE_ZweiterWeihnachtstag   = ECB_ChristmasHoliday

	// Holidays in the Netherlands
	NLNieuwjaar       = US_NewYear
	NLGoedeVrijdag    = ECB_GoodFriday
	NLPaasMaandag     = ECB_EasterMonday
	NLKoningsDag      = NewHolidayFunc(calculateKoningsDag)
	NLBevrijdingsDag  = NewHoliday(time.May, 5)
	NLHemelvaart      = DE_Himmelfahrt
	NLPinksterMaandag = DE_Pfingstmontag
	NLEersteKerstdag  = ECB_ChristmasDay
	NLTweedeKerstdag  = ECB_ChristmasHoliday

	// Holidays in Great Britain
	GB_NewYear       = NewHolidayFunc(calculateNewYearsHoliday)
	GB_GoodFriday    = ECB_GoodFriday
	GB_EasterMonday  = ECB_EasterMonday
	GB_EarlyMay      = NewHolidayFloat(time.May, time.Monday, 1)
	GB_SpringHoliday = NewHolidayFloat(time.May, time.Monday, -1)
	GB_SummerHoliday = NewHolidayFloat(time.August, time.Monday, -1)
	GB_ChristmasDay  = ECB_ChristmasDay
	GB_BoxingDay     = ECB_ChristmasHoliday
)

// HolidayFn calculates the occurrence of a holiday for the given year.
// This is useful for holidays like Easter that depend on complex rules.
type HolidayFn func(year int, loc *time.Location) (month time.Month, day int)

// Holiday holds information about the yearly occurrence of a holiday.
//
// A valid Holiday consists of one of the following:
// - Month and Day (such as March 14 for Pi Day)
// - Month, Weekday, and Offset (such as the second Monday of October for Columbus Day)
// - Offset (such as the 183rd day of the year for the start of the second half)
// - Func (to calculate the holiday)
type Holiday struct {
	Month   time.Month
	Weekday time.Weekday
	Day     int
	Offset  int
	Func    HolidayFn

	// last values used to calculate month and day with Func
	lastYear int
	lastLoc  *time.Location
}

func calculateGoodFriday(year int, loc *time.Location) (time.Month, int) {
	easter := calculateEaster(year, loc)
	//Go the the day before yesterday
	gf := easter.AddDate(0, 0, -2)
	return gf.Month(), gf.Day()
}

func calculateEasterMonday(year int, loc *time.Location) (time.Month, int) {
	easter := calculateEaster(year, loc)
	//Go the the day after Easter
	em := easter.AddDate(0, 0, +1)
	return em.Month(), em.Day()
}

func calculateEaster(year int, loc *time.Location) time.Time {
	// Meeus/Jones/Butcher algorithm
	y := year
	a := y % 19
	b := y / 100
	c := y % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451

	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

func calculateHimmelfahrt(year int, loc *time.Location) (time.Month, int) {
	easter := calculateEaster(year, loc)
	//Go the the day after Easter
	em := easter.AddDate(0, 0, +39)
	return em.Month(), em.Day()
}

func calculatePfingstMontag(year int, loc *time.Location) (time.Month, int) {
	easter := calculateEaster(year, loc)
	//Go the the day after Easter
	em := easter.AddDate(0, 0, +50)
	return em.Month(), em.Day()
}

//KoningsDag (kingsday) is April 27th, 26th if the 27th is a Sunday
func calculateKoningsDag(year int, loc *time.Location) (time.Month, int) {
	koningsDag := time.Date(year, time.April, 27, 0, 0, 0, 0, loc)
	if koningsDag.Weekday() == time.Sunday {
		koningsDag = koningsDag.AddDate(0, 0, -1)
	}
	return koningsDag.Month(), koningsDag.Day()
}

// NewYearsDay is the 1st of January unless the 1st is a Saturday or Sunday in which case it occurs on the following Monday.
func calculateNewYearsHoliday(year int, loc *time.Location) (time.Month, int) {
	day := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
	switch day.Weekday() {
	case time.Saturday:
		day = day.AddDate(0, 0, 2)
	case time.Sunday:
		day = day.AddDate(0, 0, 1)
	}
	return time.January, day.Day()
}

// NewHoliday creates a new Holiday instance for an exact day of a month.
func NewHoliday(month time.Month, day int) Holiday {
	return Holiday{Month: month, Day: day}
}

// NewHolidayFloat creates a new Holiday instance for an offset-based day of
// a month.
func NewHolidayFloat(month time.Month, weekday time.Weekday, offset int) Holiday {
	return Holiday{Month: month, Weekday: weekday, Offset: offset}
}

// NewHolidayFunc creates a new Holiday instance that uses a function to
// calculate the day and month.
func NewHolidayFunc(fn HolidayFn) Holiday {
	return Holiday{Func: fn}
}

// matches determines whether the given date is the one referred to by the
// Holiday.
func (h *Holiday) matches(date time.Time) bool {

	if h.Func != nil && (date.Year() != h.lastYear || date.Location() != h.lastLoc) {
		h.Month, h.Day = h.Func(date.Year(), date.Location())
		h.lastYear = date.Year()
		h.lastLoc = date.Location()
	}

	if h.Month > 0 {
		if date.Month() != h.Month {
			return false
		}
		if h.Day > 0 {
			return date.Day() == h.Day
		}
		if h.Weekday > 0 && h.Offset != 0 {
			return IsWeekdayN(date, h.Weekday, h.Offset)
		}
	} else if h.Offset > 0 {
		return date.YearDay() == h.Offset
	}
	return false
}

//AddGermanHolidays adds all German Holdays to Calendar
func AddGermanHolidays(c *Calendar) {
	c.AddHoliday(DE_Neujahr)
	c.AddHoliday(DE_KarFreitag)
	c.AddHoliday(DE_Ostermontag)
	c.AddHoliday(DE_TagderArbeit)
	c.AddHoliday(DE_Himmelfahrt)
	c.AddHoliday(DE_Pfingstmontag)
	c.AddHoliday(DE_TagderDeutschenEinheit)
	c.AddHoliday(DE_ErsterWeihnachtstag)
	c.AddHoliday(DE_ZweiterWeihnachtstag)
}

//AddDutchHolidays adds all Dutch Holdays to Calendar
func AddDutchHolidays(c *Calendar) {
	c.AddHoliday(NLNieuwjaar)
	c.AddHoliday(NLGoedeVrijdag)
	c.AddHoliday(NLPaasMaandag)
	c.AddHoliday(NLKoningsDag)
	c.AddHoliday(NLBevrijdingsDag)
	c.AddHoliday(NLHemelvaart)
	c.AddHoliday(NLPinksterMaandag)
	c.AddHoliday(NLEersteKerstdag)
	c.AddHoliday(NLTweedeKerstdag)
}

// AddBritishHolidays add all British holidays to Calender
func AddBritishHolidays(c *Calendar) {
	c.AddHoliday(GB_NewYear)
	c.AddHoliday(GB_GoodFriday)
	c.AddHoliday(GB_EasterMonday)
	c.AddHoliday(GB_EarlyMay)
	c.AddHoliday(GB_SpringHoliday)
	c.AddHoliday(GB_SummerHoliday)
	c.AddHoliday(GB_ChristmasDay)
	c.AddHoliday(GB_BoxingDay)
}
