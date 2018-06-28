package models

import (
	"database/sql"
	"fmt"
	"github.com/qichengzx/shortme/utils"
	"time"
)

type URL struct {
	ID        int64      `json:"-"`
	LongUrl   string     `json:"long_url"`
	HASH      string     `json:"hahs"`
	Clicks    int64      `json:"clicks"`
	DeletedAt *time.Time `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
}

func Find(hash string) *URL {
	var u = new(URL)
	u.HASH = hash
	u.Find()

	return u
}

func Save(longUrl, customHash string) (*URL, error) {
	var url = new(URL)
	url.CreatedAt = time.Now().UTC()
	res, err := db.Exec("INSERT INTO `links` (`long_url`,`hash`,`created_at`) values (?,?,?)", longUrl, customHash, url.CreatedAt)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	url.ID = id
	url.LongUrl = longUrl

	if customHash == "" {
		hash := utils.Encode(id)
		db.Exec("UPDATE `links` SET `hash`=? WHERE `id`=?", hash, id)
		url.HASH = hash
	} else {
		url.HASH = customHash
	}

	return url, err
}

func FindLongUrl(hash string) string {
	var u = new(URL)
	u.HASH = hash
	u.Find()

	return u.LongUrl
}

func (u *URL) Find() *URL {
	var url string
	var clicks int64
	var created_at time.Time

	err := db.QueryRow("SELECT `long_url`,`clicks`,`created_at` FROM `links` WHERE `hash` = ?", u.HASH).Scan(&url, &clicks, &created_at)
	if err != nil {
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
		u.Clicks = clicks
		u.CreatedAt = created_at
	}

	return u
}
