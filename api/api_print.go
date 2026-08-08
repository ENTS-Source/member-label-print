package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ents-source/member-label-print/printer"
)

var AMPCategoryId int

func doPrint(w http.ResponseWriter, r *http.Request) {
	fob := r.URL.Query().Get("fob")
	if fob != "" {
		if err := doPrintFob(fob); err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	ballot := r.URL.Query().Get("ballot")
	if ballot != "" {
		if err := doPrintBallot(); err != nil {
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

func doPrintBallot() error {
	log.Println("PRINT BALLOT")
	return printer.DoPrintBallot(2026, false, []string{
		"Travis Ralston",
		"Michael Kulpa",
		"Karly Kulpa",
		"Rudy Janvier",
		"Chris Bacivan",
	})
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

func doCheck(w http.ResponseWriter, r *http.Request) {
	fob := r.URL.Query().Get("fob")
	ok, name, err := doCheckMembership(fob)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	_, _ = w.Write([]byte(name))
}

func doCheckMembership(fob string) (bool, string, error) {
	log.Println("CHECK MEMBERSHIP: ", fob)
	users, err := amemberProApi.FindUsersByFob(fob)
	if err != nil {
		return false, fob, err
	}
	if len(users) != 1 {
		return false, fob, fmt.Errorf("expected 1 user record, got %d", len(users))
	}

	productIds, err := amemberProApi.GetProductIdsForCategory(AMPCategoryId)
	//log.Println("PRODUCT IDS: ", productIds)
	if err != nil {
		return false, fob, err
	}
	name := fmt.Sprintf("%s %s (%s; %d)", users[0].FirstName, users[0].LastName, users[0].Nickname, users[0].Id)
	for _, access := range users[0].Nested.Access {
		ok, err := access.IsActiveOn(time.Now())
		if err != nil {
			return false, name, err
		}

		//log.Printf("ACCESS: %#v (%t)", access, ok)
		if ok {
			for _, productId := range productIds {
				if access.ProductId == strconv.Itoa(productId) {
					return true, name, nil
				}
			}
		}
	}

	return false, name, nil
}
