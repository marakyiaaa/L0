package handler

import (
	"encoding/json"
	"fmt"
	"l0/internal/model"
	"net/http"
)

type Service interface {
	//GetOrders() ([]model.Order, error)
	GetOrder(id string) (model.Order, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Получаем все заказы
//func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
//	orders, err := h.service.GetOrders()
//	if err != nil {
//		http.Error(w, fmt.Sprintf("Не удалось получить заказы: %v", err), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	if err := json.NewEncoder(w).Encode(orders); err != nil {
//		http.Error(w, fmt.Sprintf("Ошибка при кодировании ответа: %v", err), http.StatusInternalServerError)
//		return
//	}
//}

// Получаем заказ по ID
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL, например: /order/{id}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID заказа не передан", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось найти заказ с ID %s: %v", id, err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при кодировании ответа: %v", err), http.StatusInternalServerError)
		return
	}
}
