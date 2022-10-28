package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	cfg "github.com/vins7/top-up-services/config"
	cfgClient "github.com/vins7/top-up-services/config/client"
	"google.golang.org/grpc"
)

var (
	EMoneyConn *grpc.ClientConn
)

func init() {
	var err error
	config := cfg.GetConfig()
	EMoneyConn, err = OpenNewConnection(config.Client.EMoneyServices)
	if err != nil {
		log.Fatal("Not connected err =>", err)
	}
}

func OpenNewConnection(config cfgClient.Server) (*grpc.ClientConn, error) {
	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(config.Timeout)*time.Millisecond)

	conn, err := grpc.DialContext(ctx,
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
