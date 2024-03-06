package main

// *****************************************************************************
//  __________      .__  .__       .__
//  \______   \____ |  | |__| _____|  |__
//   |     ___/  _ \|  | |  |/  ___/  |  \
//   |    |  (  <_> )  |_|  |\___ \|   Y  \
//   |____|   \____/|____/__/____  >___|  /
//                               \/     \/
//
//          Polish © jpl@ozf.fr 2024
//
// *****************************************************************************

// *****************************************************************************
// IMPORTS
// *****************************************************************************
import (
	"math"
	"math/rand"
	"time"
)

// *****************************************************************************
// CONSTANTS
// *****************************************************************************
const (
	deg2rad = math.Pi / 180.0
	rad2deg = 180.0 / math.Pi
)

// *****************************************************************************
// TYPES
// *****************************************************************************
type My struct{}

// *****************************************************************************
// MyAdd()
// *****************************************************************************
func (m My) MyAdd() {
	if checkStack(2) {
		f2, _ := fs.Pop()
		f1, _ := fs.Pop()
		fs.Push(f1 + f2)
	}
}

// *****************************************************************************
// MySub()
// *****************************************************************************
func (m My) MySub() {
	if checkStack(2) {
		f2, _ := fs.Pop()
		f1, _ := fs.Pop()
		fs.Push(f1 - f2)
	}
}

// *****************************************************************************
// MyMult()
// *****************************************************************************
func (m My) MyMult() {
	if checkStack(2) {
		f2, _ := fs.Pop()
		f1, _ := fs.Pop()
		fs.Push(f1 * f2)
	}
}

// *****************************************************************************
// MyDiv()
// *****************************************************************************
func (m My) MyDiv() {
	if checkStack(2) {
		f2, _ := fs.Pop()
		f1, _ := fs.Pop()
		fs.Push(f1 / f2)
	}
}

// *****************************************************************************
// MyAbs()
// *****************************************************************************
func (m My) MyAbs() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Abs(f))
	}
}

// *****************************************************************************
// MyAcos()
// *****************************************************************************
func (m My) MyAcos() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Acos(f))
	}
}

// *****************************************************************************
// MyAcosh()
// *****************************************************************************
func (m My) MyAcosh() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Acosh(f))
	}
}

// *****************************************************************************
// MyAsin()
// *****************************************************************************
func (m My) MyAsin() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Asin(f))
	}
}

// *****************************************************************************
// MyAsinh()
// *****************************************************************************
func (m My) MyAsinh() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Asinh(f))
	}
}

// *****************************************************************************
// MyAtan()
// *****************************************************************************
func (m My) MyAtan() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Atan(f))
	}
}

// *****************************************************************************
// MyAtan2()
// *****************************************************************************
func (m My) MyAtan2() {
	if checkStack(2) {
		f2, _ := fs.Pop()
		f1, _ := fs.Pop()
		fs.Push(math.Atan2(f2, f1))
	}
}

// *****************************************************************************
// MyAtanh()
// *****************************************************************************
func (m My) MyAtanh() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Atanh(f))
	}
}

// *****************************************************************************
// MyCbrt()
// *****************************************************************************
func (m My) MyCbrt() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Cbrt(f))
	}
}

// *****************************************************************************
// MyCeil()
// *****************************************************************************
func (m My) MyCeil() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Ceil(f))
	}
}

// *****************************************************************************
// MyCos()
// *****************************************************************************
func (m My) MyCos() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Cos(f))
	}
}

// *****************************************************************************
// MyCosh()
// *****************************************************************************
func (m My) MyCosh() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Cosh(f))
	}
}

// *****************************************************************************
// MyExp()
// *****************************************************************************
func (m My) MyExp() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Exp(f))
	}
}

// *****************************************************************************
// MyExp2()
// *****************************************************************************
func (m My) MyExp2() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Exp2(f))
	}
}

// *****************************************************************************
// MyFloor()
// *****************************************************************************
func (m My) MyFloor() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Floor(f))
	}
}

// *****************************************************************************
// MyGamma()
// *****************************************************************************
func (m My) MyGamma() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Gamma(f))
	}
}

// *****************************************************************************
// MyHypot()
// *****************************************************************************
func (m My) MyHypot() {
	if checkStack(2) {
		f2, _ := fs.Pop()
		f1, _ := fs.Pop()
		fs.Push(math.Hypot(f1, f2))
	}
}

// *****************************************************************************
// MyInv()
// *****************************************************************************
func (m My) MyInv() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(1 / f)
	}
}

