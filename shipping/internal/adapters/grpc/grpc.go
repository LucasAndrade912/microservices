package grpc

import (
	"context"
	"fmt"

	"github.com/lucasandrade912/microservices-proto/golang/shipping"
	"github.com/lucasandrade912/microservices/shipping/internal/application/core/domain"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Adapter) Create(ctx context.Context, request *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	log.WithContext(ctx).Info("Creating shipping...")

	// Convert []*shipping.Item to []domain.Item
	var domainItems []domain.Item
	for _, item := range request.Items {
		domainItems = append(domainItems, domain.Item{
			ProductCode: item.ProductCode,
			UnitPrice:   float32(item.UnitPrice),
			Quantity:    int32(item.Quantity),
		})
	}

	newShipping := domain.NewShipping(request.PurchaseId, domainItems)
	result, err := a.api.Ship(ctx, newShipping)
	code := status.Code(err)

	if code == codes.InvalidArgument || code == codes.DeadlineExceeded {
		return nil, err
	} else if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to create shipping. %v ", err)).Err()
	}

	return &shipping.CreateShippingResponse{DeadlineInDays: int64(result)}, nil
}
