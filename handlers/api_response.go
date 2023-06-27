package handlers

import(
	"github.com/sptsn/sptsn-backend/elastic"
)

type ApiResponse struct {
	Results []elastic.Article `json:"results"`
	Total int `json:"total"`
}