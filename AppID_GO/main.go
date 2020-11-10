package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tkanos/gonfig"
	"golang.org/x/oauth2"
)

var public string = "public/"

var conf oauth2.Config

var appIDConfig struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	RedirectURL  string
}

type middleware func(http.HandlerFunc) http.HandlerFunc
type server struct{}

func (s *server) ServerHTTP(ses http.ResponseWriter, req *http.Request) {

}

func (s *server) addmiddleware(hf http.HandlerFunc, hfs ...middleware) http.HandlerFunc {
	for _, m := range hfs {
		hf = m(hf)
	}
	return hf
}

func auth() middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			_, err := req.Cookie("session_token")
			if err != nil {
				http.Redirect(res, req, "/login", http.StatusFound)
				return
			}
			f(res, req)
		}
	}
}

func authout() middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			_, err := req.Cookie("session_token")
			if err != nil {
				f(res, req)
				return
			}
			http.Redirect(res, req, "/", http.StatusFound)
		}
	}
}
func static(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, public+req.URL.Path[1:])
}

func login(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, public+"login.html")
}

func logout(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session_token")

	if err != nil {

		log.Println("No session cookie found:" + err.Error())

	} else {

		log.Println("Session cookie found, invalidating it.")

		// If cookie was found, let's invalidate it
		cookie.MaxAge = -1

	}

	// Setting the invalidated cookie
	http.SetCookie(res, cookie)
	http.ServeFile(res, req, public+"login.html")
}

func loginwithIBM(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, conf.AuthCodeURL("state"), http.StatusFound)
}

func callback(res http.ResponseWriter, req *http.Request) {
	keys := req.URL.Query()["code"]
	if keys == nil {
		res.Write([]byte("Error with authentication"))
		return
	}
	//get token
	fmt.Println(keys[0])

	token, _ := conf.Exchange(context.Background(), keys[0])
	client := conf.Client(context.Background(), &oauth2.Token{AccessToken: token.AccessToken})

	userinfo, _ := client.Get(strings.Replace(conf.Endpoint.AuthURL, "/authorization", "/userinfo", 1))
	var profile map[string]interface{}
	error := json.NewDecoder(userinfo.Body).Decode(&profile)
	if error != nil {
		res.Write([]byte("Error with authentication"))
		return
	}
	http.SetCookie(res, &http.Cookie{
		Name:    "session_token",
		Value:   keys[0],
		Path:    "/",
		Expires: time.Now().Add(1000 * time.Second),
	})
	http.Redirect(res, req, "/", http.StatusFound)
}

func main() {
	//Read config file
	appIDConfig := appIDConfig
	err := gonfig.GetConf("config/IBMAppID.json", &appIDConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conf.RedirectURL = appIDConfig.RedirectURL
	conf.ClientID = appIDConfig.ClientID
	conf.ClientSecret = appIDConfig.ClientSecret
	conf.Scopes = []string{"openid", "profile"}
	conf.Endpoint = oauth2.Endpoint{
		AuthURL:  appIDConfig.AuthURL + "/authorization",
		TokenURL: appIDConfig.AuthURL + "/token",
	}

	s := server{}
	http.Handle("/login", s.addmiddleware(login, authout()))
	http.Handle("/logout", s.addmiddleware(logout))
	http.Handle("/loginwithibm", s.addmiddleware(loginwithIBM, authout()))
	http.HandleFunc("/auth/callback", callback)
	http.Handle("/", s.addmiddleware(static, auth()))
	http.ListenAndServe(":3000", nil)
}
