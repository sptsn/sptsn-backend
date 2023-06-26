package handlers

import (
  "net/http"
  "encoding/json"
	"strings"
	"fmt"
	"bytes"
	"io"
	"github.com/sptsn/sptsn-backend/elastic"
)

func SlugHandler(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	var slug string = s[len(s) - 1]

	data := elastic.SearchParams{}
	data.Query = &elastic.Query{
		Match: map[string]string{"slug": slug},
	}
  json_data, _ := json.Marshal(data)
  resp, err := http.Post(elastic.BaseUrl + "/articles/_search", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
    fmt.Println(err)
  }

  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    fmt.Println(bodyString)
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

	w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)

	if len(hits) > 0 {
		hit := hits[0]
		article := elastic.Article {
			ID: hit.ID,
			Title: hit.Source.Title,
			Description: hit.Source.Description,
			Date: hit.Source.Date,
			Slug: hit.Source.Slug,
			Rewritten: hit.Source.Rewritten,
			Tags: hit.Source.Tags,
			Content: hit.Source.Content,
		}
		json.NewEncoder(w).Encode(article)
	} else {
		json.NewEncoder(w).Encode(map[string]string{})
	}
}