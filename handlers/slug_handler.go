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

	data := elastic.SearchParams{}
	data.Query = &elastic.Query{
		Match: map[string]string{"slug": slug},
	}

	elasticResponse, err := elastic.Client(data)
	if err != nil {
    fmt.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
		return
	}

  hits := elasticResponse.Hits.Hits

	w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)

	if len(hits) > 0 {
		article := elastic.BuildArticle(hits[0])
		json.NewEncoder(w).Encode(article)
	} else {
		json.NewEncoder(w).Encode(map[string]string{})
	}
}