package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var env Enviroment

type Enviroment struct {
	ClientID     string
	ClientSecret string
}

func init() {
	godotenv.Load()
	env = Enviroment{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!", r.URL.Path[1:])
	})
	http.HandleFunc("/login", githubLoginHandler)
	http.HandleFunc("/login/callback", githubCallbackHandler)
	log.Println("Listening on port 8080")
	port := ":8080"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}
	http.ListenAndServe(port, nil)
}

// githubLoginHandler handles login to github
func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, "https://github.com/login/oauth/authorize", nil)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	// TODO: consider adding a random state to query params
	q.Add("client_id", env.ClientID)
	req.URL.RawQuery = q.Encode()

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

// githubCallbackHandler handles github callback url
func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	type Res struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}
	code := r.URL.Query().Get("code")
	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	q := req.URL.Query()
	q.Add("client_id", env.ClientID)
	q.Add("client_secret", env.ClientSecret)
	q.Add("code", code)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var response Res
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &response)
	fmt.Println(varifyUser(response.AccessToken))
}

func varifyUser(token string) bool {
	type Res struct {
		EventsURL         string      `json:" events_url"`
		AvatarURL         string      `json:"avatar_url"`
		Bio               interface{} `json:"bio"`
		Blog              string      `json:"blog"`
		Company           interface{} `json:"company"`
		CreatedAt         string      `json:"created_at"`
		Email             interface{} `json:"email"`
		Followers         int64       `json:"followers"`
		FollowersURL      string      `json:"followers_url"`
		Following         int64       `json:"following"`
		FollowingURL      string      `json:"following_url"`
		GistsURL          string      `json:"gists_url"`
		GravatarID        string      `json:"gravatar_id"`
		Hireable          interface{} `json:"hireable"`
		HTMLURL           string      `json:"html_url"`
		ID                int64       `json:"id"`
		Location          interface{} `json:"location"`
		Login             string      `json:"login"`
		Name              string      `json:"name"`
		NodeID            string      `json:"node_id"`
		OrganizationsURL  string      `json:"organizations_url"`
		PublicGists       int64       `json:"public_gists"`
		PublicRepos       int64       `json:"public_repos"`
		ReceivedEventsURL string      `json:"received_events_url"`
		ReposURL          string      `json:"repos_url"`
		SiteAdmin         bool        `json:"site_admin"`
		StarredURL        string      `json:"starred_url"`
		Subscriptio_nsURL string      `json:"subscriptio ns_url"`
		TwitterUsername   interface{} `json:"twitter_username"`
		Type              string      `json:"type"`
		UpdatedAt         string      `json:"updated_at"`
		URL               string      `json:"url"`
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var response Res
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &response)
	return response.Login == "shawnyu5"
}
