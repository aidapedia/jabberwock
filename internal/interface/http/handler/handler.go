package handler

import "github.com/aidapedia/jabberwock/internal/usecase/authenticated"

type Handler struct {
	usecase authenticated.Interface
}

func NewHandler(usecase authenticated.Interface) *Handler {
	return &Handler{usecase: usecase}
}
