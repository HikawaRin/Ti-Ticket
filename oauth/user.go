package oauth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"ti-ticket/cache"
)

var (
	_github_client_id     string = "da030ce5baaa76bb362d"
	_github_client_secret string = os.Getenv("GITHUB_CLIENT_SECRET")

	_github_redirect_uri  string = "http://localhost:8080/oauth/redirect"
	_github_authtoken_uri string = "https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s"
	_github_api_uri       string = "https://api.github.com/"

	httpClient http.Client = http.Client{}

	// For Debug
	_latest_token string
)

func Authorization(c *gin.Context) {
	target := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s/", _github_client_id, _github_redirect_uri)
	log.Print("Please authorization on Github: ", target)
	c.Redirect(http.StatusFound, target)
}

func Authorized(c *gin.Context) {
	// mask for test
	code, ok := c.GetQuery("code")
	if !ok {
		c.JSON(http.StatusUnauthorized, "Invalid Request.")
		return
	}

	accessToken, err := requestGithubToken(code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprint("Unauthorized with Github. %v", err.Error()))
		return
	}

	user, err := requestGithubUser(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable get Github user.")
		return
	}

	up, ok := cache.AddUser(user.Email)
	// up, ok := cache.AddUser("dummy@dummy.com")
	if !ok {
		c.JSON(http.StatusInternalServerError, "Unable to add user.")
		return
	}
	token, err := createToken(*up)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Construct Token failure.")
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
	_latest_token = token
}

func FetchPasswd(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, "Invalid token content")
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"passwd": user.Password,
	})
}

type githubOAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func requestGithubToken(code string) (string, error) {
	requestURL := fmt.Sprintf(_github_authtoken_uri, _github_client_id, _github_client_secret, code)
	log.Print(requestURL, "\n")
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("accept", "application/json")
	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var token githubOAuthAccessResponse
	if err := json.NewDecoder(response.Body).Decode(&token); err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

type GithubUser struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `josn:"email"`
}

func requestGithubUser(token string) (*GithubUser, error) {
	requestURL := _github_api_uri + "user"
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("accept", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var user GithubUser
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

type titicketClaims struct {
	Account string `json:"account"`
	jwt.StandardClaims
}

func createToken(user cache.User) (string, error) {
	var err error

	claims := titicketClaims{
		Account: user.Account,
		StandardClaims: jwt.StandardClaims{
			Id:        user.Uid,
			ExpiresAt: user.Expire_time,
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	return token, err
}

func extractToken(bearToken string) (string, bool) {
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1], true
	}
	return "", false
}

func vertifyToken(bearToken string) (*titicketClaims, error) {
	token, err := jwt.ParseWithClaims(
		bearToken,
		&titicketClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*titicketClaims)
	if ok && token.Valid {
		if claims.ExpiresAt < time.Now().UTC().Unix() {
			log.Print("Token expired\n")
			return nil, fmt.Errorf("%s", "Token expired.")
		}
	}
	return claims, nil
}

func DEBUGToken(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{
		"token": _latest_token,
	})
}

func DEBUGDropUser(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, "Invalid token content")
		return
	}
	cache.DropUser(user)
	c.JSON(http.StatusOK, map[string]string{
		"Account": user.Account,
	})
}
