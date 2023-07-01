package handlers

import (
  "net/http"
  "fmt"
  "net/url"
  "encoding/json"
	// "time"
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
	data := buildSearchParams(query)
	articles, err := elastic.Client(data)
	if err != nil {
    fmt.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := ApiResponse{
		Results: articles,
		Total: len(articles),
	}

	// time.Sleep(1 * time.Second)

  w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(response)
}

func buildSearchParams(query string) elastic.SearchParams {
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
		data.Highlight = &elastic.Highlight {
			Fields: &elastic.Fields {
				Content: map[string]int{"number_of_fragments": 1},
			},
			TagsSchema: "styled",
		}
  }
	return data
}