package handler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"l0/internal/model"
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

	data := struct {
		Order model.Order
	}{
		Order: order,
	}

	tmpl, err := template.ParseFiles("internal/handler/templates/order.html")
	if err != nil {
		logrus.WithField("Ошибка загрузки шаблона: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		logrus.WithField("Ошибка выполнения шаблона", err).Error()
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RenderHTML(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/handler/templates/order.html")
	if err != nil {
		logrus.WithField("Ошибка загрузки шаблона", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := struct {
		Order interface{}
	}{
		Order: "",
	}

	if err := tmpl.Execute(w, data); err != nil {
		logrus.WithField("Ошибка выполнения шаблона", err).Error()
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
