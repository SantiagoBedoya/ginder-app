package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/SantiagoBedoya/ginder_oauth-api/api"
	"github.com/SantiagoBedoya/ginder_oauth-api/oauth"
	"github.com/SantiagoBedoya/ginder_oauth-api/repositories/postgres"
	"github.com/gin-gonic/gin"
)

func main() {
	tokenAccessRepo, err := postgres.NewPostgresRepository(os.Getenv("DATABASE_URL"), &oauth.AccessToken{})
	if err != nil {
		log.Fatal(err)
	}
	port := "3000"
	if strings.TrimSpace(os.Getenv("PORT")) != "" {
		port = strings.TrimSpace(os.Getenv("PORT"))
	}
	router := gin.Default()
	api.SetupRoutes(router, tokenAccessRepo)

	errs := make(chan error, 2)
	go func() {
		log.Println("starting application on port :" + port)
		errs <- http.ListenAndServe(":"+port, router)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Printf("Termininated %s\n", <-errs)
}
