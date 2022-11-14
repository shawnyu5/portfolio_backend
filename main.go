package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!", r.URL.Path[1:])
	})
	http.HandleFunc("/login", githubLoginHandler)
	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, "https://github.com/login/oauth/authorize", nil)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	// TODO: consider adding a random state to query params
	q.Add("client_id", "Iv1.3e29d7aba57c28f1")
	req.URL.RawQuery = q.Encode()

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, req.URL.String(), http.StatusFound)

}
