package api

import (
	"context"

	"github.com/lucasandrade912/microservices/shipping/internal/application/core/domain"
	"github.com/lucasandrade912/microservices/shipping/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (a Application) Ship(ctx context.Context, shipping domain.Shipping) (int, error) {
	err := a.db.Save(ctx, &shipping)
	if err != nil {
		return 0, err
	}

	totalUnits := 0

	for _, item := range shipping.Items {
		totalUnits += int(item.Quantity)
	}

	deadline := 1

	additionalDays := totalUnits / 5
	deadline += additionalDays

	return deadline, nil
}
