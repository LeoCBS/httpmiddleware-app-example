package api

import (
	"context"
	"net/http"
	"time"

	"github.com/LeoCBS/httpmiddleware"
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
}

func (h *Handler) getRecordByID(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
	// all advantagens what you has using httpmiddleware lib:
	// * you don't need validate named URL parameter (:id)
	// * you don't need check/handling errors types
	// * you don't need write JSON response
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
