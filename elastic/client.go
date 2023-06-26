package elastic

import(
	"net/http"
  "fmt"
  "encoding/json"
  "bytes"
)

func Client(data SearchParams) (*ElasticResponse, error) {
	json_data, _ := json.Marshal(data)
  resp, err := http.Post(BaseUrl + "/articles/_search", "application/json", bytes.NewBuffer(json_data))

  if err != nil {
		return nil, err
  }

  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    return nil, fmt.Errorf("StatusCode: %w", resp.StatusCode)
  }

  elasticResponse := &ElasticResponse{}
  err = json.NewDecoder(resp.Body).Decode(&elasticResponse)
  if err != nil {
    return nil, err
  }

	return elasticResponse, nil
}

func BuildArticle(hit Hit) Article {
	return Article {
		ID: hit.ID,
		Title: hit.Source.Title,
		Description: hit.Source.Description,
		Date: hit.Source.Date,
		Slug: hit.Source.Slug,
		Rewritten: hit.Source.Rewritten,
		Tags: hit.Source.Tags,
		Content: hit.Source.Content,
	}
}