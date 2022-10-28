package client

import proto "github.com/vins7/module-protos/app/interface/grpc/proto/e_money_service"

type EMoneyRepo interface {
	GetBalance(in interface{}) (*proto.GetBalanceResponse, error)
	DetailBiller(in interface{}) (*proto.Biller, error)
}
