//go:build integration
// +build integration

package api_test

import (
	"fmt"
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
}

func setUp(t *testing.T) *fixture {
	mc, err := storage.NewMongoClient(getMongoURI())
	test.AssertNoError(t, err)
	test.AssertNotNil(t, mc)
	a := api.New(mc)
	l := logrus.New()
	md := httpmiddleware.New(l)
	a.AddHandlers(md)
	return &fixture{md: md}
}

func getMongoURI() string {
	mongoPort := strings.ReplaceAll(os.Getenv("MONGO_PORT"), "tcp://", "")
	return fmt.Sprintf("mongodb://%s", mongoPort)
}

func TestNewApiHandler(t *testing.T) {
	setUp(t)
}
