package graphql

import (
	"context"
	"encoding/json"
	"github.com/aeramu/menfess-backend/constants"
	handler "github.com/aeramu/menfess-backend/handler/graphql"
	"github.com/aeramu/menfess-backend/service"
	"github.com/graph-gophers/graphql-go"
	"io/ioutil"
	"net/http"
	"os"
)

func NewServer(svc service.Service) (*server, error) {
	f, err := os.Open("schema.graphql")
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	schema, err := graphql.ParseSchema(string(b), handler.NewResolver(svc), graphql.UseFieldResolvers())
	if err != nil {
		return nil, err
	}

	return &server{
		Schema: schema,
	}, nil
}

type server struct {
	*graphql.Schema
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), constants.AuthorizationKey, r.Header.Get("Authorization"))

	response := s.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJSON)
}
