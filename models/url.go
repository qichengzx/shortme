package models

import (
	"database/sql"
	"fmt"
	"github.com/qichengzx/shortme/utils"
	"time"
)

type URL struct {
	ID        int64
	LongUrl   string
	HASH      string
	Clicks    int64
	DeletedAt *time.Time
	CreatedAt time.Time
}

func Find(hash string) (string, error) {
	var u = new(URL)
	u.HASH = hash
	u.Find()

	return u.LongUrl, nil
}

func Save(longUrl string) (*URL, error) {
	var url = new(URL)
	url.CreatedAt = time.Now().UTC()
	res, err := db.Exec("INSERT INTO `links` (`long_url`,`hash`,`created_at`) values (?,?,?)", longUrl, "", url.CreatedAt)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	hash := utils.Encode(id)
	db.Exec("UPDATE `links` SET `hash`=? WHERE `id`=?", hash, id)

	url.ID = id
	url.LongUrl = longUrl
	url.HASH = hash

	return url, err
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
		u.LongUrl = url
	}

	return u
}
