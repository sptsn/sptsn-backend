package elastic

type Hit struct {
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
  Highlight struct {
    Content  []string `json:"content"`
  } `json:"highlight"`
}

type ElasticResponse struct {
  Took     int  `json:"took"`
  TimedOut bool `json:"timed_out"`
  Shards   struct {
  } `json:"_shards"`
  Hits struct {
    Total struct {
    } `json:"total"`
    MaxScore float32 `json:"max_score"`
    Hits     []Hit `json:"hits"`
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
  Highlight   string `json:"highlight"`
}

type MultiMatch struct {
  Query string `json:"query"`
  Fields []string `json:"fields"`
}

type Query struct {
  MultiMatch *MultiMatch `json:"multi_match,omitempty"`
  Match map[string]string `json:"match,omitempty"`
}

type Fields struct {
  Content map[string]int `json:"content"`
}

type Highlight struct {
  Fields *Fields `json:"fields"`
}

type SearchParams struct {
  Source []string `json:"_source,omitempty"`
  Sort map[string]string `json:"sort,omitempty"`
  Query *Query `json:"query,omitempty"`
  Highlight *Highlight `json:"highlight,omitempty"`
}

const BaseUrl = "https://sptsn.ru/elastic"
// const BaseUrl = "http://localhost:9200"