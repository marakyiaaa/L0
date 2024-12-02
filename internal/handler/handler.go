package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"l0/internal/model"
	"log"
	"net/http"
)

type Service interface {
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

// Получаем заказ по ID
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
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

	// Возвращаем JSON для API
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при кодировании ответа: %v", err), http.StatusInternalServerError)
			return
		}
		return
	}
	// Рендерим страницу с данными заказа
	tmpl, err := template.ParseFiles("internal/handler/templates/order.html")
	if err != nil {
		log.Printf("Ошибка загрузки шаблона: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Параметры для рендера
	data := struct {
		Order interface{}
		Error string
	}{
		Order: order,
		Error: "",
	}

	tmpl.Execute(w, data)
}

func (h *Handler) RenderHTML(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/handler/templates/order.html")
	if err != nil {
		log.Printf("Ошибка загрузки шаблона", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Параметры для рендера
	data := struct {
		Order interface{}
		Error string
	}{
		Order: nil,
		Error: "",
	}

	tmpl.Execute(w, data)
}
