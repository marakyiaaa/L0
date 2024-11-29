package handler

import (
	"l0/internal/model"
	"net/http"
)

type Service interface {
	GetOrders() ([]model.Order, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	if resp, err := h.service.GetOrders(); err != nil {
		// status err
	}
	w.Write(resp)
}
