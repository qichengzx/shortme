package models

import (
	"database/sql"
	"fmt"
	"time"
)

type URL struct {
	ID        int
	URL       string
	HASH      string
	Clicks    int64
	DeletedAt *time.Time
	CreatedAt time.Time
}

func Find(hash string) (string, error) {
	var u = new(URL)
	u.HASH = hash
	u.Find()

	return u.URL, nil
}

func (u *URL) Find() *URL {
	var url string

	err := db.QueryRow("SELECT `long_url` FROM `links` WHERE `hash` = ?", u.HASH).Scan(&url)
	if err == nil {
		return nil
	}
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		fmt.Printf("SQL query got error : %+v\n", err)
		return nil
	default:
		u.URL = url
	}

	return u
}
