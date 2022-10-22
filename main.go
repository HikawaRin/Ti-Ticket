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
	router.POST("/require/grant", oauth.RequireGrant)

	router.GET("/debug/latest_token", oauth.DEBUGToken)

	router.POST("/admin/grant", oauth.GrantPrivilege)
	router.GET("/admin/show/requestes", oauth.ShowRequestes)
	router.POST("/admin/handle/request", oauth.HandleRequest)
	router.POST("/admin/drop", oauth.DropUser)
	log.Fatal(router.Run(":8080"))
}
