package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ents-source/member-label-print/printer"
)

func doPrint(w http.ResponseWriter, r *http.Request) {
	fob := r.URL.Query().Get("fob")
	if fob != "" {
		if err := doPrintFob(fob); err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	label := r.URL.Query().Get("label")
	if label != "" {
		expireDate := r.URL.Query().Get("expireDate")
		if expireDate != "" {
			ms, err := strconv.ParseInt(expireDate, 10, 64)
			if err != nil {
				panic(err)
			}
			expireDate = time.UnixMilli(ms).Format("Jan 02")
		}
		if err := doPrintLabel(r.URL.Query().Get("kind"), label, r.URL.Query().Get("paragraph"), r.URL.Query().Get("expireLabel"), expireDate); err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func doPrintFob(fob string) error {
	log.Println("PRINT FOB: ", fob)
	users, err := amemberProApi.FindUsersByFob(fob)
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return fmt.Errorf("expected 1 user record, got %d", len(users))
	}

	return printer.DoPrint("", users[0].Name(), "ID: "+strconv.Itoa(users[0].Id), "", "", "")
}

func doPrintLabel(embossed string, label string, paragraph string, expireLabel string, expireDate string) error {
	log.Println("PRINT LABEL: ", embossed, label, paragraph)
	return printer.DoPrint(embossed, label, "", paragraph, expireLabel, expireDate)
}
