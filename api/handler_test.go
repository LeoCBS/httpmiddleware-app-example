//go:build integration
// +build integration

package api_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/LeoCBS/httpmiddleware"
	"github.com/LeoCBS/httpmiddleware-app-example/api"
	"github.com/LeoCBS/httpmiddleware-app-example/storage"
	"github.com/LeoCBS/httpmiddleware-app-example/test"
	"github.com/sirupsen/logrus"
)

type fixture struct {
	md *httpmiddleware.Middleware
	mc *storage.MongoClient
}

func setUp(t *testing.T) *fixture {
	mc, err := storage.NewMongoClient(getMongoURI())
	test.AssertNoError(t, err)
	test.AssertNotNil(t, mc)
	a := api.New(mc)
	l := logrus.New()
	md := httpmiddleware.New(l)
	a.AddHandlers(md)
	return &fixture{
		md: md,
		mc: mc,
	}
}

func getMongoURI() string {
	mongoPort := strings.ReplaceAll(os.Getenv("MONGO_PORT"), "tcp://", "")
	return fmt.Sprintf("mongodb://%s", mongoPort)
}

func TestNewApiHandler(t *testing.T) {
	setUp(t)
}

func TestInsertAndFindRecord(t *testing.T) {
	f := setUp(t)
	database := "whatever"
	collection := "kmelo"
	defer tearDown(t, f.mc, database)
	URL := fmt.Sprintf(
		"/storage/mongodb/database/%s/collection/%s/",
		database,
		collection,
	)
	bodyStr := []byte(`{"name":"leonardo"}`)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(bodyStr))
	test.AssertNoError(t, err)

	recorder := httptest.NewRecorder()
	f.md.ServeHTTP(recorder, req)
	resp := recorder.Result()
	test.AssertEqual(t, http.StatusCreated, resp.StatusCode)

	location := resp.Header.Get("Location")
	req, err = http.NewRequest("GET", location, nil)
	test.AssertNoError(t, err)

	recorder = httptest.NewRecorder()
	f.md.ServeHTTP(recorder, req)
	resp = recorder.Result()
	test.AssertEqual(t, http.StatusOK, resp.StatusCode)

}

func tearDown(t *testing.T, mc *storage.MongoClient, database string) {
	test.AssertNoError(t, mc.DropDatabase(context.Background(), database))
}

func TestInsertBadRequest(t *testing.T) {
	f := setUp(t)
	database := "whatever"
	collection := "kmelo"
	defer tearDown(t, f.mc, database)
	URL := fmt.Sprintf(
		"/storage/mongodb/database/%s/collection/%s/",
		database,
		collection,
	)
	bodyStrWrong := []byte(`{"name":"leonardo"`)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(bodyStrWrong))
	test.AssertNoError(t, err)

	recorder := httptest.NewRecorder()
	f.md.ServeHTTP(recorder, req)
	resp := recorder.Result()
	test.AssertEqual(t, http.StatusBadRequest, resp.StatusCode)
	expectedResponseBody := `{"error":"error to decode request JSON body / err = unexpected EOF"}`
	test.AssertBodyContains(t, resp.Body, expectedResponseBody)
}

func TestFindNotFound(t *testing.T) {
	f := setUp(t)
	database := "whatever"
	collection := "kmelo"
	defer tearDown(t, f.mc, database)
	URL := fmt.Sprintf(
		"/storage/mongodb/database/%s/collection/%s/id/64307e0190256830aa282c31",
		database,
		collection,
	)
	req, err := http.NewRequest("GET", URL, nil)
	test.AssertNoError(t, err)

	recorder := httptest.NewRecorder()
	f.md.ServeHTTP(recorder, req)
	resp := recorder.Result()
	test.AssertEqual(t, http.StatusNotFound, resp.StatusCode)
	expectedResponseBody := `{"error":"id {64307e0190256830aa282c31} not found"}`
	test.AssertBodyContains(t, resp.Body, expectedResponseBody)
}
