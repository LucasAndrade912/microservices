package db

import (
	"fmt"
	"strings"

	"github.com/lucasandrade912/microservices/order/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	ProductCode string `gorm:"uniqueIndex;size:100"`
	Name        string `gorm:"size:255"`
	UnitPrice   float32
	Stock       int32
	Active      bool `gorm:"default:true"`
}

type Order struct {
	gorm.Model
	CustomerID int64
	Status     string
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderID     uint
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	// Adicionar parseTime=true automaticamente se não estiver presente
	if !strings.Contains(dataSourceUrl, "parseTime=true") {
		if strings.Contains(dataSourceUrl, "?") {
			dataSourceUrl += "&parseTime=true"
		} else {
			dataSourceUrl += "?parseTime=true"
		}
	}

	db, openErr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})

	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	err := db.AutoMigrate(&Product{}, &Order{}, &OrderItem{})

	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	// Inserir produtos de exemplo se não existirem
	adapter := &Adapter{db: db}
	adapter.seedProducts()

	return adapter, nil
}

func (a Adapter) Get(id string) (domain.Order, error) {
	var orderEntity Order

	res := a.db.First(&orderEntity, id)

	var orderItems []domain.OrderItem

	for _, orderItem := range orderEntity.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	order := domain.Order{
		ID:         int64(orderEntity.ID),
		CustomerID: orderEntity.CustomerID,
		Status:     orderEntity.Status,
		OrderItems: orderItems,
		CreatedAt:  orderEntity.CreatedAt.UnixNano(),
	}

	return order, res.Error
}

func (a Adapter) Save(order *domain.Order) error {
	var orderItems []OrderItem

	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	orderModel := Order{
		CustomerID: order.CustomerID,
		Status:     order.Status,
		OrderItems: orderItems,
	}

	res := a.db.Create(&orderModel)

	if res.Error == nil {
		order.ID = int64(orderModel.ID)
	}

	return res.Error
}

func (a Adapter) GetProduct(productCode string) (domain.Product, error) {
	var productEntity Product
	res := a.db.Where("product_code = ? AND active = ?", productCode, true).First(&productEntity)
	
	if res.Error != nil {
		return domain.Product{}, domain.ErrProductNotFound
	}

	product := domain.Product{
		ID:          int64(productEntity.ID),
		ProductCode: productEntity.ProductCode,
		Name:        productEntity.Name,
		UnitPrice:   productEntity.UnitPrice,
		Stock:       productEntity.Stock,
		Active:      productEntity.Active,
	}

	return product, nil
}

func (a Adapter) GetProducts(productCodes []string) ([]domain.Product, error) {
	var productEntities []Product
	res := a.db.Where("product_code IN ? AND active = ?", productCodes, true).Find(&productEntities)
	
	if res.Error != nil {
		return nil, res.Error
	}

	var products []domain.Product
	for _, entity := range productEntities {
		products = append(products, domain.Product{
			ID:          int64(entity.ID),
			ProductCode: entity.ProductCode,
			Name:        entity.Name,
			UnitPrice:   entity.UnitPrice,
			Stock:       entity.Stock,
			Active:      entity.Active,
		})
	}

	return products, nil
}

func (a Adapter) UpdateProductStock(productCode string, quantity int32) error {
	res := a.db.Model(&Product{}).Where("product_code = ?", productCode).Update("stock", gorm.Expr("stock - ?", quantity))
	return res.Error
}

func (a Adapter) seedProducts() {
	// Verificar se já existem produtos
	var count int64
	a.db.Model(&Product{}).Count(&count)
	
	if count > 0 {
		return // Já existem produtos
	}

	// Criar produtos de exemplo (sem timestamps do gorm.Model)
	products := []Product{
		{ProductCode: "ABC123", Name: "Produto A", UnitPrice: 10.50, Stock: 100, Active: true},
		{ProductCode: "XYZ789", Name: "Produto B", UnitPrice: 20.00, Stock: 50, Active: true},
		{ProductCode: "DEF456", Name: "Produto C", UnitPrice: 15.75, Stock: 75, Active: true},
		{ProductCode: "GHI999", Name: "Produto D", UnitPrice: 5.25, Stock: 200, Active: true},
		{ProductCode: "JKL111", Name: "Produto E", UnitPrice: 30.00, Stock: 25, Active: true},
	}

	for _, product := range products {
		a.db.Create(&product)
	}
}