// *****************************************************************************
// MyLog()
// *****************************************************************************
func (m My) MyLog() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Log(f))
	}
}

// *****************************************************************************
// MyLog10()
// *****************************************************************************
func (m My) MyLog10() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Log10(f))
	}
}

// *****************************************************************************
// MyRnd()
// *****************************************************************************
func (m My) MyRnd() {
	fs.Push(rand.Float64())
}

// *****************************************************************************
// MyRound()
// *****************************************************************************
func (m My) MyRound() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Round(f))
	}
}

// *****************************************************************************
// MySin()
// *****************************************************************************
func (m My) MySin() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Sin(f))
	}
}

// *****************************************************************************
// MySinh()
// *****************************************************************************
func (m My) MySinh() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Sinh(f))
	}
}

// *****************************************************************************
// MySqrt()
// *****************************************************************************
func (m My) MySqrt() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Sqrt(f))
	}
}

// *****************************************************************************
// MyTan()
// *****************************************************************************
func (m My) MyTan() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Tan(f))
	}
}

// *****************************************************************************
// MyTanh()
// *****************************************************************************
func (m My) MyTanh() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Tanh(f))
	}
}

// *****************************************************************************
// MyTrunc()
// *****************************************************************************
func (m My) MyTrunc() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Trunc(f))
	}
}

// *****************************************************************************
// MyFrac()
// *****************************************************************************
func (m My) MyFrac() {
	if checkStack(1) {
		f, _ := fs.Pop()
		_, frac := math.Modf(f)
		fs.Push(frac)
	}
}

// *****************************************************************************
// MyPi()
// *****************************************************************************
func (m My) MyPi() {
	fs.Push(math.Pi)
}

// *****************************************************************************
// MyE()
// *****************************************************************************
func (m My) MyE() {
	fs.Push(math.E)
}

// *****************************************************************************
// MyPhi()
// *****************************************************************************
func (m My) MyPhi() {
	fs.Push(math.Phi)
}

// *****************************************************************************
// MyPow()
// *****************************************************************************
func (m My) MyPow() {
	if checkStack(2) {
		f2, _ := fs.Pop()
		f1, _ := fs.Pop()
		fs.Push(math.Pow(f1, f2))
	}
}

// *****************************************************************************
// MySqr()
// *****************************************************************************
func (m My) MySqr() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(math.Pow(f, 2))
	}
}

// *****************************************************************************
// MyRoot()
// *****************************************************************************
func (m My) MyRoot() {
	if checkStack(2) {
		f2, _ := fs.Pop()
		f1, _ := fs.Pop()
		fs.Push(math.Pow(f1, 1/f2))
	}
}

// *****************************************************************************
// factorial()
// *****************************************************************************
func factorial(number float64) float64 {

	// if the number has reached 1 then we have to
	// return 1 as 1 is the minimum value we have to multiply with
	if number == 1 {
		return 1
	}

	// multiplying with the current number and calling the function
	// for 1 lesser number
	factorialOfNumber := number * factorial(number-1)

	// return the factorial of the current number
	return factorialOfNumber
}

// *****************************************************************************
// MyFact()
// *****************************************************************************
func (m My) MyFact() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(factorial(f))
	}
}

// *****************************************************************************
// MyTorad()
// *****************************************************************************
func (m My) MyTorad() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(f * deg2rad)
	}
}

// *****************************************************************************
// MyTodeg()
// *****************************************************************************
func (m My) MyTodeg() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(f * rad2deg)
	}
}

// *****************************************************************************
// MyTodms()
// *****************************************************************************
func (m My) MyTodms() {
	if checkStack(1) {
		f, _ := fs.Pop()
		sign := 1.0
		if f < 0 {
			sign = -1.0
		}
		f = math.Abs(f)
		d := math.Floor(f)
		m1 := (f - d) * 60.0
		mm := math.Floor(m1)
		ss := (m1 - mm) * 60.0

		dms := d + (mm / 100.0) + (ss / 10000.0)
		dms = sign * dms
		fs.Push(dms)
	}
}

// *****************************************************************************
// MyTodec()
// *****************************************************************************
func (m My) MyTodec() {
	if checkStack(1) {
		f, _ := fs.Pop()
		sign := 1.0
		if f < 0 {
			sign = -1.0
		}
		f = math.Abs(f)
		d := math.Floor(f)

		mm := math.Floor((f - d) * 100.0)

		ss := (((f - d) * 100.0) - mm) * 100.0

		dec := d + (mm / 60.0) + (ss / 3600.0)
		dec = sign * dec
		fs.Push(dec)
	}
}

