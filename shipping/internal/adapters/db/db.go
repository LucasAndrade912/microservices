package db

import (
	"context"
	"fmt"

	"github.com/lucasandrade912/microservices/shipping/internal/application/core/domain"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Shipping struct {
	gorm.Model
	PurchaseID int64
	Items      []Item `gorm:"foreignKey:ShippingID"`
}

type Item struct {
	gorm.Model
	ShippingID  uint
	ProductCode string
	UnitPrice   float64
	Quantity    int32
}

type Adapter struct {
	db *gorm.DB
}

func (a Adapter) Get(ctx context.Context, id string) (domain.Shipping, error) {
	var shippingEntity Shipping
	res := a.db.WithContext(ctx).First(&shippingEntity, id)
	var items []domain.Item
	for _, item := range shippingEntity.Items {
		items = append(items, domain.Item{
			ProductCode: item.ProductCode,
			UnitPrice:   float32(item.UnitPrice),
			Quantity:    item.Quantity,
		})
	}
	shipping := domain.Shipping{
		PurchaseID: shippingEntity.PurchaseID,
		Items:      items,
	}
	return shipping, res.Error
}

func (a Adapter) Save(ctx context.Context, shipping *domain.Shipping) error {
	var items []Item
	for _, item := range shipping.Items {
		items = append(items, Item{
			ProductCode: item.ProductCode,
			UnitPrice:   float64(item.UnitPrice),
			Quantity:    item.Quantity,
		})
	}
	shippingModel := Shipping{
		PurchaseID: shipping.PurchaseID,
		Items:      items,
	}
	res := a.db.WithContext(ctx).Create(&shippingModel)
	if res.Error == nil {
		shipping.ID = int64(shippingModel.ID)
	}
	return res.Error
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName("shipping"))); err != nil {
		return nil, fmt.Errorf("db otel plugin error: %v", err)
	}

	err := db.AutoMigrate(&Shipping{}, &Item{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}
