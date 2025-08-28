package ports

import "github.com/lucasandrade912/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	Ship(purchaseId int, items []domain.OrderItem) error
}
