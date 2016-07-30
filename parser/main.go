package main

import (
	"fmt"
	"github.com/depili/go-rgb-led-matrix/assembly"
	"github.com/depili/go-rgb-led-matrix/bdf"
	"github.com/depili/go-rgb-led-matrix/matrix"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"os"
	"os/signal"
	"time"
)

var Options struct {
	SmallFont string `short:"f" long:"smallfont" description:"Font for clock and countdown" default:"fonts/6x10.bdf"`
	Font      string `short:"F" long:"font" description:"Font for event name" default:"fonts/7x13B.bdf"`
	Url       string `short:"u" long:"url" description:"Schedule url" default:"http://schedule.assembly.org/asms16/schedules/events.json"`
	Matrix    string `short:"m" long:"matrix" description:"Matrix to connect to" required:"true"`
}

var parser = flags.NewParser(&Options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		panic(err)
	}

	url := Options.Url
	shutdown := make(chan bool)
	schedChan := make(chan *assembly.Schedule)

	font, err := bdf.Parse(Options.Font)
	if err != nil {
		panic(err)
	}
	smallFont, err := bdf.Parse(Options.SmallFont)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Fonts loaded.\n")

	m := matrix.Init(Options.Matrix, 32, 128)
	defer m.Close()

	m.Fill(matrix.ColorBlack())
	// Start schedule updater go routine
	go assembly.ScheduleWorker(url, schedChan, shutdown)

	var sched *assembly.Schedule
	var ev assembly.Event

	// Colors
	errorColor := [3]byte{0, 255, 255}   // cyan
	pastColor := [3]byte{217, 28, 227}   // pink-ish
	futureColor := [3]byte{28, 227, 190} // Turquise-ish
	fuseColor := [3]byte{74, 35, 17}     // dark brown

	// Initial bitmaps
	errorBitmap := font.TextBitmap("Schedule not imported yet.. waiting...  ")
	evBitmap := font.TextBitmap("Event name")
	ttgBitmap := smallFont.TextBitmap("123")
	clockBitmap := smallFont.TextBitmap("15:04:05")

	// Build intial flame effect bitmaps and palette
	m.InitFlame()

	// Status variables
	ttg := "123"
	inPast := false
	haveSched := false
	haveEvent := false
	minutesToGo := float64(32)
	step := 0

	// Initialize tickers for various tasks
	eventTicker := time.NewTicker(time.Millisecond * 100)
	scrollTicker := time.NewTicker(time.Millisecond * 10)

	// Time to show past events for
	delta, _ := time.ParseDuration("-15m")

	// Lengths for the clocks
	clockX := 127 - (8 * smallFont.Width)
	ttgLength := clockX - 5

	// Trap SIGINT aka Ctrl-C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Main message loop
	for {
		select {
		case sched = <-schedChan: // Schedule update one per 5min
			fmt.Printf("New schedule parsed! Events: %d\n", len(sched.Events))
			haveSched = true
		case <-eventTicker.C: // Check for the current event
			if haveSched {
				if e, found := sched.NextEvent(time.Now().Add(delta), "bigscreen"); found {
					haveEvent = true
					ev = e
					// Generate bitmaps
					evBitmap = font.TextBitmap(fmt.Sprintf("%s  ", ev.Name))
					ttg, inPast = ev.TimeToGo(time.Now())
					ttgBitmap = smallFont.TextBitmap(ttg)
					clockBitmap = smallFont.TextBitmap(time.Now().Format("15:04:05"))

					// Calculate minutes till the event
					minutesToGo = ev.Start_time.Sub(time.Now()).Minutes()
					minutesToGo -= float64(int(minutesToGo)/32) * 32 // Debug
				} else {
					// No event found
					haveEvent = false
				}
			}
		case <-scrollTicker.C: // Advance animations
			if haveSched && haveEvent {
				m.Fill(matrix.ColorBlack())
				// Burning fuse showing last 32min prior to event
				if minutesToGo <= 31.5 && minutesToGo > 0 {
					for i := 0; i < 128; i++ {
						m.FlameClear(29, i)
						m.FlameClear(30, i)
						m.FlameClear(31, i)
					}
					for i := int(((32 - minutesToGo) * 4)); i > int(((32 - minutesToGo - 1) * 4)); i-- {
						m.FlameSet(29, i)
						m.FlameSet(30, i)
						m.FlameSet(31, i)
					}
					m.FlameFill()
				}

				for i := 127; i > int(((32 - minutesToGo) * 4)); i-- {
					m.SetPixel(29, i, fuseColor)
					m.SetPixel(30, i, fuseColor)
					m.SetPixel(31, i, fuseColor)
				}

				// Scroll the event name
				m.Scroll(evBitmap, matrix.ColorWhite(), 0, 0, step/2, 128)

				// T±12:04:05 output, set color based on ±
				if inPast {
					m.Scroll(ttgBitmap, pastColor, 16, 0, 0, ttgLength)
				} else {
					m.Scroll(ttgBitmap, futureColor, 16, 0, 0, ttgLength)
				}

				// Clock
				m.ScrollPlasma(clockBitmap, 16, clockX, step/5, 56)
			} else {
				// Schedule not loaded yet
				m.Fill(errorColor)
				m.Scroll(errorBitmap, matrix.ColorWhite(), 5, 0, step, 128)
			}
			m.Send()
			step++

		case <-sigChan:
			// SIGINT received, shutdown gracefully
			m.Close()
			shutdown <- true
			os.Exit(1)
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
