package handlers

import (
  "net/http"
  "fmt"
  "net/url"
  "encoding/json"
	"github.com/sptsn/sptsn-backend/elastic"
)

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
  u, err := url.Parse(r.URL.String())
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("Internal server error"))
    return
  }

  query := u.Query().Get("q")
  data := elastic.SearchParams {
    Source: []string{"title", "slug", "date", "description", "tags"},
    Sort: map[string]string{"date": "desc"},
  }
  if query != "" {
    data.Query = &elastic.Query {
      MultiMatch: &elastic.MultiMatch{
        Query: query,
        Fields: []string{"content"},
      },
    }
  }

	elasticResponse, err := elastic.Client(data)
	if err != nil {
    fmt.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
		return
	}

  hits := elasticResponse.Hits.Hits
  var output = make([]elastic.Article, 0)
  for _, hit := range hits {
    output = append(output, elastic.BuildArticle(hit))
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(output)
}