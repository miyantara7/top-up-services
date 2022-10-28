package client

import (
	proto "github.com/vins7/module-protos/app/interface/grpc/proto/e_money_service"
	"github.com/vins7/top-up-services/app/interface/model"
	"github.com/vins7/top-up-services/app/util"
	cfgServer "github.com/vins7/top-up-services/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EMoneyClient struct {
	client proto.BillerServiceClient
}

func NewEMoneyClient(c *grpc.ClientConn) EMoneyRepo {
	return &EMoneyClient{
		client: proto.NewBillerServiceClient(c),
	}
}

func (e *EMoneyClient) GetBalance(req interface{}) (*proto.GetBalanceResponse, error) {

	data, ok := req.(*model.GetBalance)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Error casting GetBalance")
	}

	res, err := e.client.GetBalance(
		util.BuildContext(cfgServer.GetConfig().Server.App.Name),
		&proto.GetBalanceRequest{
			UserId:   data.UserId,
			UserName: data.UserName,
			NoKartu:  data.NoKartu,
		})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *EMoneyClient) DetailBiller(req interface{}) (*proto.Biller, error) {
	data, ok := req.(*model.BillerRequest)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Error casting DetailBiller")
	}

	res, err := e.client.DetailBiller(
		util.BuildContext(cfgServer.GetConfig().Server.App.Name),
		&proto.BillerRequest{
			ID: data.ID,
		})
	if err != nil {
		return nil, err
	}
	return res, nil
}
