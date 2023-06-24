package main

import (
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/malarcon-79/microservices-lab/02-servicios-backend/custody-service/controllers"

	pb "github.com/malarcon-79/microservices-lab/grpc-protos-go/system/custody"

	"go.uber.org/zap"

	"google.golang.org/grpc"

	"github.com/malarcon-79/microservices-lab/orm-go/dao"
	"github.com/malarcon-79/microservices-lab/orm-go/initdb"
	"github.com/malarcon-79/microservices-lab/orm-go/model"
)

func main() {
	_logger, _ := zap.NewProduction()
	logger := _logger.Sugar()

	var opts []grpc.ServerOption
	opts = append(opts, grpc.MaxRecvMsgSize(math.MaxInt32))
	opts = append(opts, grpc.MaxSendMsgSize(math.MaxInt32))
	server := grpc.NewServer(opts...)

	// registrar stubs Servidor
	controller, _ := controllers.NewCustodyServiceController()
	pb.RegisterCustodyServiceServer(server, &controller)

	errChannel := make(chan error)
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, syscall.SIGTERM, syscall.SIGINT)

	dsn := os.Getenv("dsn")
	if len(dsn) == 0 {
		dsn = "host=psql-backend.backend port=5432 user=postgres password=postgres dbname=lab sslmode=disable"
	}

	if err := initdb.DAOInit(dsn, "psql"); err != nil {
		logger.Panic(err)
	}
	if err := dao.DB.AutoMigrate(&model.Custody{}); err != nil {
		logger.Panic("no se pudo migrar tablas", err)
	}

	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:5000")
		if err != nil {
			errChannel <- err
			return
		}

		logger.Info("Iniciando Servidor gRPC")
		if err := server.Serve(listener); err != nil {
			errChannel <- err
			return
		}
	}()

	// Shutdown hook
	defer func() {
		server.GracefulStop()
	}()

	select {
	case err := <-errChannel:
		logger.Error("Error fatal", err)
	case <-stopChannel:
		logger.Warn("Finalizando servidor")
	}
}
