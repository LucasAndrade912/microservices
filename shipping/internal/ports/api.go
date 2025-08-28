package ports

import (
	"context"

	"github.com/lucasandrade912/microservices/shipping/internal/application/core/domain"
)

type APIPort interface {
	Ship(ctx context.Context, shipping domain.Shipping) (int, error)
}
