package storage_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/LeoCBS/httpmiddleware-app-example/storage"
	"github.com/LeoCBS/httpmiddleware-app-example/test"
)

func TestNewMongoClient(t *testing.T) {
	mc, err := storage.NewMongoClient(getMongoURI())
	test.AssertNoError(t, err)
	test.AssertNotNil(t, mc)
}

func getMongoURI() string {
	mongoPort := strings.ReplaceAll(os.Getenv("MONGO_PORT"), "tcp://", "")
	fmt.Println(mongoPort)
	return fmt.Sprintf("mongodb://%s", mongoPort)
}
