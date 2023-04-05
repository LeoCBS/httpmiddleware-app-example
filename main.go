package main

import (
	"net/http"

	"github.com/LeoCBS/httpmiddleware"
	"github.com/sirupsen/logrus"
)

type user struct {
	Login string `json:"login"`
}

func main() {
	l := logrus.New()
	md := httpmiddleware.New(l)
	md.POST("/name/:name", createUser)
	md.GET("/name/:name", getUser)

	s := &http.Server{
		Addr:    ":8080",
		Handler: md,
	}
	panic(s.ListenAndServe())
}

func createUser(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
	return httpmiddleware.Response{
		StatusCode: http.StatusCreated,
		Body:       user{Login: ps.ByName("name")},
	}
}

func getUser(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
	return httpmiddleware.Response{
		StatusCode: http.StatusCreated,
		Body:       user{Login: ps.ByName("name")},
	}
}
