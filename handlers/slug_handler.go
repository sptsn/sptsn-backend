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

	var results []elastic.Article
	if len(articles) == 0 {
		results = make([]elastic.Article, 0)
	} else {
		results = articles[:1]
	}

	response := ApiResponse{
		Results: results,
		Total: len(results),
	}

	w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}