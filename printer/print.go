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

	if _, err = p.PrintImage(img); err != nil {
		return err
	}
	if embossed != "" {
		if _, err = p.LineFeed(); err != nil {
			return err
		}
		if _, err = p.Reverse(true).Size(3, 3).Bold(true).Write(embossed); err != nil {
			return err
		}
		if _, err = p.Reverse(false).LineFeed(); err != nil {
			return err
		}
	}
	if _, err = p.Bold(true).Size(2, 2).Write(header); err != nil {
		return err
	}
	if _, err = p.LineFeed(); err != nil {
		return err
	}
	if additional != "" {
		additional = " - " + additional
	}
	now := time.Now().Format("January 02, 2006")
	if _, err = p.Bold(false).Size(1, 1).Justify(escpos.JustifyLeft).Write("Printed " + now + additional); err != nil {
		return err
	}
	if paragraph != "" {
		for range 2 {
			if _, err = p.LineFeed(); err != nil {
				return err
			}
		}
		p.Bold(false).Size(1, 1).Justify(escpos.JustifyLeft).Write(paragraph)
	}
	for range 5 {
		if _, err = p.LineFeed(); err != nil {
			return err
		}
	}

	// You need to use either p.Print() or p.PrintAndCut() at the end to send the data to the printer.
	if err = p.PrintAndCut(); err != nil {
		return err
	}

	return nil
}
