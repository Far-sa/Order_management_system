package gateway

import (
	pb "github.com/Far-sa/commons/api"

	"context"
)

type OrdersGateway interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error)
}
