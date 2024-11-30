package handler

import (
	"encoding/json"
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
	resp, err := h.service.GetOrders()
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	//w.Write(resp)
}

//func (h *Handler) GetOrderID()  {
//
//}
