package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	proto "github.com/vins7/module-protos/app/interface/grpc/proto/top_up_service"
	client "github.com/vins7/top-up-services/app/adapter/client"
	dbTopUp "github.com/vins7/top-up-services/app/adapter/db/top_up"
	"github.com/vins7/top-up-services/app/infrastructure/connection/db"
	conn "github.com/vins7/top-up-services/app/infrastructure/connection/grpc"
	svcUser "github.com/vins7/top-up-services/app/service/top_up"
	uc "github.com/vins7/top-up-services/app/usecase/top_up"
	cfg "github.com/vins7/top-up-services/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServer() {

	config := cfg.GetConfig()
	grpcServer := grpc.NewServer()

	Apply(grpcServer)
	reflection.Register(grpcServer)

	svcHost := config.Server.Grpc.Host
	svcPort := config.Server.Grpc.Port

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", svcHost, svcPort))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start notpool Service gRPC server: %v", err)
		}
	}()

	fmt.Printf("gRPC server is running at %s:%d\n", svcHost, svcPort)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal := <-c
	log.Fatalf("process killed with signal: %v\n", signal.String())

}

func Apply(server *grpc.Server) {
	proto.RegisterTopUpServicesServer(server, svcUser.NewTopUpService(uc.NewTopUpUsecase(client.NewEMoneyClient(conn.EMoneyConn), dbTopUp.NewTopUpDB(db.TopUpDB, db.EMoneyDB))))
}
