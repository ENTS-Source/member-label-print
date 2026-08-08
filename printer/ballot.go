package printer

import (
	"fmt"
	"log"
	"os"

	"github.com/kenshaw/escpos"
	"github.com/tarm/serial"
)

func DoPrintBallot(year int, contested bool, candidates []string) error {
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

	p.SetFontSize(3, 3)
	p.Text(map[string]string{
		"em": "true",
	}, fmt.Sprintf("AGM %d Ballot", year))
	p.Linefeed()

	p.Linefeed()
	p.SetFontSize(1, 1)
	body := "Clearly mark the box next to each of the\ncandidates you are voting to become a Director.\n\nYou may choose multiple options."
	if contested {
		body = fmt.Sprintf("Clearly rank each of the candidates in order of preference for becoming a Director.\n\n1 is your first choice.\n%d is your last choice.", len(candidates))
	}
	p.Text(map[string]string{}, body)
	p.Linefeed()

	p.Linefeed()
	for _, name := range candidates {
		p.SetFontSize(2, 2)
		p.Text(map[string]string{}, fmt.Sprintf("[   ] %s", name))
		p.SetFontSize(1, 1)
		p.Linefeed()
		p.Text(map[string]string{}, "")
		p.Linefeed()
	}
	p.Linefeed()

	for range 5 {
		p.Linefeed()
	}
	p.Cut()
	p.End()
	return socket.Flush()
}
