package grpc

import (
	"net"
	"order/proto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type OrderGRPCServer struct {
	proto.UnimplementedOrderServiceServer
}

func GRPCListen() {
	conn, err := net.Listen("tcp", ":30001")
	if err != nil {
		logrus.Fatalf("failed to listen for gRPC: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterOrderServiceServer(srv, OrderGRPCServer{})
	logrus.Infof("gRPC Server started on port 30001")

	if err := srv.Serve(conn); err != nil {
		logrus.Fatalf("error when connect to gRPC Server  %v", err)
	}
}
