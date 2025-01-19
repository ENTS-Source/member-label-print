package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/ents-source/go-amember-api/amember"
	"github.com/ents-source/member-label-print/api"
	"github.com/ents-source/member-label-print/assets"
	"github.com/ents-source/member-label-print/printer"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	HttpBind string `envconfig:"http_bind" default:"0.0.0.0:8080"`

	PrinterPort string `envconfig:"serial_port"`

	AmpApiKey     string `envconfig:"amp_api_key"`
	AmpApiKeyFile string `envconfig:"amp_api_key_file"`
	AmpApiUrl     string `envconfig:"amp_api_url"`
}

func main() {
	var c config
	err := envconfig.Process("mlp", &c)
	if err != nil {
		log.Fatal(err)
	}

	devMode := os.Getenv("DEV") == "true"

	webPath := assets.SetupWeb()
	if devMode {
		webPath = "./web"
	}

	paymentsApi := amember.NewClient(c.AmpApiUrl, getPassword(c.AmpApiKey, c.AmpApiKeyFile))

	printer.SerialPort = c.PrinterPort
	err = printer.LoadLogo(path.Join(webPath, "logo.png"))
	if err != nil {
		log.Fatal(err)
	}

	wg := api.Start(c.HttpBind, webPath, paymentsApi)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer close(stop)
		<-stop

		log.Println("Stopping api...")
		api.Stop()

		log.Println("Cleaning up...")
		_ = os.RemoveAll(webPath)

		log.Println("Done stopping")
	}()

	wg.Add(1)
	wg.Wait()

	log.Println("Goodbye!")
}

func getPassword(in string, f string) string {
	passwd := in
	if f != "" {
		b, err := os.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		passwd = string(b)
	}
	return passwd
}
