package main

import (
	"fmt"
	"github.com/depili/go-rgb-led-matrix/assembly"
	"github.com/depili/go-rgb-led-matrix/bdf"
	"github.com/depili/go-rgb-led-matrix/matrix"
	"io/ioutil"
	"time"
)

func main() {

	url := "http://schedule.assembly.org/asms16/schedules/events.json"
	shutdown := make(chan bool)
	schedChan := make(chan *assembly.Schedule)

	font, err := bdf.Parse("fonts/7x13B.bdf")
	if err != nil {
		panic(err)
	}
	smallFont, err := bdf.Parse("fonts/6x10.bdf")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Fonts loaded.\n")

	m := matrix.Init("tcp://192.168.0.30:5555", 32, 128)
	defer m.Close()

	m.Fill(matrix.ColorBlack())

	go assembly.ScheduleWorker(url, schedChan, shutdown)
	var sched *assembly.Schedule
	var ev assembly.Event
	errorColor := [3]byte{0, 255, 255}   // cyan
	pastColor := [3]byte{217, 28, 227}   // pink-ish
	futureColor := [3]byte{28, 227, 190} // Turquise-ish
	errorBitmap := font.TextBitmap("Schedule not imported yet.. waiting...  ")
	evBitmap := font.TextBitmap("Event name")
	ttgBitmap := smallFont.TextBitmap("123")
	ttg := "123"
	clockBitmap := smallFont.TextBitmap("15:04:05")
	inPast := false
	haveSched := false
	haveEvent := false
	i := 0
	eventTicker := time.NewTicker(time.Millisecond * 100)
	scrollTicker := time.NewTicker(time.Millisecond * 10)
	delta, _ := time.ParseDuration("-15m")
	clockX := 127 - (8 * smallFont.Width)
	ttgLength := clockX - 5

	for {
		select {
		case sched = <-schedChan:
			fmt.Printf("New schedule parsed! Events: %d\n", len(sched.Events))
			haveSched = true
		case <-eventTicker.C:
			if haveSched {
				if e, found := sched.NextEvent(time.Now().Add(delta), "bigscreen"); found {
					haveEvent = true
					ev = e
					evBitmap = font.TextBitmap(fmt.Sprintf("%s  ", ev.Name))
					ttg, inPast = ev.TimeToGo(time.Now())
					ttgBitmap = smallFont.TextBitmap(ttg)
				} else {
					haveEvent = false
				}
				clockBitmap = smallFont.TextBitmap(time.Now().Format("15:04:05"))
			}
		case <-scrollTicker.C:
			if haveSched && haveEvent {
				m.Fill(matrix.ColorBlack())
				m.Scroll(evBitmap, matrix.ColorWhite(), 0, 0, i/2, 128)
				if inPast {
					m.Scroll(ttgBitmap, pastColor, 21, 0, 0, ttgLength)
				} else {
					m.Scroll(ttgBitmap, futureColor, 21, 0, 0, ttgLength)
				}
				m.ScrollPlasma(clockBitmap, 21, clockX, i/5, 56)

			} else {
				// Schedule not loaded yet
				m.Fill(errorColor)
				m.Scroll(errorBitmap, matrix.ColorWhite(), 5, 0, i, 128)
			}
			m.Send()
			i++
		}
	}
}

func ParseSampleFile() *assembly.Schedule {
	data, err := ioutil.ReadFile("json.sample")
	if err != nil {
		panic(err)
	}

	return assembly.ParseSchedule(data)
}
