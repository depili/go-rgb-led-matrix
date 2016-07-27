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
		fmt.Printf("Event: %s\n", event.Name)
		fmt.Printf("\tStarts: %s\n", event.Start_time.String())
		fmt.Printf("\tEnds: %s\n", event.End_time.String())
		duration := event.End_time.Sub(event.Start_time.Time)
		fmt.Printf("\tDuration: %s\n", duration.String())
		fmt.Printf("\n\n")

	}

	ev := sched.NextEvent(time.Now(), "bigscreen")
	fmt.Printf("Next big screen event: %s", ev.Name)
}
