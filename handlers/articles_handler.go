package handlers

import (
  "net/http"
  "os"
  "fmt"
  "net/url"
  "encoding/json"
  "bytes"
)

type ElasticResponse struct {
  Took     int  `json:"took"`
  TimedOut bool `json:"timed_out"`
  Shards   struct {
  } `json:"_shards"`
  Hits struct {
    Total struct {
    } `json:"total"`
    MaxScore float32 `json:"max_score"`
    Hits     []struct {
      Index  string `json:"_index"`
      ID     string `json:"_id"`
      Score  float32    `json:"_score"`
      Source struct {
        Title       string   `json:"title"`
        Description string   `json:"description"`
        Date        string   `json:"date"`
        Slug        string   `json:"slug"`
        Rewritten   bool     `json:"rewritten"`
        Content     string   `json:"content"`
        Tags        []string `json:"tags"`
      } `json:"_source"`
    } `json:"hits"`
  } `json:"hits"`
}

type Article struct {
  ID          string `json:"id"`
  Title       string `json:"title"`
  Description string `json:"description"`
  Date        string  `json:"date"`
  Slug        string  `json:"slug"`
  Rewritten   bool    `json:"rewritten"`
  Content     string  `json:"content"`
  Tags        []string `json:"tags"`
}

type MultiMatch struct {
  Query string `json:"query"`
  Fields []string `json:"fields"`
}

type Query struct {
  MultiMatch MultiMatch `json:"multi_match"`
}

type ElasicParams struct {
  Source []string `json:"_source"`
  Sort map[string]string `json:"sort"`
  Query *Query `json:"query,omitempty"`
}

// const baseUrl = "https://sptsn.ru/elastic"
const baseUrl = "http://localhost:9200"

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
  u, err := url.Parse(r.URL.String())
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("Internal server error"))
    return
  }

  query := u.Query().Get("q")
  data := ElasicParams {
    Source: []string{"title", "slug", "date", "description", "tags"},
    Sort: map[string]string{"date": "desc"},
  }
  if query != "" {
    data.Query = &Query {
      MultiMatch: MultiMatch{
        Query: query,
        Fields: []string{"content"},
      },
    }
  }

  json_data, _ := json.Marshal(data)
  resp, err := http.Post(baseUrl + "/articles/_search", "application/json", bytes.NewBuffer(json_data))

  if err != nil {
    fmt.Println(err)
  }

  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  elasticResponse := &ElasticResponse{}
  err = json.NewDecoder(resp.Body).Decode(&elasticResponse)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Println(err)
    return
  }

  hits := elasticResponse.Hits.Hits
  var output = make([]Article, 0)
  for _, hit := range hits {
    output = append(output, Article {
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