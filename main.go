package main

import (
	"context"
	"log"
	"net"
	"sync"

	api "your_module/3_1"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	api.UnimplementedOrderServiceServer
	orders map[string]*api.Order
	mu     sync.Mutex
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[string]*api.Order),
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *api.CreateOrderRequest) (*api.CreateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	s.orders[id] = &api.Order{
		Id:       id,
		Item:     req.Item,
		Quantity: req.Quantity,
	}

	return &api.CreateOrderResponse{Id: id}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *api.GetOrderRequest) (*api.GetOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "")
	}
	return &api.GetOrderResponse{Order: order}, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, req *api.UpdateOrderRequest) (*api.UpdateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "")
	}
	order.Item = req.Item
	order.Quantity = req.Quantity

	return &api.UpdateOrderResponse{Order: order}, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, req *api.DeleteOrderRequest) (*api.DeleteOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.orders[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "")
	}

	delete(s.orders, req.Id)
	return &api.DeleteOrderResponse{Success: true}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, req *api.ListOrdersRequest) (*api.ListOrdersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := make([]*api.Order, 0, len(s.orders))
	for _, order := range s.orders {
		list = append(list, order)
	}

	return &api.ListOrdersResponse{Orders: list}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterOrderServiceServer(grpcServer, NewOrderService())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
