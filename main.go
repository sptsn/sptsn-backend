package main

import (
  "net/http"
  handlers "github.com/sptsn/sptsn-backend/handlers/articles_handler"
)

func main() {
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

  mux := http.NewServeMux()

  mux.HandleFunc("/articles", handlers.ArticlesHandler)
  http.ListenAndServe(":"+port, mux)
}
