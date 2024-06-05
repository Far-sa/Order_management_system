package service

import (
	"context"

	pb "github.com/Far-sa/commons/api"
	"github.com/Far-sa/payment/contract"
)

type service struct {
	processor contract.PaymentProcessor
	gateway   contract.OrdersGateway
}

func NewService(processor contract.PaymentProcessor, gateway contract.OrdersGateway) *service {
	return &service{processor: processor, gateway: gateway}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(ctx, o)
	if err != nil {
		return "", err
	}

	//TODO update order with link

	err = s.gateway.UpdateOrderAfterPaymentLink(ctx, o.ID, link)
	if err != nil {
		return "", err
	}

	return link, nil
}
