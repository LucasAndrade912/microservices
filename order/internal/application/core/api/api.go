package api

import (
	"github.com/lucasandrade912/microservices/order/internal/application/core/domain"
	"github.com/lucasandrade912/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db:       db,
		payment:  payment,
		shipping: shipping,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	total := 0

	for _, orderItem := range order.OrderItems {
		total += int(orderItem.Quantity)
	}

	if total > 50 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, "Total itens over 50 is not allowed")
	}

	err := a.db.Save(&order)

	if err != nil {
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(&order)

	if paymentErr != nil {
		return domain.Order{}, paymentErr
	}

	shippingErr := a.shipping.Ship(int(order.ID), order.OrderItems)

	if shippingErr != nil {
		return domain.Order{}, shippingErr
	}

	return order, nil
}
