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
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, defaultUrl)
	})
	r.GET("/:hash", expandUrl)

	r.Run(port)
}

// look for long url from redis by hash
// it will redirect to default host when short url is not exist
func expandUrl(c *gin.Context) {
	hash := c.Param("hash")
	if hash == except {
		return
	}

	url, _ := models.Find(hash)
	if url != "" {
		c.Redirect(http.StatusFound, url)
		return
	}

	c.Redirect(http.StatusFound, defaultUrl)
	return
}
