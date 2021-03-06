package main

import (
  "database/sql"
  "fmt"
  "log"
  "net/http"

  _ "github.com/lib/pq"

  "github.com/gorilla/mux"
  "github.com/picolloo/go-market/product/infra"
  "github.com/picolloo/go-market/product/infra/ports"
  "github.com/picolloo/go-market/product/usecases"
)

func main() {
  dbUrl := fmt.Sprintf(
    "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    "go-market-db", 5432, "docker", "docker", "go-market",
  )

  db, err := sql.Open("postgres", dbUrl)
  if err != nil {
    log.Fatalf("Unable to connect to database: %s", err.Error())
  }
  defer db.Close()

  postgresProductRepo := product_infra.NewPostgresProductRepository(db)
  service := product_usecase.NewService(postgresProductRepo)

  rootRouter := mux.NewRouter()

  rootRouter.Handle(
    "/healthcheck",
    http.HandlerFunc(
      func(rw http.ResponseWriter, r *http.Request) {
        rw.Write([]byte("Healthy"))
      },
    ))


    apiRouter := rootRouter.PathPrefix("/api").Subrouter()
    product_ports.MakeProductHandlers(apiRouter, service)

    http.Handle("/", rootRouter)

    server := &http.Server{
      Addr: ":9000",
    }

    if err := server.ListenAndServe(); err != nil {
      log.Fatal(err)
    }
  }
