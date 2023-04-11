package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/LeoCBS/httpmiddleware"
	"github.com/LeoCBS/httpmiddleware-app-example/api"
	"github.com/LeoCBS/httpmiddleware-app-example/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	mc, err := storage.NewMongoClient(getMongoURI())
	if err != nil {
		panic(err)
	}
	a := api.New(mc)
	l := logrus.New()
	md := httpmiddleware.New(l)
	a.AddHandlers(md)

	s := &http.Server{
		Addr:    ":8080",
		Handler: md,
	}
	l.Info("starting server on port 8080")
	panic(s.ListenAndServe())
}

func getMongoURI() string {
	mongoPort := strings.ReplaceAll(os.Getenv("MONGO_PORT"), "tcp://", "")
	return fmt.Sprintf("mongodb://%s", mongoPort)
}
