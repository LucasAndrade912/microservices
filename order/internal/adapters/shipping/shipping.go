package shipping_adapter

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/lucasandrade912/microservices-proto/golang/shipping"
	"github.com/lucasandrade912/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	shipping shipping.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(
		opts,
		grpc.WithUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
				grpc_retry.WithMax(5),
				grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
			),
		),
	)

	conn, err := grpc.Dial(shippingServiceUrl, opts...)

	if err != nil {
		return nil, err
	}

	client := shipping.NewShippingClient(conn)
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) Ship(purchaseId int, items []domain.OrderItem) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	shippingItems := make([]*shipping.Item, len(items))
	for i, item := range items {
		shippingItems[i] = &shipping.Item{
			ProductCode: item.ProductCode,
			UnitPrice:   item.UnitPrice,
			Quantity:    int32(item.Quantity),
		}
	}

	_, err := a.shipping.Create(ctx, &shipping.CreateShippingRequest{
		PurchaseId: int64(purchaseId),
		Items:      shippingItems,
	})
	return err
}
