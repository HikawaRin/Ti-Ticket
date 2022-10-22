package oauth

import (
	"fmt"
	"log"
	"net/http"
	"ti-ticket/DAO"
	"ti-ticket/cache"

	"github.com/gin-gonic/gin"
)

type peivilegeCode struct {
	User string `json:"user"`
	Code int    `json:"code"`
}

func GrantPrivilege(c *gin.Context) {
	var pc peivilegeCode
	if err := c.ShouldBindJSON(&pc); err != nil {
		c.JSON(http.StatusBadRequest, "Post Body format error")
		return
	}

	bearToken, ok := extractToken(c.Request.Header.Get("Authorization"))
	if !ok {
		c.JSON(http.StatusUnauthorized, "Invalid token format.")
		return
	}
	claims, err := vertifyToken(bearToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	user, ok := cache.FetchUser(claims.Id)
	if !ok {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("%s", "Invalid user"))
	}
	log.Print("Operator: ", user.Account, "\n")
	// TODO Need check if is admin token
	priv, err := DAO.GrantUser(pc.User, pc.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, priv)
}

func RequireGrant(c *gin.Context) {
	var pc peivilegeCode
	if err := c.ShouldBindJSON(&pc); err != nil {
		c.JSON(http.StatusBadRequest, "Post Body format error")
		return
	}

	bearToken, ok := extractToken(c.Request.Header.Get("Authorization"))
	if !ok {
		c.JSON(http.StatusUnauthorized, "Invalid token format.")
		return
	}
	claims, err := vertifyToken(bearToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	user, ok := cache.FetchUser(claims.Id)
	if !ok || user.Account != pc.User {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("%s", "Invalid user"))
	}
	cache.ReceiveRequest(pc.User, pc.Code)
	c.JSON(http.StatusOK, "Wait Admin approve")
}

func ShowRequestes(c *gin.Context) {
	bearToken, ok := extractToken(c.Request.Header.Get("Authorization"))
	if !ok {
		c.JSON(http.StatusUnauthorized, "Invalid token format.")
		return
	}
	claims, err := vertifyToken(bearToken)
	if err != nil || claims.Account == "" {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	// TODO user should be admin
	// user, ok := cache.FetchUser(claims.Id)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, fmt.Errorf("%s", "Invalid user"))
	// }
	requests := cache.ListRequest()
	c.JSON(http.StatusOK, *requests)
}

func HandleRequest(c *gin.Context) {
	type adminOp struct {
		Op bool `json:"op"`
	}
	var op adminOp
	if err := c.ShouldBindJSON(&op); err != nil {
		c.JSON(http.StatusBadRequest, "Post Body format error")
		return
	}

	bearToken, ok := extractToken(c.Request.Header.Get("Authorization"))
	if !ok {
		c.JSON(http.StatusUnauthorized, "Invalid token format.")
		return
	}
	claims, err := vertifyToken(bearToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	user, ok := cache.FetchUser(claims.Id)
	if !ok {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("%s", "Invalid user"))
	}
	log.Print("Operator: ", user.Account, "\n")
	// TODO Need check if is admin token
	priv, err := cache.HandleRequest(op.Op)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, priv)
}
