package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/SantiagoBedoya/ginder_users-api/api"
	"github.com/SantiagoBedoya/ginder_users-api/repositories/mongodb"
	"github.com/gin-gonic/gin"
)

func main() {
	port := "3000"
	if strings.TrimSpace(os.Getenv("PORT")) != "" {
		port = strings.TrimSpace(os.Getenv("PORT"))
	}
	mongoTimeout := 10
	if strings.TrimSpace(os.Getenv("MONGO_TIMEOUT")) != "" {
		var err error
		mongoTimeout, err = strconv.Atoi(strings.TrimSpace(os.Getenv("MONGO_TIMEOUT")))
		if err != nil {
			log.Fatal("mongoTime should be a number")
		}
	}
	repo, err := mongodb.NewMongoRepository(os.Getenv("MONGO_URL"), os.Getenv("MONGO_DB"), mongoTimeout)
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	api.SetupRoutes(router, repo)

	errs := make(chan error, 2)
	go func() {
		log.Println("application is running on port :" + port)
		errs <- http.ListenAndServe(":"+port, router)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Printf("Terminated %s\n", <-errs)
}
