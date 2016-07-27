package main

import (
	"fmt"
	"github.com/depili/go-rgb-led-matrix/assembly"
	"io/ioutil"
	"time"
)

func main() {
	data, err := ioutil.ReadFile("json.sample")
	if err != nil {
		panic(err)
	}
	sched := assembly.ParseSchedule(data)
	fmt.Printf("Parsed! Events: %d", len(sched.Events))
	for _, event := range sched.Events {
		fmt.Printf("%s\n\n", event.String())
	}

	stime := "2016-08-04T16:02+0300"
	timeLayout := "2006-01-02T15:04-0700"
	timeDate, _ := time.Parse(timeLayout, stime)

	delta, _ := time.ParseDuration("-5m")
	ev := sched.NextEvent(timeDate.Add(delta), "bigscreen")
	fmt.Printf("Next big screen event:\n%s", ev.String())
	fmt.Printf("time.Now(): %s\n", time.Now())
	ttg, inPast := ev.TimeToGo(timeDate)
	fmt.Printf("%s in past: %t\n", ttg, inPast)
}
