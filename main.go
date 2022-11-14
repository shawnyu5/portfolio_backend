package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(http.MethodGet, "https://github.com/login/oauth/authorize", nil)
		if err != nil {
			panic(err)
		}
		q := req.URL.Query()
		q.Add("client_id", "Iv1.3e29d7aba57c28f1")
		req.URL.RawQuery = q.Encode()

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		fmt.Printf("main res: %v\n", res) // __AUTO_GENERATED_PRINT_VAR__
		http.Redirect(w, r, req.URL.String(), http.StatusFound)
	})
	http.ListenAndServe(":8080", nil)
}
