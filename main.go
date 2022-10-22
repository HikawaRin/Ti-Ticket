package main

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ti-ticket/DAO"
	"ti-ticket/cache"
	"ti-ticket/oauth"
)

var (
	router = gin.Default()
)

func main() {
	var err error
	if err = DAO.InitConnect(); err != nil {
		log.Print("Initial Database Connection failure.")
	}
	if err = cache.Init(); err != nil {
		log.Print("Initial service cache failure.")
	}

	router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/login") })
	router.GET("/login", oauth.Authorization)
	router.GET("/oauth/redirect", oauth.Authorized)
	router.GET("/fetch_passwd", oauth.FetchPasswd)

	router.GET("/debug/latest_token")
	router.GET("/debug/drop_me")
	log.Fatal(router.Run(":8080"))
}
