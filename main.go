package main

import (
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/speps/go-hashids"
	"net/http"
	"time"
)

const (
	hdSalt        = "mysalt"
	hdMinLength   = 5
	defaultDomain = "http://example.com/"
)

var (
	// 定义常量
	RedisClient *redis.Pool
	REDIS_HOST  = "127.0.0.1:6379"
	REDIS_DB    = 0

	db      *sql.DB
	DB_HOST = "tcp(127.0.0.1:3306)"
	DB_NAME = "short"
	DB_USER = "root"
	DB_PASS = ""
)

func main() {
	initRedis()
	initMysql()

	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		// 生成
		v1.POST("/short", shortUrl)

		// 根据HASH取网址
		v1.GET("/expand/:hash", expandUrlApi)
	}

	// 根据HASH跳转
	r.GET("/s/:hash", expandUrl)
	// 主页跳转
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, defaultDomain)
	})

	r.Run(":8080")
}

func initRedis() {
	// 建立连接池
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     1,
		MaxActive:   10,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}
}

func initMysql() {
	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"
	db, _ = sql.Open("mysql", dsn)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(20)
	db.Ping()
}

// 接受生成短网址请求
// 参数: 长URL
// 返回值: HASH
// @TODO 校验URL格式
func shortUrl(c *gin.Context) {
	longUrl := c.PostForm("url")

	if len(longUrl) > 0 {

		if hash, ok := insert(longUrl); ok {
			c.JSON(200, gin.H{
				"status":  200,
				"message": "ok",
				"hash":    hash,
			})
		}

	} else {
		c.JSON(200, gin.H{
			"status":  500,
			"message": "请传入网址",
		})
	}
}

// 根据HASH解析并跳转到对应的长URL
// 不存在则跳转到默认地址
func expandUrl(c *gin.Context) {
	hash := c.Param("hash")

	if url, ok := findByHash(hash); ok {
		c.Redirect(http.StatusMovedPermanently, url)
	}
	// 注意:
	// 	实际中，此应用的运行域名可能与默认域名不同，如a.com运行此程序，默认域名为b.com
	// 	当访问一个不存在的HASH或a.com时，可以跳转到任意域名，即defaultDomain
	c.Redirect(http.StatusMovedPermanently, defaultDomain)
}

// 根据HASH在redis中查找并返回结果
// 不存在则返回404状态和默认网址
// 
// 不存在可能是redis中没有，此处没有再次检查MySQL中是否存在
// 可以根据实际情况做调整
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

	// 此处可以尝试在MySQL中再次查询

	c.JSON(200, gin.H{
		"status":  404,
		"message": "url of hash is not exist",
		"data":    defaultDomain,
	})
}

// HASH转ID后查找
// 与上一方法类似，将HASH转成ID后在数据库中查找
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
		"data":    defaultDomain,
	})
}

// 将ID转换成对应的HASH值
// hdSalt与hdMinLength 会影响生成结果，确定后尽量不要改动
func short(id int) string {
	hd := hashids.NewData()
	hd.Salt = hdSalt
	hd.MinLength = hdMinLength

	h := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{id})

	return e
}

// 根据HASH解析出对应的ID值
// hdSalt与hdMinLength 会影响生成结果，确定后尽量不要改动
func expand(hash string) int {
	hd := hashids.NewData()
	hd.Salt = hdSalt
	hd.MinLength = hdMinLength

	h := hashids.NewWithData(hd)
	d, _ := h.DecodeWithError(hash)

	return d[0]
}

// 数据库中根据ID查找
func find(id int) (string, bool) {
	var url string
	err := db.QueryRow("SELECT url FROM url WHERE id = ?", id).Scan(&url)
	if err == nil {
		return url, true
	} else {
		return "", false
	}
}

// redis中根据HASH查找
func findByHash(h string) (string, bool) {
	rc := RedisClient.Get()

	defer rc.Close()
	url, _ := redis.String(rc.Do("GET", "URL:"+h))

	if url != "" {
		return url, true
	}

	id := expand(h)
	if urldb, ok := find(id); ok {
		return urldb, true
	}

	return "", false
}

// 将长网址插入到数据库中
// 并把返回的ID生成HASH和长网址存入redis
func insert(url string) (string, bool) {
	stmt, _ := db.Prepare(`INSERT INTO url (url) values (?)`)
	res, err := stmt.Exec(url)
	checkErr(err)

	id, _ := res.LastInsertId()

	rc := RedisClient.Get()
	defer rc.Close()

	hash := short(int(id))
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
