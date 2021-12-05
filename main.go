package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/jalaali/go-jalaali"
	uuid "github.com/satori/go.uuid"
	"pip2badesabe.ir/utils"
)

type Event struct {
	Event   string `json:"event"`
	Holiday bool   `json:"holiday"`
}

type Events struct {
	Events []Event `json:"events"`
	Date   string  `json:"date"`
	JDate  string  `json:"jDate"`
	JDay   string  `json:"jDay"`
	JMonth string  `json:"jmonth"`
}

const (
	layoutISO          = "2006-01-02"
	layoutMonthName    = "January"
	badesabaAPIAddress = "https://badesaba.ir/api/site/getDataCalendar"
)

func main() {

	// TODO: CLI, calendar.json like time.ir, git, github, refactor nested loop
	// go run, memory pointer

	var getYears = []string{"1400", "1401", "1402", "1403"}

	for _, year := range getYears {
		var years map[string][]Events

		years, err := getEventsOfTheYear(year)
		if err != nil {
			log.Fatal(err)
		}

		icsEvents := createICSString(year, years)
		jsonEvents, _ := json.Marshal(years)

		createFile(year, icsEvents, jsonEvents)
	}

}

func createFile(year string, icsData []byte, jsonData []byte) {
	ioutil.WriteFile(fmt.Sprintf("dist/data-%s.json", year), jsonData, 0644)
	ioutil.WriteFile(fmt.Sprintf("dist/events-%s.ics", year), icsData, 0644)
}

func getEventsOfTheYear(year string) (map[string][]Events, error) {
	months := make(map[string][]Events)
	intYear, _ := strconv.Atoi(year)
	month := 1
	for i := 1; i <= 12; i++ {
		if events, err := getEvents(intYear, month); err != nil {
			return nil, err
		} else if events == nil {
			return nil, fmt.Errorf("event is nil")
		} else {
			month++
			months[strconv.Itoa(i)] = events
		}
	}
	return months, nil
}

func getEvents(year, month int) ([]Events, error) {

	url := fmt.Sprintf("%s/%d/%d", badesabaAPIAddress, month, year)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var Error = fmt.Errorf("httpstatus: %d", resp.StatusCode)
		return nil, Error
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []Events

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}

	for k, _ := range result {
		t, _ := time.Parse(layoutISO, result[k].Date)
		JDate, _ := jalaali.From(t).JFormat(layoutISO)
		Jmonth, _ := jalaali.From(t).JFormat(layoutMonthName)
		result[k].JDate = utils.FaToEn.Replace(JDate)
		result[k].JDay = utils.FaToEn.Replace(JDate)[8:10]
		result[k].JMonth = Jmonth
	}

	return result, nil
}

func createICSString(year string, events map[string][]Events) []byte {

	cal := ics.NewCalendar()
	cal.SetCalscale("GREGORIAN")
	cal.SetProductId("badesamab/ics")
	cal.SetMethod("PUBLISH")
	cal.SetMethod(ics.MethodRequest)

	// TODO: wrong
	for _, v := range events {
		for _, o := range v {
			for _, d := range o.Events {
				var summery string

				if d.Holiday {
					summery = fmt.Sprintf("تعطیل - %s %s %s", o.JDay, o.JMonth, d.Event)
				} else {
					summery = fmt.Sprintf("%s %s %s", o.JDay, o.JMonth, d.Event)
				}

				summery = utils.TrimString(summery)

				timse, _ := time.Parse(layoutISO, o.Date)

				description := fmt.Sprintf("%s %s %s", o.JDay, o.JMonth, d.Event)
				description = utils.TrimString(description)

				event := cal.AddEvent(uuid.NewV4().String())
				event.SetSummary(summery)
				event.SetDtStampTime(timse.Add(4 * time.Hour))
				event.SetAllDayStartAt(timse)
				event.SetDescription(description)
				if d.Holiday {
					event.SetProperty("X-MICROSOFT-CDO-BUSYSTATUS", "FREE")
				}
			}
		}
	}

	return []byte(cal.Serialize())
}
