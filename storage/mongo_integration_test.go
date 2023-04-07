package storage_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/LeoCBS/httpmiddleware-app-example/storage"
	"github.com/LeoCBS/httpmiddleware-app-example/test"
	errors "github.com/LeoCBS/httpmiddleware/errors"
)

func TestNewMongoClient(t *testing.T) {
	setUp(t)
}

type fixture struct {
	mc *storage.MongoClient
}

func setUp(t *testing.T) *fixture {
	mc, err := storage.NewMongoClient(getMongoURI())
	test.AssertNoError(t, err)
	test.AssertNotNil(t, mc)
	return &fixture{mc: mc}
}

func getMongoURI() string {
	mongoPort := strings.ReplaceAll(os.Getenv("MONGO_PORT"), "tcp://", "")
	return fmt.Sprintf("mongodb://%s", mongoPort)
}

func TestInsertRecordAndFindByID(t *testing.T) {
	f := setUp(t)
	database := "testdb"
	collection := "person"
	type person struct {
		Name string
	}
	defer tearDown(t, f, database)
	expectedRecord := person{Name: "leo"}
	_id, err := f.mc.InsertRecord(
		context.Background(),
		database,
		collection,
		expectedRecord,
	)
	test.AssertNoError(t, err)

	recordRetrieved, err := f.mc.FindByID(context.Background(), database, collection, _id)
	test.AssertNoError(t, err)
	test.AssertNotNil(t, recordRetrieved)
}

func tearDown(t *testing.T, f *fixture, database string) {
	test.AssertNoError(t, f.mc.DropDatabase(context.Background(), database))
}

func TestFindByIDReturnInvalidHexID(t *testing.T) {
	f := setUp(t)
	database := "testdb"
	collection := "person"
	_, err := f.mc.FindByID(context.Background(), database, collection, "kmelo")
	expectedError := errors.New("error on create objectID / err = {the provided hex string is not a valid ObjectID}")
	test.AssertEqual(t, err, expectedError)
}

func TestFindByIDReturnNotFound(t *testing.T) {
	f := setUp(t)
	database := "testdb"
	collection := "person"
	_, err := f.mc.FindByID(context.Background(), database, collection, "64307e0190256830aa282c31")
	expectedError := errors.NewNotFound("id {64307e0190256830aa282c31} not found")
	test.AssertEqual(t, err, expectedError)
}
