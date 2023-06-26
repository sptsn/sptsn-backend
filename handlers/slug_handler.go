package handlers

import (
  "net/http"
  "encoding/json"
	"strings"
	"fmt"
	"github.com/sptsn/sptsn-backend/elastic"
)

func SlugHandler(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	var slug string = s[len(s) - 1]

	data := elastic.SearchParams{
		Query: &elastic.Query{
			Match: map[string]string{"slug": slug},
		},
	}

	articles, err := elastic.Client(data)
	if err != nil {
    fmt.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)

	if len(articles) > 0 {
		json.NewEncoder(w).Encode(articles[0])
	} else {
		json.NewEncoder(w).Encode(map[string]string{})
	}
}