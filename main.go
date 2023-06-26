package main

import (
  "net/http"
  "os"
  "github.com/sptsn/sptsn-backend/handlers"
)

func main() {
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

  mux := http.NewServeMux()

  mux.HandleFunc("/articles", handlers.ArticlesHandler)
  mux.HandleFunc("/articles/", handlers.SlugHandler)

  http.ListenAndServe(":"+port, mux)
}
