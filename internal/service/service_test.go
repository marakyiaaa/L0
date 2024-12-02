package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"l0/internal/model"
	"testing"
)

type MockOrderCache struct {
	mock.Mock
}

func (m *MockOrderCache) IsEmpty() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockOrderCache) CreateOrder(order model.Order) {
	m.Called(order)
}

func (m *MockOrderCache) GetOrder(id string) (model.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return model.Order{}, args.Error(1)
	}
	return args.Get(0).(model.Order), args.Error(1)
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateOrder(order model.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockRepository) GetOrder(id string) (model.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return model.Order{}, args.Error(1)
	}
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *MockRepository) GetOrders() ([]model.Order, error) {
	args := m.Called()
	return args.Get(0).([]model.Order), args.Error(1)
}

func TestGetOrderFromCache(t *testing.T) {
	mockCache := new(MockOrderCache)
	mockRepo := new(MockRepository)
	order := model.Order{
		Order_uid: "zDxf8iwo5VUJ3zDRp0s",
		Items: []model.Items{
			{OrderUID: "item1"},
		},
	}

	mockCache.On("GetOrder", "zDxf8iwo5VUJ3zDRp0s").Return(order, nil)

	service := &Service{repository: mockRepo, orderCache: mockCache}
	result, err := service.GetOrder("zDxf8iwo5VUJ3zDRp0s")

	assert.NoError(t, err)
	assert.Equal(t, order, result)

	mockCache.AssertExpectations(t)
}

func TestCreateOrder_InvalidPaymentAmount(t *testing.T) {
	mockCache := new(MockOrderCache)
	mockRepo := new(MockRepository)

	order := model.Order{
		Order_uid: "123",
		Payment: model.Payment{
			Amount: -100,
		},
		Items: []model.Items{
			{OrderUID: "item1"},
		},
	}

	service := &Service{repository: mockRepo, orderCache: mockCache}

	err := service.CreateOrder(order)
	assert.Error(t, err)
	assert.Equal(t, "невалидный заказ", err.Error())

	mockRepo.AssertNotCalled(t, "CreateOrder")
	mockCache.AssertNotCalled(t, "CreateOrder")
}

func TestNew(t *testing.T) {
	mockCache := new(MockOrderCache)
	mockRepo := new(MockRepository)

	orders := []model.Order{
		{Order_uid: "order1"},
		{Order_uid: "order2"},
	}

	mockCache.On("IsEmpty").Return(true)
	mockRepo.On("GetOrders").Return(orders, nil)
	mockCache.On("CreateOrder", orders[0]).Return(nil)
	mockCache.On("CreateOrder", orders[1]).Return(nil)

	service := New(mockRepo, mockCache)

	assert.NotNil(t, service)
	mockCache.AssertCalled(t, "IsEmpty")
	mockRepo.AssertCalled(t, "GetOrders")
	mockCache.AssertCalled(t, "CreateOrder", orders[0])
	mockCache.AssertCalled(t, "CreateOrder", orders[1])
}

func TestGetOrder_FromDatabase(t *testing.T) {
	mockCache := new(MockOrderCache)
	mockRepo := new(MockRepository)

	order := model.Order{
		Order_uid: "order123",
		Items:     []model.Items{{OrderUID: "item1"}},
	}

	mockCache.On("GetOrder", "order123").Return(model.Order{}, errors.New("not found"))
	mockRepo.On("GetOrder", "order123").Return(order, nil)
	mockCache.On("CreateOrder", order).Return(nil)

	service := &Service{repository: mockRepo, orderCache: mockCache}
	result, err := service.GetOrder("order123")

	assert.NoError(t, err)
	assert.Equal(t, order, result)

	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestGetOrder_NotFound(t *testing.T) {
	mockCache := new(MockOrderCache)
	mockRepo := new(MockRepository)

	mockCache.On("GetOrder", "order123").Return(model.Order{}, errors.New("not found"))
	mockRepo.On("GetOrder", "order123").Return(model.Order{}, errors.New("not found"))

	service := &Service{repository: mockRepo, orderCache: mockCache}
	_, err := service.GetOrder("order123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не удалось найти заказ")
}
func TestGetOrders_Success(t *testing.T) {
	mockCache := new(MockOrderCache)
	mockRepo := new(MockRepository)

	orders := []model.Order{
		{Order_uid: "order1"},
		{Order_uid: "order2"},
	}
	mockRepo.On("GetOrders").Return(orders, nil)

	service := &Service{repository: mockRepo, orderCache: mockCache}
	result, err := service.GetOrders()

	assert.NoError(t, err)
	assert.Equal(t, orders, result)

	mockRepo.AssertExpectations(t)
}
