package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LeoCBS/httpmiddleware"
	"github.com/LeoCBS/httpmiddleware/errors"
)

type Storage interface {
	InsertRecord(ctx context.Context, database, collection string, record interface{}) (string, error)
	FindByID(ctx context.Context, database, collection, _id string) (map[string]interface{}, error)
}

type Handler struct {
	s Storage
}

func New(s Storage) *Handler {
	return &Handler{s: s}
}

func (h *Handler) AddHandlers(md *httpmiddleware.Middleware) {
	md.GET("/storage/mongodb/database/:database/collection/:collection/id/:id", h.getRecordByID)
	md.POST("/storage/mongodb/database/:database/collection/:collection/", h.insertRecord)
}

func (h *Handler) getRecordByID(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
	// all advantagens what you has using httpmiddleware lib:
	// * you don't need validate named URL parameter (:id)
	// * you don't need check/handling errors types
	// * you don't need write JSON response
	// * you don't need manipulate response headers
	// * you don't need validate if client send correct HTTP method
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	result, err := h.s.FindByID(ctxTimeout, ps.ByName("database"), ps.ByName("collection"), ps.ByName("id"))
	if err != nil {
		return httpmiddleware.Response{Error: err}
	}
	return httpmiddleware.Response{
		Body:       result,
		StatusCode: http.StatusOK,
	}

}

func (h *Handler) insertRecord(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	record, err := decodeRequestBody(r)
	if err != nil {
		return httpmiddleware.Response{Error: errors.NewBadRequest(err.Error())}
	}
	recordID, err := h.s.InsertRecord(
		ctxTimeout,
		ps.ByName("database"),
		ps.ByName("collection"),
		record,
	)
	if err != nil {
		return httpmiddleware.Response{Error: err}
	}
	return httpmiddleware.Response{
		StatusCode: http.StatusCreated,
		Headers: getResponseHeaderLocation(
			recordID,
			ps.ByName("database"),
			ps.ByName("collection")),
	}

}

func getResponseHeaderLocation(recordID, database, collection string) map[string]string {
	headerKey := "Location"
	headerValue := fmt.Sprintf(
		"/storage/mongodb/database/%s/collection/%s/id/%s",
		database,
		collection,
		recordID,
	)
	return map[string]string{
		headerKey: headerValue,
	}
}

func decodeRequestBody(r *http.Request) (interface{}, error) {
	var target interface{}
	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		return nil, fmt.Errorf("error to decode request JSON body / err = %w", err)
	}
	return target, nil
}
