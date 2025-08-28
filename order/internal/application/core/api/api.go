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
	// Validar produtos antes de processar o pedido
	err := a.validateOrderItems(order.OrderItems)
	if err != nil {
		return domain.Order{}, err
	}

	total := 0
	for _, orderItem := range order.OrderItems {
		total += int(orderItem.Quantity)
	}

	if total > 50 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, "Total itens over 50 is not allowed")
	}

	// Atualizar estoque dos produtos
	err = a.updateProductStock(order.OrderItems)
	if err != nil {
		return domain.Order{}, err
	}

	err = a.db.Save(&order)
	if err != nil {
		// Reverter estoque em caso de erro
		a.revertProductStock(order.OrderItems)
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		// Reverter estoque em caso de erro de pagamento
		a.revertProductStock(order.OrderItems)
		return domain.Order{}, paymentErr
	}

	shippingErr := a.shipping.Ship(int(order.ID), order.OrderItems)
	if shippingErr != nil {
		return domain.Order{}, shippingErr
	}

	return order, nil
}

func (a Application) validateOrderItems(orderItems []domain.OrderItem) error {
	// Extrair códigos dos produtos
	var productCodes []string
	productQuantityMap := make(map[string]int32)

	for _, item := range orderItems {
		if item.Quantity <= 0 {
			return status.Errorf(codes.InvalidArgument, "quantity must be greater than zero for product %s", item.ProductCode)
		}
		productCodes = append(productCodes, item.ProductCode)
		productQuantityMap[item.ProductCode] += item.Quantity
	}

	// Buscar produtos no banco
	products, err := a.db.GetProducts(productCodes)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to validate products: %v", err)
	}

	// Criar mapa de produtos encontrados
	foundProducts := make(map[string]domain.Product)
	for _, product := range products {
		foundProducts[product.ProductCode] = product
	}

	// Validar se todos os produtos existem e têm estoque suficiente
	for _, productCode := range productCodes {
		product, exists := foundProducts[productCode]
		if !exists {
			return status.Errorf(codes.NotFound, "product with code '%s' not found", productCode)
		}

		if !product.Active {
			return status.Errorf(codes.FailedPrecondition, "product '%s' is not active", productCode)
		}

		requiredQuantity := productQuantityMap[productCode]
		if product.Stock < requiredQuantity {
			return status.Errorf(codes.FailedPrecondition, "insufficient stock for product '%s': required %d, available %d",
				productCode, requiredQuantity, product.Stock)
		}
	}

	return nil
}

func (a Application) updateProductStock(orderItems []domain.OrderItem) error {
	productQuantityMap := make(map[string]int32)

	for _, item := range orderItems {
		productQuantityMap[item.ProductCode] += item.Quantity
	}

	for productCode, quantity := range productQuantityMap {
		err := a.db.UpdateProductStock(productCode, quantity)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to update stock for product %s: %v", productCode, err)
		}
	}

	return nil
}

func (a Application) revertProductStock(orderItems []domain.OrderItem) {
	productQuantityMap := make(map[string]int32)

	for _, item := range orderItems {
		productQuantityMap[item.ProductCode] += item.Quantity
	}

	for productCode, quantity := range productQuantityMap {
		// Reverter o estoque (adicionar de volta)
		a.db.UpdateProductStock(productCode, -quantity)
	}
}
