package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qichengzx/shortme/config"
	"github.com/qichengzx/shortme/models"
	"log"
	"net/http"
)

var (
	defaultUrl string
	port       string
	except     = "favicon.ico"
	err        error
)

func init() {
	defaultUrl, err = config.GetByBlock("common", "defaulturl")
	port, err = config.GetByBlock("common", "appport")
	if err != nil {
		log.Printf("func config.Get got err : %+v", err)
	}
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("source/template/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	// TODO
	//r.GET("/:hash", expandUrl)

	r.POST("shorten", shortenweb)
	r.GET("lookup/:hash", lookupWeb)

	api := r.Group("/api")
	{
		api.POST("shorten", shorten)
		api.GET("lookup/:hash", lookup)
	}

	r.Run(port)
}

func shortenweb(c *gin.Context) {
	longUrl := c.PostForm("long_url")
	hash := c.PostForm("hash")
	url, err := models.Save(longUrl, hash)
	if err != nil {
		panic(err)
	}

	c.HTML(200, "result.html", gin.H{
		"ok":     true,
		"result": defaultUrl + url.HASH,
	})
}

// look for long url from redis by hash
// it will redirect to default host when short url is not exist
func expandUrl(c *gin.Context) {
	hash := c.Param("hash")
	if hash == except {
		return
	}

	url := models.FindLongUrl(hash)
	if url != "" {
		c.Redirect(http.StatusFound, url)
		return
	}

	c.Redirect(http.StatusFound, defaultUrl)
	return
}

func shorten(c *gin.Context) {
	longUrl := c.PostForm("long_url")
	hash := c.PostForm("hash")
	url, err := models.Save(longUrl, hash)
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"action": "shorten",
		"result": defaultUrl + url.HASH,
	})
}

func lookupWeb(c *gin.Context) {
	hash := c.Param("hash")

	url := models.Find(hash)

	c.HTML(200, "result.html", gin.H{
		"ok":     true,
		"result": defaultUrl + url.HASH,
	})
}

func lookup(c *gin.Context) {
	hash := c.Param("hash")

	url := models.Find(hash)
	c.JSON(200, gin.H{
		"action": "lookup",
		"result": url,
	})
}
