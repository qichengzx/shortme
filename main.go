package main

import (
	"fmt"
	"strings"
	"time"
	"net/http"
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"

	"github.com/speps/go-hashids"
)

const (
	hdSalt        = "mysalt"
	hdMinLength   = 5
	defaultDomain = "http://localhost/"
)

var (
	RedisClient *redis.Pool
	RedisHost   = "127.0.0.1:6379"
	RedisDb     = 0
	RedisPwd    = ""

	db      *sql.DB
	DbHost = "tcp(127.0.0.1:3306)"
	DbName = "short"
	DbUser = "root"
	DbPass = ""
)

func main() {
	initRedis()
	initMysql()

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		//http code can be StatusFound or StatusMovedPermanently 
		c.Redirect(http.StatusFound, defaultDomain)
	})
	r.GET("/:hash", expandUrl)
	r.GET("/:hash/info", expandUrlApi)
	r.POST("/short", shortUrl)

	r.Run(":8000")
}

func initRedis() {
	RedisClient = &redis.Pool{
		MaxIdle:     1,
		MaxActive:   10,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", RedisHost)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", RedisPwd); err != nil {
				c.Close()
				return nil, err
			}
			c.Do("SELECT", RedisDb)
			return c, nil
		},
	}
}

func initMysql() {
	dsn := DbUser + ":" + DbPass + "@" + DbHost + "/" + DbName + "?charset=utf8"
	db, _ = sql.Open("mysql", dsn)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(20)
	db.Ping()
}

// receive a url from www-form encoded
// return a shortlen url
func shortUrl(c *gin.Context) {
	longUrl := c.PostForm("url")

	if longUrl == "" {
		c.JSON(200, gin.H{
			"status":  500,
			"message": "请传入网址",
		})
		return
	}

	if !strings.HasPrefix(longUrl, "http") {
		longUrl = "http://" + longUrl
	}

	if hash, ok := insert(longUrl); ok {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "ok",
			"short":   defaultDomain + hash,
		})
	}
}

// look for long url from redis by hash
// it will redirect to default host when long url is not exist
func expandUrl(c *gin.Context) {
	hash := c.Param("hash")

	if url, ok := findByHash(hash); ok {
		c.Redirect(http.StatusFound, url)
	}

	c.Redirect(http.StatusFound, defaultDomain)
}

// look for long url from redis by hash
// return the long url if result exist or 404 if result is null
func expandUrlApi(c *gin.Context) {
	hash := c.Param("hash")

	if url, ok := findByHash(hash); ok {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "ok",
			"data":    url,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  404,
		"message": "url of hash is not exist",
	})
}

// look for long url from MySQL by id
func expandUrlApi2(c *gin.Context) {
	hash := c.Param("hash")
	id := expand(hash)

	if url, ok := find(id); ok {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "ok",
			"data":    url,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  404,
		"message": "url of hash is not exist",
	})
}

// shortlen a url by id
// the result depends on hdSalt and hdMinLength
func shortenURL(id int) string {
	hd := hashids.NewData()
	hd.Salt = hdSalt
	hd.MinLength = hdMinLength

	h := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{id})

	return e
}

// generate a ID from HASH by hashids
// the result depends on hdSalt and hdMinLength
func expand(hash string) int {
	hd := hashids.NewData()
	hd.Salt = hdSalt
	hd.MinLength = hdMinLength

	h := hashids.NewWithData(hd)
	d, _ := h.DecodeWithError(hash)

	return d[0]
}

// look for url in the mysql by id
func find(id int) (string, bool) {
	var url string
	err := db.QueryRow("SELECT url FROM url WHERE id = ?", id).Scan(&url)
	if err == nil {
		return url, true
	} else {
		return "", false
	}
}

// findByHash for find in redis by hash
func findByHash(h string) (string, bool) {
	rc := RedisClient.Get()

	defer rc.Close()
	url, _ := redis.String(rc.Do("GET", "URL:"+h))

	if url != "" {
		return url, true
	}

	// if the redis result is null ,  in the mysql
	id := expand(h)
	if urldb, ok := find(id); ok {
		return urldb, true
	}

	return "", false
}

// add long url to mysql and use the result ID to generate hash
// add hash and long url to redis
func insert(url string) (string, bool) {
	stmt, _ := db.Prepare(`INSERT INTO url (url) values (?)`)
	res, err := stmt.Exec(url)
	checkErr(err)

	id, _ := res.LastInsertId()

	rc := RedisClient.Get()
	defer rc.Close()

	hash := shortenURL(int(id))
	rc.Do("SET", "URL:"+hash, url)

	return hash, true
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
