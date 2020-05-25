package submerge

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type subLine struct {
	Num   int
	Time  string
	Text1 string
	Text2 string
}

func (s *subLine) isAfter(sub2 *subLine) bool {
	if sub2 == nil {
		return false
	}
	if s == nil {
		return true
	}
	times := []string{s.Time, sub2.Time}
	sort.Strings(times)
	return times[1] == s.Time
}

func (s *subLine) String() string {
	return fmt.Sprintf("[%d] %s   ((  %s |  %s ))", s.Num, s.Time, s.Text1, s.Text2)
}

func (s *subLine) toFormat() string {
	wr := strings.Builder{}
	wr.Write([]byte(fmt.Sprintf("%d\n", s.Num)))
	wr.Write([]byte(fmt.Sprintf("%s\n", s.Time)))
	wr.Write([]byte(fmt.Sprintf("%s\n", s.Text1)))
	if s.Text2 != "" {
		wr.Write([]byte(fmt.Sprintf("%s\n", s.Text2)))
	}
	wr.Write([]byte("\n"))
	return wr.String()
}

func (s *subLine) addColor(color string) {
	if s.Text1 != "" {
		s.Text1 = fmt.Sprintf(`<font color="%s">%s</font>`, color, s.Text1)
	}

	if s.Text2 != "" {
		s.Text2 = fmt.Sprintf(`<font color="%s">%s</font>`, color, s.Text2)
	}
}

func (s *subLine) addDelay(hours, mins, secs, ms int64, delayText bool) {
	// 00:03:35,954 --> 00:03:37,834
	times := strings.Split(s.Time, " --> ")
	if len(times) < 2 {
		fmt.Println(s.Time)
		fmt.Println(s.Num)
		fmt.Println(s.Text1)
		fmt.Println(s.Text2)
	}
	hours1, mins1, secs1, ms1 := s.timeAsInts(times[0])
	hours2, mins2, secs2, ms2 := s.timeAsInts(times[1])

	m1 := ms1 + (secs1 * 1000) + (mins1 * 1000 * 60) + (hours1 * 1000 * 60 * 60)
	m2 := ms2 + (secs2 * 1000) + (mins2 * 1000 * 60) + (hours2 * 1000 * 60 * 60)

	d1 := ms + (secs * 1000) + (mins * 1000 * 60) + (hours * 1000 * 60 * 60)
	d2 := ms + (secs * 1000) + (mins * 1000 * 60) + (hours * 1000 * 60 * 60)

	if !delayText {
		d1 *= -1
		d2 *= -1
	}
	t1 := time.Duration((m1 + d1) * 1000 * 1000)
	t2 := time.Duration((m2 + d2) * 1000 * 1000)
	if t1.Nanoseconds() < 0 {
		t1 = 0
	}
	if t2.Nanoseconds() < 0 {
		t2 = 0
	}
	s.setNewTimes(t1, t2)
}

func (s *subLine) setNewTimes(t1, t2 time.Duration) {

	hour1 := t1 / time.Hour
	t1 -= hour1 * time.Hour
	min1 := t1 / time.Minute
	t1 -= min1 * time.Minute
	sec1 := t1 / time.Second
	t1 -= sec1 * time.Second
	ms1 := t1 / time.Millisecond

	hour2 := t2 / time.Hour
	t2 -= hour2 * time.Hour
	min2 := t2 / time.Minute
	t2 -= min2 * time.Minute
	sec2 := t2 / time.Second
	t2 -= sec2 * time.Second
	ms2 := t2 / time.Millisecond
	t2 -= ms2 * time.Millisecond

	timeFormat := "%02d:%02d:%02d,%03d --> %02d:%02d:%02d,%03d"
	s.Time = fmt.Sprintf(timeFormat, hour1, min1, sec1, ms1, hour2, min2, sec2, ms2)
}

func (s *subLine) timeAsInts(time string) (hours, mins, secs, ms int64) {
	timeItems := strings.Split(time, ",")
	hoursMinsSecs := strings.Split(timeItems[0], ":")
	ms, err := strconv.ParseInt(timeItems[1], 10, 64)
	if err != nil {

		panic(err)
	}
	hours, err = strconv.ParseInt(hoursMinsSecs[0], 10, 64)
	if err != nil {
		panic(err)
	}
	mins, err = strconv.ParseInt(hoursMinsSecs[1], 10, 64)
	if err != nil {
		panic(err)
	}
	secs, err = strconv.ParseInt(hoursMinsSecs[2], 10, 64)
	if err != nil {
		panic(err)
	}
	return
}

func adjustNums(lines []*subLine) {
	for i, line := range lines {
		if line == nil {
			fmt.Printf("\n Missing: %d", i)
			continue
		}
		line.Num = i
	}
}
