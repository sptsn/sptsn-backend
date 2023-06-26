package handlers

import (
  "net/http"
  "fmt"
  "net/url"
  "encoding/json"
  "bytes"
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

  json_data, _ := json.Marshal(data)
  resp, err := http.Post(elastic.BaseUrl + "/articles/_search", "application/json", bytes.NewBuffer(json_data))

  if err != nil {
    fmt.Println(err)
  }

  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  elasticResponse := &elastic.ElasticResponse{}
  err = json.NewDecoder(resp.Body).Decode(&elasticResponse)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Println(err)
    return
  }

  hits := elasticResponse.Hits.Hits
  var output = make([]elastic.Article, 0)
  for _, hit := range hits {
    output = append(output, elastic.Article {
      ID: hit.ID,
      Title: hit.Source.Title,
      Description: hit.Source.Description,
      Date: hit.Source.Date,
      Slug: hit.Source.Slug,
      Rewritten: hit.Source.Rewritten,
      Tags: hit.Source.Tags,
      Content: hit.Source.Content,
    })
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(output)
}