package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "test-repo/pkg/api/test"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderService struct {
	pb.UnimplementedOrderServiceServer
	mu     sync.Mutex
	orders map[string]*pb.Order
}

func newOrderService() *orderService {
	return &orderService{
		orders: make(map[string]*pb.Order),
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()

	order := &pb.Order{
		Id:       id,
		Item:     req.Item,
		Quantity: req.Quantity,
	}

	s.orders[id] = order

	return &pb.CreateOrderResponse{
		Id: id,
	}, nil
}

func (s *orderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}

	return &pb.GetOrderResponse{
		Order: order,
	}, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}

	order.Item = req.Item
	order.Quantity = req.Quantity

	return &pb.UpdateOrderResponse{
		Order: order,
	}, nil
}

func (s *orderService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.orders[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}

	delete(s.orders, req.Id)

	return &pb.DeleteOrderResponse{
		Success: true,
	}, nil
}

func (s *orderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var orders []*pb.Order

	for _, order := range s.orders {
		orders = append(orders, order)
	}

	return &pb.ListOrdersResponse{
		Orders: orders,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	service := newOrderService()

	pb.RegisterOrderServiceServer(grpcServer, service)

	log.Println("Server started on port 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
