package api

import (
	"chatchat/app/api/middleware"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"net/http"
	"time"
)

func InitRouter() error {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.POST("/register", register)
	r.POST("/login", login)
	r.POST("/verificationID", SendMail)
	UserRouter := r.Group("/user")
	{
		UserRouter.Use(middleware.JWTAuthMiddleware())
		UserRouter.POST("/:uid/changePassword", ChangePassword)
		UserRouter.POST("/:uid/changeNickname", ChangeNickname)
		UserRouter.POST("/:uid/changeIntroduction", ChangeIntroduction)
		UserRouter.POST("/:uid/getUser", GetUser)
	}

	r.GET("/oauth2login", func(c *gin.Context) {
		u := config.AuthCodeURL("xyz",
			oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"))
		http.Redirect(c.Writer, c.Request, u, http.StatusFound)
	})

	r.GET("/oauth2", func(c *gin.Context) {
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
	})

	r.GET("oauth2/refresh", func(c *gin.Context) {
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
	})

	r.GET("oauth2/try", func(c *gin.Context) {
		if globalToken == nil {
			http.Redirect(c.Writer, c.Request, "/oauth2login", http.StatusFound)
			return
		}

		resp, err := http.Get(fmt.Sprintf("%s/verify?access_token=%s", authServerURL, globalToken.AccessToken))
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		io.Copy(c.Writer, resp.Body)
	})

	r.GET("oauth2/pwd", func(c *gin.Context) {
		token, err := config.PasswordCredentialsToken(context.Background(), "2022214740", "666666666")
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(c.Writer)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	r.GET("oauth2/client", func(c *gin.Context) {
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
	})
	err := r.Run(":8088")
	if err != nil {
		return err
	} else {
		return nil
	}
}
