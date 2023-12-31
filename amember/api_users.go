package amember

import (
	"net/url"
)

type User struct {
	Id        int      `json:"user_id"`
	Fob       string   `json:"fob"`
	FobAccess string   `json:"fob_access"`
	Announce  []string `json:"announce"`
	Nickname  string   `json:"nickname"`
	FirstName string   `json:"name_f"`
	LastName  string   `json:"name_l"`
}

func FindUsersByFob(fob string) ([]User, error) {
	query := url.Values{}
	query.Set("_filter[fob]", fob)
	return allPages[User]("/api/users", query)
}
