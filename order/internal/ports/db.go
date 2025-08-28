package ports

import "github.com/lucasandrade912/microservices/order/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error
	GetProduct(productCode string) (domain.Product, error)
	GetProducts(productCodes []string) ([]domain.Product, error)
	UpdateProductStock(productCode string, quantity int32) error
}