// *****************************************************************************
// MyNow()
// *****************************************************************************
func (m My) MyNow() {
	now := time.Now()
	yyyy := float64(now.Year())
	mm := float64(int(now.Month()))
	dd := float64(now.Day())
	hh := float64(now.Hour())
	nn := float64(now.Minute())
	ss := float64(now.Second())

	dh := yyyy*10000.0 + mm*100.0 + dd
	dh += hh/100.0 + nn/10000.0 + ss/1000000.0

	fs.Push(dh)
}

// *****************************************************************************
// getDayFromDate()
// *****************************************************************************
func getDayFromDate(dh float64) float64 {
	_, dd := math.Modf(math.Floor(dh) / 100.0)
	dd *= 100.0
	return dd
}

// *****************************************************************************
// MyDay()
// *****************************************************************************
func (m My) MyDay() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(f)
		fs.Push(getDayFromDate(f))
	}
}

// *****************************************************************************
// getMonthFromDate()
// *****************************************************************************
func getMonthFromDate(dh float64) float64 {
	_, mm := math.Modf(math.Floor(dh) / 10000.0)
	mm *= 100.0
	mm = math.Floor(mm)
	return mm
}

// *****************************************************************************
// MyMonth()
// *****************************************************************************
func (m My) MyMonth() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(f)
		fs.Push(getMonthFromDate(f))
	}
}

// *****************************************************************************
// getYearFromDate()
// *****************************************************************************
func getYearFromDate(dh float64) float64 {
	yy := math.Floor(math.Floor(dh) / 10000.0)
	return yy
}

// *****************************************************************************
// MyYear()
// *****************************************************************************
func (m My) MyYear() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(f)
		fs.Push(getYearFromDate(f))
	}
}

// *****************************************************************************
// getHourFromDate()
// *****************************************************************************
func getHourFromDate(dh float64) float64 {
	_, hh := math.Modf(dh)
	hh *= 100.0
	hh = math.Floor(hh)
	return hh
}

// *****************************************************************************
// MyHour()
// *****************************************************************************
func (m My) MyHour() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(f)
		fs.Push(getHourFromDate(f))
	}
}

// *****************************************************************************
// getMinuteFromDate()
// *****************************************************************************
func getMinuteFromDate(dh float64) float64 {
	_, mm := math.Modf(dh)
	mm *= 100.0
	_, mm = math.Modf(mm)
	mm *= 100.0
	mm = math.Floor(mm)
	return mm
}

// *****************************************************************************
// MyMinute()
// *****************************************************************************
func (m My) MyMinute() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(f)
		fs.Push(getMinuteFromDate(f))
	}
}

// *****************************************************************************
// getSecondFromDate()
// *****************************************************************************
func getSecondFromDate(dh float64) float64 {
	_, ss := math.Modf(dh)
	ss *= 10000.0
	_, ss = math.Modf(ss)
	ss *= 100.0
	ss = math.Floor(ss)
	return ss
}

// *****************************************************************************
// MySecond()
// *****************************************************************************
func (m My) MySecond() {
	if checkStack(1) {
		f, _ := fs.Pop()
		fs.Push(f)
		fs.Push(getSecondFromDate(f))
	}
}

// *****************************************************************************
// MyDate()
// *****************************************************************************
func (m My) MyDate() {
	if checkStack(6) {
		ss, _ := fs.Pop()
		nn, _ := fs.Pop()
		hh, _ := fs.Pop()
		dd, _ := fs.Pop()
		mm, _ := fs.Pop()
		yyyy, _ := fs.Pop()
		fs.Push(yyyy*10000.0 + mm*100.0 + dd + hh/100.0 + nn/10000.0 + ss/1000000.0)
	}
}

// *****************************************************************************
// leapYears()
// *****************************************************************************
func leapYears(date time.Time) (leaps int) {
	// https://www.tutorialspoint.com/golang-program-to-calculate-difference-between-two-time-periods
	y, m, _ := date.Date()
	if m <= 2 {
		y--
	}
	leaps = y/4 + y/400 - y/100
	return leaps
}

