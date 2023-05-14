package api

import (
	"chatchat/app/api/middleware"
	"chatchat/app/global"
	"chatchat/dao/mysql"
	"chatchat/dao/redis"
	"chatchat/model"
	"chatchat/utils"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	authServerURL = "http://162.14.78.209:9096"
)

var (
	config = oauth2.Config{
		ClientID:     "test_client_1",
		ClientSecret: "test_secret_1",
		Scopes:       []string{"all"},
		RedirectURL:  "http://162.14.78.209:8088/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}
	globalToken *oauth2.Token // Non-concurrent security
)

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
func Oauth2Register(c *gin.Context) {
	username, f := c.GetPostForm("username")
	mailID, f := c.GetPostForm("mailID")
	if f == false {
		utils.ResponseFail(c, "verification failed")
	}
	var user model.OauthUser
	c.ShouldBind(&user)
	uid, _ := redis.Get(c, fmt.Sprintf("Rmail:%s", username))
	if uid != mailID {
		utils.ResponseFail(c, "wrong mailID")
		return
	}
	ID := global.Rdb.ZCard(c, "userID").Val() + 1
	flag, msg := mysql.AddOauth2User(username, strconv.FormatInt(user.Oauth2Username, 10))
	flag, msg = mysql.AddUser(username, "", user.Nickname, ID) //写入数据库
	if !flag {
		utils.ResponseFail(c, fmt.Sprintf("write into mysql failed,%s", msg))
		return
	}
	err := redis.ZSetUserID(c, username)
	if err != nil {
		utils.ResponseFail(c, err.Error())
	}
	redis.HSet(c, "Oauth2User", user.Oauth2Username, username)
	redis.HSet(c, fmt.Sprintf("user:%s", username), "nickname", user.Nickname)
	redis.HSet(c, fmt.Sprintf("user:%s", username), "avatar", user.Avatar)
	redis.HSet(c, fmt.Sprintf("user:%s", username), "id", ID)
	claim := model.MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    "Wzy",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, _ := token.SignedString(middleware.Secret)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "register success",
		"token":  tokenString,
	})
}
func Oauth2Try(c *gin.Context) {
	if globalToken == nil {
		http.Redirect(c.Writer, c.Request, "/oauth2login", http.StatusFound)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/verify?access_token=%s", authServerURL, globalToken.AccessToken))
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	var user model.OauthUser
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		utils.ResponseFail(c, err.Error())
		return
	}
	//println(bodyBytes)
	//println(user)
	//io.Copy(c.Writer, resp.Body)
	username, _ := redis.HGet(c, "Oauth2User", strconv.FormatInt(user.Oauth2Username, 10))
	if username != "" {
		claim := model.MyClaims{
			Username: username.(string),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
				Issuer:    "Wzy",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenString, _ := token.SignedString(middleware.Secret)
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "login success",
			"token":   tokenString,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "base_user doesn't exist, please register a new one",
		"data":    user,
	})
}

func Oauth2Pwd(c *gin.Context) {
	token, err := config.PasswordCredentialsToken(context.Background(), "2022214740", "666666666")
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	globalToken = token
	e := json.NewEncoder(c.Writer)
	e.SetIndent("", "  ")
	e.Encode(token)
}

func Oauth2Client(c *gin.Context) {
	cfg := clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     config.Endpoint.TokenURL,
	}

	token, err := cfg.Token(context.Background())
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(c.Writer)
	e.SetIndent("", "  ")
	e.Encode(token)
}

func Oauth2Refresh(c *gin.Context) {
	if globalToken == nil {
		http.Redirect(c.Writer, c.Request, "/", http.StatusFound)
		return
	}

	globalToken.Expiry = time.Now()
	token, err := config.TokenSource(context.Background(), globalToken).Token()
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	globalToken = token
	e := json.NewEncoder(c.Writer)
	e.SetIndent("", "  ")
	e.Encode(token)
}

func Oauth2(c *gin.Context) {
	c.Request.ParseForm()
	state := c.Request.Form.Get("state")
	if state != "xyz" {
		http.Error(c.Writer, "State invalid", http.StatusBadRequest)
		return
	}
	code := c.Request.Form.Get("code")
	if code == "" {
		http.Error(c.Writer, "Code not found", http.StatusBadRequest)
		return
	}
	token, err := config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", "s256example"))
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	globalToken = token

	e := json.NewEncoder(c.Writer)
	e.SetIndent("", "  ")
	e.Encode(token)
}

func Oauth2Login(c *gin.Context) {
	u := config.AuthCodeURL("xyz",
		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	http.Redirect(c.Writer, c.Request, u, http.StatusFound)
}
