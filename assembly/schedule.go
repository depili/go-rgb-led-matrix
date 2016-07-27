package assembly

import (
	"encoding/json"
	"fmt"
	"time"
)

type schedule struct {
	Locations  map[string]location `json:"locations"`
	Events     []event             `json:"events"`
	Event_keys map[string]event
}

type location struct {
	Name_fi string `json:"name_fi"`
	Url     string `json:"url"`
	Name    string `json:"name"`
}

type event struct {
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

func ParseSchedule(data []byte) *schedule {
	var sched schedule
	if err := json.Unmarshal(data, &sched); err != nil {
		panic(err)
	}
	sched.Event_keys = make(map[string]event)
	for _, ev := range sched.Events {
		sched.Event_keys[ev.Key] = ev
	}

	return &sched
}

func (sched *schedule) NextEvent(t time.Time, flag string) event {
	var ret event
	for _, ev := range sched.Events {
		if ev.HasFlag(flag) && ev.Start_time.After(t) {
			ret = ev
			break
		}
	}
	if ret.Name == "" {
		panic("no event found with the flag!")
	}

	for _, ev := range sched.Events {
		if ev.HasFlag(flag) {
			if ev.Start_time.After(t) && ev.Start_time.Before(ret.Start_time.Time) {
				ret = ev
			}
		}
	}

	return ret
}

func (ev *event) HasFlag(flag string) bool {
	for _, f := range ev.Flags {
		if f == flag {
			return true
		}
	}
	return false
}

func (ev *event) String() string {
	return fmt.Sprintf("Event: \t%s\n\tStarts: %s\n\tEnds: %s\n\tDuration: %s\n",
		ev.Name, ev.Start_time.String(), ev.End_time.String(),
		ev.End_time.Sub(ev.Start_time.Time).String())
}

func (ev *event) TimeToGo(t time.Time) (string, bool) {
	ttg := ev.Start_time.Sub(t)
	if ev.Start_time.After(t) {
		return fmt.Sprintf("T-%2.f:%02d:%02d",
			ttg.Hours(), int(ttg.Minutes())%60, int(ttg.Seconds())%60), false
	} else {
		return fmt.Sprintf("T+%02.f:%02d:%02d",
			-ttg.Hours(), -int(ttg.Minutes())%60, -int(ttg.Seconds())%60), true
	}
}