// *****************************************************************************
// getDifference()
// *****************************************************************************
func getDifference(a, b time.Time) (days, hours, minutes, seconds int) {
	// https://www.tutorialspoint.com/golang-program-to-calculate-difference-between-two-time-periods
	monthDays := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	h1, min1, s1 := a.Clock()
	h2, min2, s2 := b.Clock()
	totalDays1 := y1*365 + d1
	for i := 0; i < (int)(m1)-1; i++ {
		totalDays1 += monthDays[i]
	}
	totalDays1 += leapYears(a)
	totalDays2 := y2*365 + d2
	for i := 0; i < (int)(m2)-1; i++ {
		totalDays2 += monthDays[i]
	}
	totalDays2 += leapYears(b)
	days = totalDays2 - totalDays1
	hours = h2 - h1
	minutes = min2 - min1
	seconds = s2 - s1
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	return days, hours, minutes, seconds
}

// *****************************************************************************
// MyDdiff()
// *****************************************************************************
func (m My) MyDdiff() {
	if checkStack(2) {
		d1, _ := fs.Pop()
		d2, _ := fs.Pop()

		date1 := time.Date(int(getYearFromDate(d1)), time.Month(getMonthFromDate(d1)), int(getDayFromDate(d1)), int(getHourFromDate(d1)), int(getMinuteFromDate(d1)), int(getSecondFromDate(d1)), 0, time.UTC)
		date2 := time.Date(int(getYearFromDate(d2)), time.Month(getMonthFromDate(d2)), int(getDayFromDate(d2)), int(getHourFromDate(d2)), int(getMinuteFromDate(d2)), int(getSecondFromDate(d2)), 0, time.UTC)
		if date1.After(date2) {
			date1, date2 = date2, date1
		}
		dd, hh, nn, ss := getDifference(date1, date2)
		fs.Push(float64(ss))
		fs.Push(float64(nn))
		fs.Push(float64(hh))
		fs.Push(float64(dd))
	}
}

/*
date 	( yyyy mm dd hh nn ss -- d.h )
now		( -- d.h )						20240327.202346
day		( d.h -- d.h dd )				27
month	( d.h -- d.h mm )				3
year	( d.h -- d.h yyyy )				2024
hour	( d.h -- d.h hh )				20
minute	( d.h -- d.h nn )				23
second	( d.h -- d.h ss )				46

+day	( d.h n -- d.h )
+month	( d.h n -- d.h )
+year	( d.h n -- d.h )
+hour	( d.h n -- d.h )
+minute	( d.h n -- d.h )
+second	( d.h n -- d.h )

-day	( d.h n -- d.h )
-month	( d.h n -- d.h )
-year	( d.h n -- d.h )
-hour	( d.h n -- d.h )
-minute	( d.h n -- d.h )
-second	( d.h n -- d.h )

ddiff	( d.h d.h -- d h m s )

https://www.tutorialspoint.com/golang-program-to-calculate-difference-between-two-time-periods
func leapYears(date time.Time) (leaps int) {
   y, m, _ := date.Date()
   if m <= 2 {
      y--
   }
   leaps = y/4 + y/400 - y/100
   return leaps
}
func getDifference(a, b time.Time) (days, hours, minutes, seconds int) {
   monthDays := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
   y1, m1, d1 := a.Date()
   y2, m2, d2 := b.Date()
   h1, min1, s1 := a.Clock()
   h2, min2, s2 := b.Clock()
   totalDays1 := y1*365 + d1
   for i := 0; i < (int)(m1)-1; i++ {
      totalDays1 += monthDays[i]
   }
   totalDays1 += leapYears(a)
   totalDays2 := y2*365 + d2
   for i := 0; i < (int)(m2)-1; i++ {
      totalDays2 += monthDays[i]
   }
   totalDays2 += leapYears(b)
   days = totalDays2 - totalDays1
   hours = h2 - h1
   minutes = min2 - min1
   seconds = s2 - s1
   if seconds < 0 {
      seconds += 60
      minutes--
   }
   if minutes < 0 {
      minutes += 60
      hours--
   }
   if hours < 0 {
      hours += 24
      days--
   }
   return days, hours, minutes, seconds
}
func main() {
   date1 := time.Date(2020, 4, 27, 23, 35, 0, 0, time.UTC)
   date2 := time.Date(2018, 5, 12, 12, 43, 23, 0, time.UTC)
   if date1.After(date2) {
      date1, date2 = date2, date1
   }
   days, hours, minutes, seconds := getDifference(date1, date2)
   fmt.Println("The difference between dates", date1, "and", date2, "is: ")
   fmt.Printf("%v days\n%v hours\n%v minutes\n%v seconds", days, hours, minutes, seconds)
}

*/
