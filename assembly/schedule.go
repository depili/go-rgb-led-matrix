package assembly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

type Schedule struct {
	Locations  map[string]location `json:"locations"`
	Events     events              `json:"events"`
	Event_keys map[string]Event
}

type location struct {
	Name_fi string `json:"name_fi"`
	Url     string `json:"url"`
	Name    string `json:"name"`
}

type Event struct {
	Location_key        string       `json:"locations_key"`
	Name                string       `json:"name"`
	Start_time          AssemblyTime `json:"start_time"`
	Original_start_time AssemblyTime `json:"original_start_time"`
	Flags               []string     `json:"flags"`
	End_time            AssemblyTime `json:"end_time"`
	Key                 string       `json:"key"`
	Name_fi             string       `json:"name_fi"`
	Categories          []string     `json:"categories"`
}

type events []Event

func (slice events) Len() int {
	return len(slice)
}

func (slice events) Less(i, j int) bool {
	return slice[i].Start_time.Before(slice[j].Start_time.Time)
}

func (slice events) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func ParseSchedule(data []byte) *Schedule {
	var sched Schedule
	if err := json.Unmarshal(data, &sched); err != nil {
		panic(err)
	}
	sched.Event_keys = make(map[string]Event)
	for _, ev := range sched.Events {
		sched.Event_keys[ev.Key] = ev
	}
	sort.Sort(sched.Events)

	return &sched
}

func (sched *Schedule) NextEvent(t time.Time, flag string) (Event, bool) {
	for _, ev := range sched.Events {
		if ev.HasFlag(flag) && ev.Start_time.After(t) {
			return ev, true
		}
	}
	var ret Event
	return ret, false
}

func (ev *Event) HasFlag(flag string) bool {
	for _, f := range ev.Flags {
		if f == flag {
			return true
		}
	}
	return false
}

func (ev *Event) Equal(otherEvent Event) bool {
	return ev.Key == otherEvent.Key
}

func (ev *Event) String() string {
	return fmt.Sprintf("Event: \t%s\n\tStarts: %s\n\tEnds: %s\n\tDuration: %s\n",
		ev.Name, ev.Start_time.String(), ev.End_time.String(),
		ev.End_time.Sub(ev.Start_time.Time).String())
}

func (ev *Event) TimeToGo(t time.Time) (string, bool) {
	ttg := ev.Start_time.Sub(t)
	if ev.Start_time.After(t) {
		return fmt.Sprintf("-%02d:%02d:%02d",
			int(ttg.Hours()), int(ttg.Minutes())%60, int((ttg.Seconds())+1)%60), false
	} else {
		return fmt.Sprintf("+%02d:%02d:%02d",
			-int(ttg.Hours()), -int(ttg.Minutes())%60, -int((ttg.Seconds())+1)%60), true
	}
}

func ScheduleWorker(url string, schedChan chan *Schedule, shutdown chan bool) {
	fetchSchedule(url, schedChan)
	updateTicker := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-shutdown:
			return
		case <-updateTicker.C:
			fetchSchedule(url, schedChan)
		}
	}
}

func fetchSchedule(url string, schedChan chan *Schedule) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching schedule data from url: %s", url)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		sched := ParseSchedule(body)
		schedChan <- sched
	}
	resp.Body.Close()
}
