package printer

import (
	"github.com/hennedo/escpos"
	"github.com/tarm/serial"
	"time"
)

func DoPrint(header string, additional string) error {
	config := &serial.Config{
		Name:        SerialPort,
		Baud:        9600,
		ReadTimeout: 1,
		Size:        8,
		StopBits:    serial.Stop1,
		Parity:      serial.ParityNone,
	}

	socket, err := serial.OpenPort(config)
	if err != nil {
		return err
	}
	defer socket.Close()

	p := escpos.New(socket)
	p.SetConfig(escpos.ConfigEpsonTMT20II) // Removes trailing zeros from lines

	p.PrintImage(img)
	p.Bold(true).Size(2, 2).Write(header)
	p.LineFeed()
	if additional != "" {
		additional = " - " + additional
	}
	now := time.Now().Format("January 02, 2006")
	p.Bold(false).Size(1, 1).Justify(escpos.JustifyLeft).Write("Printed " + now + additional)
	p.LineFeed()
	p.LineFeed()
	p.LineFeed()
	p.LineFeed()
	p.LineFeed()

	// You need to use either p.Print() or p.PrintAndCut() at the end to send the data to the printer.
	p.PrintAndCut()

	return nil
}
