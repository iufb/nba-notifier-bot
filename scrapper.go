package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

var zeroTime time.Time

type ScheduleSc struct {
	team string
	ot   string
	date time.Time
}

func TimeParser(timeStr string) time.Time {
	layout := "2006/1/2 3:04 pm"
	d, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Println("Parser error : ", err)
		return zeroTime
	}
	return d
}

func Scrapper(url string) ScheduleSc {
	c := colly.NewCollector()
	s := []ScheduleSc{}

	c.OnHTML("div.Schedule__Game__Wrapper", func(e *colly.HTMLElement) {
		t := ScheduleSc{}
		e.ForEach("div.Schedule__Info", func(i int, e *colly.HTMLElement) {
			t.ot = e.ChildText("span.Schedule__Team")
		})
		e.ForEach("div.Schedule__Meta", func(i int, e *colly.HTMLElement) {
			var date string
			e.ForEach("span.Schedule__Time", func(i int, e *colly.HTMLElement) {
				if i == 0 {
					date += e.Text + " "
					return
				}
				date += e.Text
			})
			if date == "" {
				t.date = zeroTime
				return
			}
			date = fmt.Sprintf("%d/%s", time.Now().Year(), date)
			t.date = TimeParser(date)
		})
		if t.date == zeroTime {
			return
		}

		s = append(s, t)
	})
	c.Visit(url)
	return s[len(s)-1]
}
