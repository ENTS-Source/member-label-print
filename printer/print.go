package printer

import (
	"time"

	"github.com/hennedo/escpos"
	"github.com/tarm/serial"
)

func DoPrint(embossed string, header string, additional string, paragraph string) error {
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
	if embossed != "" {
		p.LineFeed()
		p.Reverse(true).Size(3, 3).Bold(true).Write(embossed)
		p.Reverse(false).LineFeed()
	}
	p.Bold(true).Size(2, 2).Write(header)
	p.LineFeed()
	if additional != "" {
		additional = " - " + additional
	}
	now := time.Now().Format("January 02, 2006")
	p.Bold(false).Size(1, 1).Justify(escpos.JustifyLeft).Write("Printed " + now + additional)
	if paragraph != "" {
		p.LineFeed()
		p.LineFeed()
		p.Bold(false).Size(1, 1).Justify(escpos.JustifyLeft).Write(paragraph)
	}
	p.LineFeed()
	p.LineFeed()
	p.LineFeed()
	p.LineFeed()
	p.LineFeed()

	// You need to use either p.Print() or p.PrintAndCut() at the end to send the data to the printer.
	p.PrintAndCut()

	return nil
}
