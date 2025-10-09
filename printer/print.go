package printer

import (
	"image"
	"log"
	"os"
	"time"

	"github.com/kenshaw/escpos"
	"github.com/tarm/serial"
)

func printImage(img image.Image, printer *escpos.Escpos) error {
	xl, xh, yl, yh, data := renderImage(img)
	_, err := printer.WriteRaw(append([]byte{0x1D, 'v', 48, 0, xl, xh, yl, yh}, data...))
	return err
}

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
	p.SetLogger(log.New(os.Stdout, "ESCPOS: ", 0))
	p.Init()

	if err = printImage(img, p); err != nil {
		return err
	}
	if embossed != "" {
		p.Linefeed()
		p.SetFontSize(3, 3)
		p.Text(map[string]string{
			"em":      "true",
			"reverse": "true",
		}, embossed)
		p.SetReverse(0)
		p.Linefeed()
	}
	p.SetFontSize(2, 2)
	p.Text(map[string]string{
		"em": "true",
	}, header)
	p.Linefeed()

	if additional != "" {
		additional = " - " + additional
	}
	now := time.Now().Format("January 02, 2006")
	p.SetFontSize(1, 1)
	p.Text(map[string]string{}, "Printed "+now+additional)
	p.Linefeed()
	if paragraph != "" {
		p.Linefeed()
		p.Linefeed()
		p.Text(map[string]string{}, paragraph)
	}
	for range 5 {
		p.Linefeed()
	}
	p.Cut()
	p.End()
	return socket.Flush()
}
