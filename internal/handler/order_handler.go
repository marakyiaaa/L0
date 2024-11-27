package handler

import (
	"encoding/json"
	"l0/internal/service"
	"net/http"
)

type OrderHandler struct {
	service *service.OrderService
}

// консруктор
// принимает объект OrderService, который будет использоваться для обработки запросов в этом обработчике.
func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("id")

	//ищем по id если нет то ошибка 400
	if orderUID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	//пытаемся получить заказ из бд или кэша, если ошибка то 404
	order, err := h.service.GetOrderByID(orderUID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)

}
