package graphql

import "github.com/aeramu/menfess-backend/service"

func NewResolver(svc service.Service) *Resolver {
	return &Resolver{svc: svc}
}

type Resolver struct {
	svc service.Service
}
