package top_up

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	proto "github.com/vins7/module-protos/app/interface/grpc/proto/top_up_service"
	uc "github.com/vins7/top-up-services/app/usecase/top_up"
)

type TopUpService struct {
	uc uc.TopUp
}

func NewTopUpService(uc uc.TopUp) *TopUpService {
	return &TopUpService{
		uc: uc,
	}
}

func (t *TopUpService) TopUp(ctx context.Context, req *proto.TopUpRequest) (*empty.Empty, error) {
	if err := t.uc.TopUpBalance(req); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (t *TopUpService) Payment(ctx context.Context, req *proto.PaymentRequest) (*empty.Empty, error) {
	if err := t.uc.Payment(req); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
