package payment_adapter

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/lucasandrade912/microservices-proto/golang/payment"
	"github.com/lucasandrade912/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment payment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
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

	conn, err := grpc.Dial(paymentServiceUrl, opts...)

	if err != nil {
		return nil, err
	}

	client := payment.NewPaymentClient(conn)
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	return err
}
