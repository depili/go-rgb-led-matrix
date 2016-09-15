package main

import (
	"fmt"
	"github.com/depili/go-rgb-led-matrix/bdf"
	"github.com/depili/go-rgb-led-matrix/matrix"
	"github.com/jessevdk/go-flags"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
	"github.com/tarm/serial"
	"os"
	"os/signal"
	"time"
)

var Options struct {
	Font        string `short:"F" long:"font" description:"Font for event name" default:"fonts/6x12.bdf"`
	Matrix      string `short:"m" long:"matrix" description:"Matrix to connect to" required:"true"`
	SerialName  string `long:"serial-name" description:"Serial device for arduino" default:"/dev/ttyUSB0"`
	SerialBaud  int    `long:"serial-baud" value-name:"BAUD" default:"57600"`
	TextRed     int    `short:"r" long:"red" description:"Red component of text color" default:"128"`
	TextGreen   int    `short:"g" long:"green" description:"Green component of text color" default:"128"`
	TextBlue    int    `short:"b" long:"blue" description:"Blue component of text color" default:"0"`
	LocalTime   string `short:"t" long:"local-time" description:"Local timezone" default:"Europe/Helsinki"`
	ForeignTime string `short:"T" long:"foreign-time" description:"Foreign timezone" default:"Europe/Moscow"`
	TimePin     int    `short:"p" long:"time-pin" description:"Pin to select foreign timezone, active low" default:"15"`
}

var parser = flags.NewParser(&Options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		panic(err)
	}

	textColor := [3]byte{byte(Options.TextRed), byte(Options.TextGreen), byte(Options.TextBlue)}

	serialConfig := serial.Config{
		Name: Options.SerialName,
		Baud: Options.SerialBaud,
		// ReadTimeout:    options.SerialTimeout,
	}

	serial, err := serial.OpenPort(&serialConfig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Serial open.\n")

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}

	timePin, err := embd.NewDigitalPin(Options.TimePin)
	if err != nil {
		panic(err)
	} else if err := timePin.SetDirection(embd.In); err != nil {
		panic(err)
	}

	fmt.Printf("GPIO initialized.\n")

	local, err := time.LoadLocation(Options.LocalTime)
	if err != nil {
		panic(err)
	}

	foreign, err := time.LoadLocation(Options.ForeignTime)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Timezones loaded.\n")

	font, err := bdf.Parse(Options.Font)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Fonts loaded.\n")

	m := matrix.Init(Options.Matrix, 32, 32)
	defer m.Close()

	// Trap SIGINT aka Ctrl-C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	clockBitmap := font.TextBitmap("15:04")
	secondBitmap := font.TextBitmap("05")

	updateTicker := time.NewTicker(time.Millisecond * 10)
	send := make([]byte, 1)

	for {
		select {
		case <-sigChan:
			// SIGINT received, shutdown gracefully
			m.Close()
			os.Exit(1)
		case <-updateTicker.C:
			var t time.Time
			if i, _ := timePin.Read(); err != nil {
				panic(err)
			} else if i == 1 {
				t = time.Now().In(local)
			} else {
				t = time.Now().In(foreign)
			}
			clockBitmap = font.TextBitmap(t.Format("15:04"))
			secondBitmap = font.TextBitmap(t.Format("05"))
			seconds := t.Second()
			m.Fill(matrix.ColorBlack())
			m.Scroll(clockBitmap, textColor, 10, 0, 0, 32)
			m.Scroll(secondBitmap, textColor, 20, 10, 0, 12)
			send[0] = byte(seconds)
			_, err := serial.Write(send)
			if err != nil {
				panic(err)
			}
			m.Send()
		}
	}
}
