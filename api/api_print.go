package api

import (
	"fmt"
	"github.com/ents-source/member-label-print/amember"
	"github.com/ents-source/member-label-print/printer"
	"log"
	"net/http"
	"strconv"
)

func doPrint(w http.ResponseWriter, r *http.Request) {
	fob := r.URL.Query().Get("fob")
	if fob != "" {
		doPrintFob(fob)
		w.WriteHeader(http.StatusOK)
		return
	}

	label := r.URL.Query().Get("label")
	if label != "" {
		doPrintFob(label)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func doPrintFob(fob string) error {
	log.Println("PRINT FOB: ", fob)
	users, err := amember.FindUsersByFob(fob)
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return fmt.Errorf("expected 1 user record, got %d", len(users))
	}

	name := users[0].Nickname
	if name == "" {
		name = users[0].FirstName + " " + users[0].LastName
	}

	return printer.DoPrint(users[0].Nickname, "ID: "+strconv.Itoa(users[0].Id))
}

func doPrintLabel(label string) error {
	// TODO
	return nil
}
